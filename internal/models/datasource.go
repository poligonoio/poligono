package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataSourceType string

const (
	PostgreSQL DataSourceType = "PostgreSQL"
	MySQL      DataSourceType = "MySQL"
	MariaDB    DataSourceType = "MariaDB"
)

type EngineType string

const (
	Trino     EngineType = "Trino"
	Starburst EngineType = "Starburst"
)

type DataSource struct {
	ID             primitive.ObjectID `json:"-" bson:"_id" swaggerignore:"true"`
	Name           string             `json:"name"  bson:"name" validate:"required"`
	OrganizationId string             `json:"organization_id" bson:"organization_id" validate:"required" swaggerignore:"true"`
	CreatedBy      string             `json:"-" bson:"created_by" validate:"required"`
	Type           DataSourceType     `json:"type" bson:"type" validate:"required,oneof=PostgreSQL MySQL MariaDB"`
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
	Host     string `json:"hostname" validate:"required,hostname"`
	Port     int    `json:"port" validate:"required,number"`
	User     string `json:"username" validate:"required"`
	Database string `json:"database" validate:"required"`
	Password string `json:"password" validate:"required"`
	SSL      bool   `json:"ssl" validate:"boolean"`
}

type MySQLSecret struct {
	Host     string `json:"hostname" validate:"required,hostname"`
	Port     int    `json:"port" validate:"required,number"`
	User     string `json:"username" validate:"required"`
	Database string `json:"database" validate:"required"`
	Password string `json:"password" validate:"required"`
	SSL      bool   `json:"ssl" validate:"boolean"`
}
