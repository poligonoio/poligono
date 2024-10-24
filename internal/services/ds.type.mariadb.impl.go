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

type MariaDBDataSourceTypeImpl struct {
	ctx           context.Context
	engineService EngineService
	schemaService SchemaService
	validate      *validatorv10.Validate
}

func NewMariaDBDataSourceDatabase(ctx context.Context, engineService EngineService, schemaService SchemaService, validate *validatorv10.Validate) DataSourceTypeInter {
	return &MariaDBDataSourceTypeImpl{
		ctx:           ctx,
		engineService: engineService,
		schemaService: schemaService,
		validate:      validate,
	}
}

func (self *MariaDBDataSourceTypeImpl) Sync(dataSourceId primitive.ObjectID) error {
	catalogName := dataSourceId.Hex()
	var sqlSchemas []models.SQLSchema
	err := self.engineService.Query(fmt.Sprintf("SELECT schema_name AS name FROM data_source_%s.information_schema.schemata WHERE schema_name NOT IN ('sys', 'performance_schema')", catalogName), &sqlSchemas)
	if err != nil {
		return err
	}

	for _, sqlSchema := range sqlSchemas {
		var sqlTables []models.SQLSchema

		err := self.engineService.Query(fmt.Sprintf("SELECT table_name AS name FROM data_source_%s.information_schema.tables WHERE table_schema = '%s'", catalogName, sqlSchema.Name), &sqlTables)
		if err != nil {
			return err
		}

		var tables []models.Table
		for _, sqlTable := range sqlTables {
			var sqlFields []models.SQLField
			err := self.engineService.Query(fmt.Sprintf("SELECT column_name AS name FROM data_source_%s.information_schema.columns WHERE table_name = '%s'", catalogName, sqlTable.Name), &sqlFields)
			if err != nil {
				return err
			}

			var fields []models.Field
			for _, sqlField := range sqlFields {
				field := models.Field{
					Name: sqlField.Name,
				}
				fields = append(fields, field)
			}

			table := models.Table{
				Name:   sqlTable.Name,
				Fields: fields,
			}
			tables = append(tables, table)
		}

		schema := models.Schema{
			Name:   sqlSchema.Name,
			Tables: tables,
		}

		err = self.schemaService.Create(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *MariaDBDataSourceTypeImpl) CreateCatalog(catalogName string, dataSourceType models.DataSourceType, secret string) error {
	mysql := models.MySQLSecret{}
	if err := json.Unmarshal([]byte(secret), &mysql); err != nil {
		return err
	}

	if err := self.validate.Struct(mysql); err != nil {
		validationErr := err.(validatorv10.ValidationErrors)
		logger.Error.Println(fmt.Printf("One or more secret fields are invalid: %s\n", validationErr))
		return validationErr
	}

	var mysqlString string

	if mysql.SSL {
		mysqlString = fmt.Sprintf("jdbc:mariadb://%s:%s/%s?sslMode=REQUIRED", mysql.Host, strconv.Itoa(mysql.Port), mysql.Database)
	} else {
		mysqlString = fmt.Sprintf("jdbc:mariadb://%s:%s/%s", mysql.Host, strconv.Itoa(mysql.Port), mysql.Database)
	}

	query := fmt.Sprintf("CREATE CATALOG data_source_%s USING mariadb WITH (\"connection-url\" = '%s', \"connection-user\" = '%s', \"connection-password\" = '%s', \"case-insensitive-name-matching\" = 'true')", catalogName, mysqlString, mysql.User, mysql.Password)

	if err := self.engineService.Query(query, nil); err != nil {
		return err
	}

	return nil
}
