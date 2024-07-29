package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SQLSchema struct {
	Name string `json:"name"`
}

type Schemas struct {
	Schemas []Schema
}

type Schema struct {
	ID             primitive.ObjectID `json:"-" bson:"omitempty,_id" yaml:"-"`
	Name           string             `json:"name" bson:"name"`
	OrganizationId string             `json:"organization_id" bson:"organization_id" yaml:"-"`
	DataSourceName string             `json:"data_source_name" bson:"data_source_name" yaml:"-"`
	Description    string             `json:"description" bson:"description" yaml:"description,omitempty"`
	Tables         []Table            `json:"tables" bson:"tables"`
}

type UpdateSchema struct {
	Name           string  `json:"name" bson:"name"`
	DataSourceName string  `json:"data_source_name" bson:"data_source_name" yaml:"-"`
	Description    string  `json:"description" bson:"description" yaml:"description,omitempty"`
	Tables         []Table `json:"tables" bson:"tables"`
}

type SQLTable struct {
	Name string `json:"name"`
}

type Table struct {
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description" yaml:"description,omitempty"`
	Fields      []Field `json:"fields" bson:"fields"`
}

type SQLField struct {
	Name string `json:"name"`
}

type Field struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description" yaml:"description,omitempty"`
}
