package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataSourceType string

const (
	PostgreSQL DataSourceType = "PostgreSQL"
)

type DataSource struct {
	ID             primitive.ObjectID `json:"-" bson:"omitempty,_id" swaggerignore:"true"`
	Name           string             `json:"name"  bson:"name" validate:"required"`
	OrganizationId string             `json:"organization_id" bson:"organization_id" validate:"required" swaggerignore:"true"`
	CreatedBy      string             `json:"-" bson:"created_by" validate:"required"`
	Type           DataSourceType     `json:"type" bson:"type" validate:"required,oneof=PostgreSQL"`
	Secret         string             `json:"secret,omitempty" bson:"-" validate:"required"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at" validate:"required" swaggerignore:"true"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at" validate:"required" swaggerignore:"true"`
}

type UpdateRequestDataSourceBody struct {
	Name      string         `json:"name"  bson:"name"`
	Type      DataSourceType `json:"type" bson:"type" validate:"omitempty,oneof=PostgreSQL"`
	Secret    string         `json:"secret" bson:"-"`
	UpdatedAt time.Time      `json:"-" bson:"updated_at"`
}

type PostgreSQLObject struct {
	Host     string `json:"hostname"`
	Port     int    `json:"port"`
	User     string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}

func GetCatalogQuery(catalogName string, dataSourceType DataSourceType, secret string) (string, error) {
	var err error
	switch dataSourceType {
	case PostgreSQL:
		psql := PostgreSQLObject{}
		if err = json.Unmarshal([]byte(secret), &psql); err != nil {
			return "", err
		}

		var psqlString string

		if psql.SSL {
			psqlString = fmt.Sprintf("jdbc:postgresql://%s:%s/%s?sslmode=require", psql.Host, strconv.Itoa(psql.Port), psql.Database)
		} else {
			psqlString = fmt.Sprintf("jdbc:postgresql://%s:%s/%s", psql.Host, strconv.Itoa(psql.Port), psql.Database)
		}

		return fmt.Sprintf("CREATE CATALOG %s USING postgresql WITH (\"connection-url\" = '%s', \"connection-user\" = '%s', \"connection-password\" = '%s', \"case-insensitive-name-matching\" = 'true', \"postgresql.include-system-tables\" = 'true')", catalogName, psqlString, psql.User, psql.Password), nil
	default:
		return "", errors.New("invalid data source type")
	}
}
