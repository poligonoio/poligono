package services

import (
	"github.com/poligonoio/vega-core/internal/models"
)

type DataSourceService interface {
	GetByName(name string, organizationId string) (models.DataSource, error)
	GetAll(organizationId string) ([]models.DataSource, error)
	Create(dataSource models.DataSource) error
	Update(name string, organizationId string, newDataSource models.UpdateRequestDataSourceBody) error
	Delete(name string, organizationId string) error
	GetDataSourceSchemas(dataSourceName string, organizationId string) ([]models.Schema, error)
	Sync(dataSourceName string, organizationId string) error
}
