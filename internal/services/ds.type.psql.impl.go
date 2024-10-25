package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/poligonoio/vega-core/internal/models"
	"github.com/poligonoio/vega-core/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostgresSQLDataSourceTypeImpl struct {
	ctx           context.Context
	engineService EngineService
	schemaService SchemaService
	validate      *validatorv10.Validate
}

func NewPostgreSQLDataSourceDatabase(ctx context.Context, engineService EngineService, schemaService SchemaService, validate *validatorv10.Validate) DataSourceTypeInter {
	return &PostgresSQLDataSourceTypeImpl{
		ctx:           ctx,
		engineService: engineService,
		schemaService: schemaService,
		validate:      validate,
	}
}

func (self *PostgresSQLDataSourceTypeImpl) Sync(dataSourceId primitive.ObjectID) error {
	catalogName := dataSourceId.Hex()
	var sqlSchemas []models.SQLSchema
	err := self.engineService.Query(fmt.Sprintf("SELECT nspname AS name FROM data_source_%s.pg_catalog.pg_namespace WHERE nspname NOT IN ('pg_toast', 'pg_catalog', 'public', 'information_schema')", catalogName), &sqlSchemas)
	if err != nil {
		return err
	}

	for _, sqlSchema := range sqlSchemas {
		var psqlTables []models.SQLTable
		err := self.engineService.Query(fmt.Sprintf("SELECT relname AS name FROM data_source_%s.pg_catalog.pg_class c INNER JOIN data_source_%s.pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = '%s'", catalogName, catalogName, sqlSchema.Name), &psqlTables)
		if err != nil {
			return err
		}

		var tables []models.Table
		for _, psqlTable := range psqlTables {

			var psqlFields []models.SQLField
			err := self.engineService.Query(fmt.Sprintf("SELECT attname AS name FROM data_source_%s.pg_catalog.pg_attribute a INNER JOIN data_source_%s.pg_catalog.pg_class c ON a.attrelid = c.oid WHERE c.relname = '%s'", catalogName, catalogName, psqlTable.Name), &psqlFields)
			if err != nil {
				return err
			}

			var fields []models.Field
			for _, psqlField := range psqlFields {
				field := models.Field{
					Name: psqlField.Name,
				}
				fields = append(fields, field)
			}

			table := models.Table{
				Name:   psqlTable.Name,
				Fields: fields,
			}
			tables = append(tables, table)
		}

		schema := models.Schema{
			Name:         sqlSchema.Name,
			DataSourceId: dataSourceId,
			Tables:       tables,
		}

		err = self.schemaService.Create(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *PostgresSQLDataSourceTypeImpl) CreateCatalog(catalogName string, dataSourceType models.DataSourceType, secret string) error {
	psql := models.PostgreSQLSecret{}
	if err := json.Unmarshal([]byte(secret), &psql); err != nil {
		return err
	}

	if err := self.validate.Struct(psql); err != nil {
		validationErr := err.(validatorv10.ValidationErrors)
		logger.Error.Println(fmt.Printf("One or more secret fields are invalid: %s\n", validationErr))
		return validationErr
	}

	var psqlString string

	if psql.SSL {
		psqlString = fmt.Sprintf("jdbc:postgresql://%s:%s/%s?sslmode=require", psql.Host, strconv.Itoa(psql.Port), psql.Database)
	} else {
		psqlString = fmt.Sprintf("jdbc:postgresql://%s:%s/%s", psql.Host, strconv.Itoa(psql.Port), psql.Database)
	}

	query := fmt.Sprintf("CREATE CATALOG data_source_%s USING postgresql WITH (\"connection-url\" = '%s', \"connection-user\" = '%s', \"connection-password\" = '%s', \"case-insensitive-name-matching\" = 'true', \"postgresql.include-system-tables\" = 'true')", catalogName, psqlString, psql.User, psql.Password)

	if err := self.engineService.Query(query, nil); err != nil {
		return err
	}

	return nil
}
