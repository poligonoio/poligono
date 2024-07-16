package services

import (
	"github.com/poligonoio/vega-core/internal/models"
)

type SchemaService interface {
	Create(schema models.Schema) error
	GetAll(dataSourceName string, organizationId string, schemas *[]models.Schema) error
	UpdateDataSourceName(dataSourceName string, newDataSourceName string, organizationId string) error
	Delete(dataSourceName string, organizationId string) error
}
