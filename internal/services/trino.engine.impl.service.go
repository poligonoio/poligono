package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/poligonoio/vega-core/pkg/logger"
	_ "github.com/trinodb/trino-go-client/trino"
)

type TrinoEngineServiceImpl struct {
	ctx context.Context
	db  *sql.DB
}

func NewTrinoEngineService(ctx context.Context) (EngineService, error) {
	database, err := sql.Open("trino", os.Getenv("QUERY_ENGINE_DSN"))
	if err != nil {
		return nil, err
	}

	return &TrinoEngineServiceImpl{
		ctx: ctx,
		db:  database,
	}, nil
}

func (self *TrinoEngineServiceImpl) RemoveCatalog(catalogName string) error {
	query := fmt.Sprintf("DROP CATALOG data_source_%s", catalogName)

	if err := self.Query(query, nil); err != nil {
		return err
	}

	return nil
}

func (self *TrinoEngineServiceImpl) Query(query string, dest interface{}) error {
	query = strings.ReplaceAll(query, ";", "")

	rows, err := self.db.Query(query)
	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
		return err
	}

	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
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

func (self *TrinoEngineServiceImpl) GetRawData(query string) ([]map[string]interface{}, error) {
	query = strings.ReplaceAll(query, ";", "")

	logger.Info.Printf("Query to be execute to get raw data from Data Source: %s\n", query)
	rows, err := self.db.Query(query)

	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
		return nil, err
	}

	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		logger.Error.Printf("Failed to get data from data source: %v\n", err)
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
