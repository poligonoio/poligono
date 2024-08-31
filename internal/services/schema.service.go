package services

import (
	"github.com/poligonoio/vega-core/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SchemaService interface {
	Create(schema models.Schema) error
	GetAll(dataSourceId primitive.ObjectID, schemas *[]models.Schema) error
	Delete(dataSourceId primitive.ObjectID) error
}
