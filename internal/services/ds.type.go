package services

import "github.com/poligonoio/vega-core/internal/models"

type DataSourceTypeInter interface {
	// GetDataSourceSchemas(dataSourceName string, organizationId string) ([]models.Schema, error)
	Sync(dataSourceName string, organizationId string) error
	CreateCatalog(catalogName string, dataSourceType models.DataSourceType, secret string) error
}
