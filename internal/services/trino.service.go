package services

import (
	"github.com/poligonoio/vega-core/internal/models"
)

type TrinoService interface {
	RemoveCatalog(catalogName string) error
	CreateCatalog(catalogName string, dataSourceType models.DataSourceType, secret string) error
	Query(query string, dest interface{}) error
	GetCatalogSchemas(catalogName string) ([]models.Schema, error)
	GetRawData(query string) ([]map[string]interface{}, error)
	GetCatalogName(name string, organizationId string) string
}
