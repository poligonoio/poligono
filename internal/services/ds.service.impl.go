package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/poligonoio/vega-core/internal/models"
	"github.com/poligonoio/vega-core/pkg/logger"
	"github.com/poligonoio/vega-core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataSourceServiceImpl struct {
	ctx                  context.Context
	dataSourceCollection *mongo.Collection
	infisicalService     InfisicalService
	trinoService         TrinoService
	schemaService        SchemaService
}

func NewDataSourceService(ctx context.Context, dataSourceCollection *mongo.Collection, infisicalService InfisicalService, trinoService TrinoService, schemaService SchemaService) DataSourceService {
	return &DataSourceServiceImpl{
		ctx:                  ctx,
		dataSourceCollection: dataSourceCollection,
		infisicalService:     infisicalService,
		trinoService:         trinoService,
		schemaService:        schemaService,
	}
}

func (self *DataSourceServiceImpl) GetByName(name string, organizationId string) (models.DataSource, error) {
	var dataSource models.DataSource

	filter := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "organization_id", Value: organizationId}}
	err := self.dataSourceCollection.FindOne(self.ctx, filter).Decode(&dataSource)
	if err != nil {
		return dataSource, err
	}

	dataSource.Secret, err = self.infisicalService.GetSecret(fmt.Sprintf("%s-%s", dataSource.Name, dataSource.OrganizationId))
	if err != nil {
		return dataSource, err
	}

	return dataSource, nil
}

func (self *DataSourceServiceImpl) GetAll(organizationId string) ([]models.DataSource, error) {
	var dataSources []models.DataSource = make([]models.DataSource, 0)

	filter := bson.D{bson.E{Key: "organization_id", Value: organizationId}}
	cursor, err := self.dataSourceCollection.Find(self.ctx, filter)
	if err != nil {
		return dataSources, err
	}

	if err = cursor.All(self.ctx, &dataSources); err != nil {
		return dataSources, err
	}

	return dataSources, nil
}

func (self *DataSourceServiceImpl) Create(dataSource models.DataSource) error {
	// check if data source already exist
	query := bson.D{bson.E{Key: "name", Value: dataSource.Name}, bson.E{Key: "organization_id", Value: dataSource.OrganizationId}}
	count, err := self.dataSourceCollection.CountDocuments(self.ctx, query)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("Data source with that name already exists")
	}

	_, err = self.dataSourceCollection.InsertOne(self.ctx, dataSource)
	if err != nil {
		return err
	}

	// Create secret
	if err = self.infisicalService.CreateSecret(fmt.Sprintf("%s-%s", dataSource.Name, dataSource.OrganizationId), dataSource.Secret); err != nil {
		cleanUpErr := self.Delete(dataSource.Name, dataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return err
	}

	// Create catalog
	catalogName := fmt.Sprintf("%s_%s", dataSource.Name, dataSource.OrganizationId)
	if err = self.trinoService.CreateCatalog(catalogName, dataSource.Type, dataSource.Secret); err != nil {
		cleanUpErr := self.Delete(dataSource.Name, dataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return err
	}

	// Sync schema
	createdDataSource, err := self.GetByName(dataSource.Name, dataSource.OrganizationId)
	if err != nil {
		cleanUpErr := self.Delete(dataSource.Name, dataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return err
	}

	err = self.Sync(createdDataSource.Name, createdDataSource.OrganizationId)
	if err != nil {
		cleanUpErr := self.Delete(dataSource.Name, dataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return err
	}

	return nil
}

func (self *DataSourceServiceImpl) Update(name string, organizationId string, newDataSource models.UpdateRequestDataSourceBody) error {
	updateFields := bson.D{}

	typeData := reflect.TypeOf(newDataSource)
	values := reflect.ValueOf(newDataSource)

	for i := 0; i < typeData.NumField(); i++ {
		field := typeData.Field(i)
		val := values.Field(i)
		tag := field.Tag.Get("bson")

		if !utils.IsZeroType(val) {
			update := primitive.E{Key: tag, Value: val.Interface()}
			updateFields = append(updateFields, update)
		}
	}

	update := bson.D{
		primitive.E{
			Key:   "$set",
			Value: updateFields,
		},
	}

	query := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "organization_id", Value: organizationId}}
	result, err := self.dataSourceCollection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return err
	}

	if result.MatchedCount < 1 {
		return fmt.Errorf("No data source found with name %s in the organization %s", name, organizationId)
	}

	// update data source secret
	if newDataSource.Secret != "" && newDataSource.Name != "" {
		if err = self.infisicalService.DeleteSecret(fmt.Sprintf("%s-%s", name, organizationId)); err != nil {
			return err
		}

		if err = self.infisicalService.CreateSecret(fmt.Sprintf("%s-%s", newDataSource.Name, organizationId), newDataSource.Secret); err != nil {
			return err
		}

		// Update trino catalog
		updatedDataSource, err := self.GetByName(newDataSource.Name, organizationId)
		if err != nil {
			return err
		}

		if err = self.trinoService.RemoveCatalog(self.trinoService.GetCatalogName(name, organizationId)); err != nil {
			return err
		}

		catalogName := fmt.Sprintf("%s_%s", newDataSource.Name, organizationId)
		if err = self.trinoService.CreateCatalog(catalogName, updatedDataSource.Type, newDataSource.Secret); err != nil {
			return err
		}

		// Update data source name on schema
		if err = self.schemaService.UpdateDataSourceName(name, newDataSource.Name, organizationId); err != nil {
			return err
		}
	} else if newDataSource.Secret != "" {
		if err := self.infisicalService.UpdateSecret(fmt.Sprintf("%s-%s", name, organizationId), newDataSource.Secret); err != nil {
			return err
		}

		// Update trino catalog
		updatedDataSource, err := self.GetByName(newDataSource.Name, organizationId)
		if err != nil {
			return err
		}

		if err = self.trinoService.RemoveCatalog(self.trinoService.GetCatalogName(name, organizationId)); err != nil {
			return err
		}

		catalogName := self.trinoService.GetCatalogName(name, organizationId)
		if err = self.trinoService.CreateCatalog(catalogName, updatedDataSource.Type, newDataSource.Secret); err != nil {
			return err
		}
	} else if newDataSource.Name != "" {
		currentSecret, err := self.infisicalService.GetSecret(fmt.Sprintf("%s-%s", name, organizationId))
		if err != nil {
			return err
		}

		if err = self.infisicalService.DeleteSecret(fmt.Sprintf("%s-%s", name, organizationId)); err != nil {
			return err
		}

		if err = self.infisicalService.CreateSecret(fmt.Sprintf("%s-%s", newDataSource.Name, organizationId), currentSecret); err != nil {
			return err
		}

		// Update trino catalog
		updatedDataSource, err := self.GetByName(newDataSource.Name, organizationId)
		if err != nil {
			return err
		}

		if err = self.trinoService.RemoveCatalog(self.trinoService.GetCatalogName(name, organizationId)); err != nil {
			return err
		}

		catalogName := self.trinoService.GetCatalogName(newDataSource.Name, organizationId)
		if err = self.trinoService.CreateCatalog(catalogName, updatedDataSource.Type, currentSecret); err != nil {
			return err
		}

		// Update data source name on schema
		if err = self.schemaService.UpdateDataSourceName(name, newDataSource.Name, organizationId); err != nil {
			return err
		}
	}

	return nil
}

func (self *DataSourceServiceImpl) Delete(name string, organizationId string) error {
	filter := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "organization_id", Value: organizationId}}
	result, err := self.dataSourceCollection.DeleteOne(self.ctx, filter)
	if err != nil {
		logger.Error.Printf("Error deleting data source document: %v", err)
	}

	if result.DeletedCount < 1 {
		logger.Error.Printf("No Data source found with name %s in the organization %s", name, organizationId)
	}

	if err = self.infisicalService.DeleteSecret(fmt.Sprintf("%s-%s", name, organizationId)); err != nil {
		logger.Error.Printf("Error deleting secret: %v", err)
	}

	if err = self.trinoService.RemoveCatalog(self.trinoService.GetCatalogName(name, organizationId)); err != nil {
		logger.Error.Printf("Error deleting catalog: %v", err)
	}

	if err = self.schemaService.Delete(name, organizationId); err != nil {
		logger.Error.Printf("Error deleting schema document: %v", err)
	}

	return nil
}

func (self *DataSourceServiceImpl) GetDataSourceSchemas(dataSourceName string, organizationId string) ([]models.Schema, error) {
	var schemas []models.Schema

	err := self.schemaService.GetAll(dataSourceName, organizationId, &schemas)
	if err != nil {
		return schemas, err
	}

	return schemas, nil
}

func (self *DataSourceServiceImpl) Sync(dataSourceName string, organizationId string) error {
	catalogName := self.trinoService.GetCatalogName(dataSourceName, organizationId)

	// use trino get catalog schema instead of writing it all again
	var psqlSchemas []models.PSQLSchema
	err := self.trinoService.Query(fmt.Sprintf("SELECT nspname FROM %s.pg_catalog.pg_namespace WHERE nspname NOT IN ('pg_toast', 'pg_catalog', 'public', 'information_schema')", catalogName), &psqlSchemas)
	if err != nil {
		return err
	}

	for _, psqlSchema := range psqlSchemas {
		var psqlTables []models.PSQLTable
		err := self.trinoService.Query(fmt.Sprintf("SELECT relname, relnamespace FROM %s.pg_catalog.pg_class c INNER JOIN %s.pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = '%s'", catalogName, catalogName, psqlSchema.Nspname), &psqlTables)
		if err != nil {
			return err
		}

		var tables []models.Table
		for _, psqlTable := range psqlTables {

			var psqlFields []models.PSQLField
			err := self.trinoService.Query(fmt.Sprintf("SELECT * FROM %s.pg_catalog.pg_attribute a INNER JOIN %s.pg_catalog.pg_class c ON a.attrelid = c.oid WHERE c.relname = '%s'", catalogName, catalogName, psqlTable.Relname), &psqlFields)
			if err != nil {
				return err
			}

			var fields []models.Field
			for _, psqlField := range psqlFields {
				field := models.Field{
					Name: psqlField.Attname,
				}
				fields = append(fields, field)
			}

			table := models.Table{
				Name:   psqlTable.Relname,
				Fields: fields,
			}
			tables = append(tables, table)
		}

		schema := models.Schema{
			Name:           psqlSchema.Nspname,
			Tables:         tables,
			OrganizationId: organizationId,
			DataSourceName: dataSourceName,
		}

		err = self.schemaService.Create(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
