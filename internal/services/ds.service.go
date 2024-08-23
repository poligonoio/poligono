package services

import (
	"github.com/poligonoio/vega-core/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataSourceService interface {
	GetByName(name string, organizationId string, includeSecret bool) (models.DataSource, error)
	GetById(id primitive.ObjectID, includeSecret bool) (models.DataSource, error)
	GetAll(organizationId string) ([]models.DataSource, error)
	Create(dataSource models.DataSource) (models.DataSource, error)
	Update(name string, organizationId string, newDataSource models.UpdateRequestDataSourceBody) (models.DataSource, error)
	Delete(name string, organizationId string) (models.DataSource, error)
	GetDataSourceSchemas(id primitive.ObjectID) ([]models.Schema, error)
	Sync(id primitive.ObjectID, dataSourceType models.DataSourceType) error
}
