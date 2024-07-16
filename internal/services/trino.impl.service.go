package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/poligonoio/vega-core/internal/models"
	"github.com/poligonoio/vega-core/pkg/logger"
	_ "github.com/trinodb/trino-go-client/trino"
)

type TrinoServiceImpl struct {
	ctx context.Context
	db  *sql.DB
}

func NewTrinoService(ctx context.Context) (TrinoService, error) {
	database, err := sql.Open("trino", os.Getenv("TRINO_DSN"))
	if err != nil {
		return nil, err
	}

	return &TrinoServiceImpl{
		ctx: ctx,
		db:  database,
	}, nil
}

func (self *TrinoServiceImpl) RemoveCatalog(catalogName string) error {
	query := fmt.Sprintf("DROP CATALOG %s", catalogName)

	if err := self.Query(query, nil); err != nil {
		return err
	}

	return nil
}

func (self *TrinoServiceImpl) CreateCatalog(catalogName string, dataSourceType models.DataSourceType, secret string) error {
	query, err := models.GetCatalogQuery(catalogName, dataSourceType, secret)
	if err != nil {
		return err
	}

	if err = self.Query(query, nil); err != nil {
		return err
	}

	return nil
}

func (self *TrinoServiceImpl) Query(query string, dest interface{}) error {
	query = strings.ReplaceAll(query, ";", "")

	rows, err := self.db.Query(query)
	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
		logger.Error.Printf("Query: \n\n%s\n\n", query)
		return err
	}

	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
		logger.Error.Printf("Query: \n\n%s\n\n", query)
		return err
	}

	colNum := len(columns)

	var results []map[string]interface{}

	for rows.Next() {
		// Prepare to read row using Scan
		r := make([]interface{}, colNum)
		for i := range r {
			r[i] = &r[i]
		}

		// Read rows using Scan
		err = rows.Scan(r...)
		if err != nil {
			logger.Error.Printf("Failed to get data from data source: %v\n", err)
			logger.Error.Printf("Query: \n\n%s\n\n", query)
			return err
		}

		// Create a row map to store row's data
		var row = map[string]interface{}{}
		for i := range r {
			row[columns[i]] = r[i]
		}

		// Append to the final results slice
		results = append(results, row)
	}

	if dest != nil {
		resultsJson, err := json.Marshal(results)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(resultsJson, dest); err != nil {
			return err
		}
	}

	return nil
}

func (self *TrinoServiceImpl) GetRawData(query string) ([]map[string]interface{}, error) {
	query = strings.ReplaceAll(query, ";", "")

	logger.Info.Printf("Query to be execute to get raw data from Data Source: %s\n", query)
	rows, err := self.db.Query(query)

	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
		logger.Error.Printf("Query: \n\n%s\n\n", query)
		return nil, err
	}

	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
		logger.Error.Printf("Query: \n\n%s\n\n", query)
		return nil, err
	}

	colNum := len(columns)

	var results []map[string]interface{}

	for rows.Next() {
		// Prepare to read row using Scan
		r := make([]interface{}, colNum)
		for i := range r {
			r[i] = &r[i]
		}

		// Read rows using Scan
		err = rows.Scan(r...)
		if err != nil {
			logger.Error.Printf("Failed to get data from data source: %v\n", err)
			logger.Error.Printf("Query: \n\n%s\n\n", query)
			return nil, err
		}

		// Create a row map to store row's data
		var row = map[string]interface{}{}
		for i := range r {
			row[columns[i]] = r[i]
		}

		// Append to the final results slice
		results = append(results, row)
	}

	return results, nil
}

func (self *TrinoServiceImpl) GetCatalogSchemas(catalogName string) ([]models.Schema, error) {
	var schemas []models.Schema

	var psqlSchemas []models.PSQLSchema
	err := self.Query(fmt.Sprintf("SELECT nspname FROM %s.pg_catalog.pg_namespace WHERE nspname NOT IN ('pg_toast', 'pg_catalog', 'public', 'information_schema')", catalogName), &psqlSchemas)
	if err != nil {
		return nil, err
	}

	for _, psqlSchema := range psqlSchemas {
		var psqlTables []models.PSQLTable
		err := self.Query(fmt.Sprintf("SELECT relname, relnamespace FROM %s.pg_catalog.pg_class c INNER JOIN %s.pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = '%s'", catalogName, catalogName, psqlSchema.Nspname), &psqlTables)
		if err != nil {
			return nil, err
		}

		var tables []models.Table
		for _, psqlTable := range psqlTables {

			var psqlFields []models.PSQLField
			err := self.Query(fmt.Sprintf("SELECT * FROM %s.pg_catalog.pg_attribute a INNER JOIN %s.pg_catalog.pg_class c ON a.attrelid = c.oid WHERE c.relname = '%s'", catalogName, catalogName, psqlTable.Relname), &psqlFields)
			if err != nil {
				return nil, err
			}

			var fields []models.Field
			for _, psqlField := range psqlFields {
				field := models.Field{
					Name: psqlField.Attname,
				}
				fields = append(fields, field)
			}

			table := models.Table{
				Name:   psqlTable.Relname,
				Fields: fields,
			}
			tables = append(tables, table)
		}

		schema := models.Schema{
			Name:   psqlSchema.Nspname,
			Tables: tables,
		}

		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func (self *TrinoServiceImpl) GetCatalogName(name string, organizationId string) string {
	if name != "" && organizationId != "" {
		return strings.ToLower(fmt.Sprintf("%s_%s", name, organizationId))
	}

	return ""
}
