package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/poligonoio/vega-core/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchemaServiceImpl struct {
	ctx              context.Context
	schemaCollection *mongo.Collection
}

func NewSchemaService(ctx context.Context, schemaCollection *mongo.Collection) SchemaService {
	return &SchemaServiceImpl{
		ctx:              ctx,
		schemaCollection: schemaCollection,
	}
}

func (self *SchemaServiceImpl) Create(schema models.Schema) error {
	query := bson.D{bson.E{Key: "name", Value: schema.Name}, bson.E{Key: "organization_id", Value: schema.OrganizationId}}
	count, err := self.schemaCollection.CountDocuments(self.ctx, query)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("Data source with that name already exists")
	}

	_, err = self.schemaCollection.InsertOne(self.ctx, schema)
	if err != nil {
		return err
	}

	return nil
}

func (self *SchemaServiceImpl) GetAll(dataSourceName string, organizationId string, schemas *[]models.Schema) error {
	filter := bson.D{bson.E{Key: "data_source_name", Value: dataSourceName}, bson.E{Key: "organization_id", Value: organizationId}}
	cursor, err := self.schemaCollection.Find(self.ctx, filter)
	if err != nil {
		return err
	}

	if err = cursor.All(self.ctx, schemas); err != nil {
		return err
	}

	return nil
}

func (self *SchemaServiceImpl) UpdateDataSourceName(oldDataSourceName string, newDataSourceName string, organizationId string) error {
	query := bson.D{bson.E{Key: "data_source_name", Value: oldDataSourceName}, bson.E{Key: "organization_id", Value: organizationId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "data_source_name", Value: newDataSourceName}}}}

	result, err := self.schemaCollection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return err
	}

	if result.MatchedCount < 1 {
		return fmt.Errorf("No schema found with data source name  %s in the organization %s", oldDataSourceName, organizationId)
	}

	return nil
}

func (self *SchemaServiceImpl) Delete(dataSourceName string, organizationId string) error {
	filter := bson.D{bson.E{Key: "data_source_name", Value: dataSourceName}, bson.E{Key: "organization_id", Value: organizationId}}
	result, err := self.schemaCollection.DeleteOne(self.ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount < 1 {
		return fmt.Errorf("No Schema found for Data Source %s in the organization %s", dataSourceName, organizationId)
	}

	return nil
}
