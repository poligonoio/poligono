package services

import (
	"context"
	"fmt"

	"github.com/poligonoio/vega-core/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	schema.ID = primitive.NewObjectID()

	_, err := self.schemaCollection.InsertOne(self.ctx, schema)
	if err != nil {
		return err
	}

	return nil
}

func (self *SchemaServiceImpl) GetAll(dataSourceId primitive.ObjectID, schemas *[]models.Schema) error {
	filter := bson.D{bson.E{Key: "data_source_id", Value: dataSourceId}}
	cursor, err := self.schemaCollection.Find(self.ctx, filter)
	if err != nil {
		return err
	}

	if err = cursor.All(self.ctx, schemas); err != nil {
		return err
	}

	return nil
}

func (self *SchemaServiceImpl) Delete(dataSourceId primitive.ObjectID) error {
	filter := bson.D{bson.E{Key: "data_source_id", Value: dataSourceId}}
	result, err := self.schemaCollection.DeleteOne(self.ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount < 1 {
		return fmt.Errorf("No Schema found for Data Source ID: %s", dataSourceId.Hex())
	}

	return nil
}
