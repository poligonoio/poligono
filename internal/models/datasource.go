package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataSourceType string

const (
	PostgreSQL DataSourceType = "PostgreSQL"
	MySQL      DataSourceType = "MySQL"
)

type EngineType string

const (
	Trino     EngineType = "Trino"
	Starburst EngineType = "Starburst"
)

type DataSource struct {
	ID             primitive.ObjectID `json:"-" bson:"omitempty,_id" swaggerignore:"true"`
	Name           string             `json:"name"  bson:"name" validate:"required"`
	OrganizationId string             `json:"organization_id" bson:"organization_id" validate:"required" swaggerignore:"true"`
	CreatedBy      string             `json:"-" bson:"created_by" validate:"required"`
	Type           DataSourceType     `json:"type" bson:"type" validate:"required,oneof=PostgreSQL MySQL"`
	Secret         string             `json:"secret,omitempty" bson:"-" validate:"required"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at" validate:"required" swaggerignore:"true"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at" validate:"required" swaggerignore:"true"`
}

type UpdateRequestDataSourceBody struct {
	Name      string         `json:"name"  bson:"name"`
	Type      DataSourceType `json:"type" bson:"type" validate:"omitempty,oneof=PostgreSQL MySQL"`
	Secret    string         `json:"secret" bson:"-"`
	UpdatedAt time.Time      `json:"-" bson:"updated_at"`
}

type PostgreSQLSecret struct {
	Host     string `json:"hostname"`
	Port     int    `json:"port"`
	User     string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}

type MySQLSecret struct {
	Host     string `json:"hostname"`
	Port     int    `json:"port"`
	User     string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}
