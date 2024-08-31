package services

import (
	"github.com/poligonoio/vega-core/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataSourceTypeInter interface {
	// GetDataSourceSchemas(dataSourceName string, organizationId string) ([]models.Schema, error)
	Sync(dataSourceId primitive.ObjectID) error
	CreateCatalog(catalogName string, dataSourceType models.DataSourceType, secret string) error
}
