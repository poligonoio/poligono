package services

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	validatorv10 "github.com/go-playground/validator/v10"
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
	engineService        EngineService
	schemaService        SchemaService
	validate             *validatorv10.Validate
}

func NewDataSourceService(ctx context.Context, dataSourceCollection *mongo.Collection, infisicalService InfisicalService, engineService EngineService, schemaService SchemaService, validate *validatorv10.Validate) DataSourceService {
	return &DataSourceServiceImpl{
		ctx:                  ctx,
		dataSourceCollection: dataSourceCollection,
		infisicalService:     infisicalService,
		engineService:        engineService,
		schemaService:        schemaService,
		validate:             validate,
	}
}

func (self *DataSourceServiceImpl) GetByName(name string, organizationId string, includeSecret bool) (models.DataSource, error) {
	var dataSource models.DataSource

	filter := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "organization_id", Value: organizationId}}
	err := self.dataSourceCollection.FindOne(self.ctx, filter).Decode(&dataSource)
	if err != nil {
		return dataSource, err
	}

	if includeSecret {
		dataSource.Secret, err = self.infisicalService.GetSecret(dataSource.ID.Hex())
		if err != nil {
			return dataSource, err
		}
	}

	return dataSource, nil
}

func (self *DataSourceServiceImpl) GetById(id primitive.ObjectID, includeSecret bool) (models.DataSource, error) {
	var dataSource models.DataSource

	filter := bson.M{"_id": id}
	err := self.dataSourceCollection.FindOne(self.ctx, filter).Decode(&dataSource)
	if err != nil {
		return dataSource, err
	}

	if includeSecret {
		dataSource.Secret, err = self.infisicalService.GetSecret(dataSource.ID.Hex())
		if err != nil {
			return dataSource, err
		}
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

func (self *DataSourceServiceImpl) Create(dataSource models.DataSource) (models.DataSource, error) {
	// check if data source already exist
	filter := bson.D{bson.E{Key: "name", Value: dataSource.Name}, bson.E{Key: "organization_id", Value: dataSource.OrganizationId}}
	count, err := self.dataSourceCollection.CountDocuments(self.ctx, filter)
	if err != nil {
		return models.DataSource{}, err
	}

	if count > 0 {
		return models.DataSource{}, errors.New("Data source with that name already exists")
	}

	dataSource.ID = primitive.NewObjectID()

	_, err = self.dataSourceCollection.InsertOne(self.ctx, dataSource)
	if err != nil {
		return models.DataSource{}, err
	}

	insertedDataSource, err := self.GetByName(dataSource.Name, dataSource.OrganizationId, false)
	if err != nil {
		_, cleanUpErr := self.Delete(insertedDataSource.Name, insertedDataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return models.DataSource{}, err
	}

	// Create secret
	if err = self.infisicalService.CreateSecret(insertedDataSource.ID.Hex(), dataSource.Secret); err != nil {
		_, cleanUpErr := self.Delete(insertedDataSource.Name, insertedDataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return models.DataSource{}, err
	}

	// Create catalog
	if err = self.CreateCatalog(insertedDataSource.ID.Hex(), dataSource.Type, dataSource.Secret); err != nil {
		_, cleanUpErr := self.Delete(dataSource.Name, dataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return models.DataSource{}, err
	}

	// Sync schema
	err = self.Sync(insertedDataSource.ID, insertedDataSource.Type)
	if err != nil {
		_, cleanUpErr := self.Delete(insertedDataSource.Name, insertedDataSource.OrganizationId)
		if cleanUpErr != nil {
			logger.Error.Printf("Data source clean up failed: %v", cleanUpErr)
		}

		return models.DataSource{}, err
	}

	return insertedDataSource, nil
}

func (self *DataSourceServiceImpl) Update(name string, organizationId string, newDataSource models.UpdateRequestDataSourceBody) (models.DataSource, error) {
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

	filter := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "organization_id", Value: organizationId}}
	result, err := self.dataSourceCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return models.DataSource{}, err
	}

	if result.MatchedCount < 1 {
		return models.DataSource{}, fmt.Errorf("No data source found with name %s in the organization %s", name, organizationId)
	}

	updatedDataSource, err := self.GetByName(newDataSource.Name, organizationId, false)
	if err != nil {
		return models.DataSource{}, err
	}

	// update data source secret
	if newDataSource.Secret != "" {
		if err := self.infisicalService.UpdateSecret(updatedDataSource.ID.Hex(), newDataSource.Secret); err != nil {
			return models.DataSource{}, err
		}

		if err = self.engineService.RemoveCatalog(updatedDataSource.ID.Hex()); err != nil {
			return models.DataSource{}, err
		}

		if err = self.CreateCatalog(updatedDataSource.ID.Hex(), updatedDataSource.Type, newDataSource.Secret); err != nil {
			return models.DataSource{}, err
		}
	}

	return updatedDataSource, nil
}

func (self *DataSourceServiceImpl) Delete(name string, organizationId string) (models.DataSource, error) {
	// get document before deletion
	filter := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "organization_id", Value: organizationId}}
	dataSource, err := self.GetByName(name, organizationId, false)
	if err != nil {
		logger.Error.Printf("Error deleting data source document: %v", err)
	}

	result, err := self.dataSourceCollection.DeleteOne(self.ctx, filter)
	if err != nil {
		logger.Error.Printf("Error deleting data source document: %v", err)
		return dataSource, err
	}

	if result.DeletedCount < 1 {
		logger.Error.Printf("No Data source found with name %s in the organization %s", name, organizationId)
		return dataSource, err
	}

	if err = self.infisicalService.DeleteSecret(dataSource.ID.Hex()); err != nil {
		logger.Error.Printf("Error deleting secret: %v", err)
		return dataSource, err
	}

	if err = self.engineService.RemoveCatalog(dataSource.ID.Hex()); err != nil {
		logger.Error.Printf("Error deleting catalog: %v", err)
		return dataSource, err
	}

	if err = self.schemaService.Delete(dataSource.ID); err != nil {
		logger.Error.Printf("Error deleting schema document: %v", err)
		return dataSource, err
	}

	return dataSource, nil
}

func (self *DataSourceServiceImpl) GetDataSourceSchemas(id primitive.ObjectID) ([]models.Schema, error) {
	var schemas []models.Schema

	err := self.schemaService.GetAll(id, &schemas)
	if err != nil {
		return schemas, err
	}

	return schemas, nil
}

func (self *DataSourceServiceImpl) CreateCatalog(catalogName string, dataSourceType models.DataSourceType, secret string) error {
	var err error

	switch dataSourceType {
	case models.PostgreSQL:
		psql := NewPostgreSQLDataSourceDatabase(self.ctx, self.engineService, self.schemaService, self.validate)
		err = psql.CreateCatalog(catalogName, dataSourceType, secret)
	case models.MySQL:
		mysql := NewMySQLDataSourceDatabase(self.ctx, self.engineService, self.schemaService, self.validate)
		err = mysql.CreateCatalog(catalogName, dataSourceType, secret)
	case models.MariaDB:
		mariadb := NewMariaDBDataSourceDatabase(self.ctx, self.engineService, self.schemaService, self.validate)
		err = mariadb.CreateCatalog(catalogName, dataSourceType, secret)
	default:
		return errors.New("Invalid Data Source Type")
	}

	if err != nil {
		return err
	}

	return nil
}

func (self *DataSourceServiceImpl) Sync(id primitive.ObjectID, dataSourceType models.DataSourceType) error {
	var err error

	switch dataSourceType {
	case models.PostgreSQL:
		psql := NewPostgreSQLDataSourceDatabase(self.ctx, self.engineService, self.schemaService, self.validate)
		err = psql.Sync(id)
	case models.MySQL:
		mysql := NewMySQLDataSourceDatabase(self.ctx, self.engineService, self.schemaService, self.validate)
		err = mysql.Sync(id)
	case models.MariaDB:
		mariadb := NewMariaDBDataSourceDatabase(self.ctx, self.engineService, self.schemaService, self.validate)
		err = mariadb.Sync(id)
	default:
		return errors.New("Invalid Data Source Type")
	}

	if err != nil {
		return err
	}

	return nil
}
