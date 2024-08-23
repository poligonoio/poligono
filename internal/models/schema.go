package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SQLSchema struct {
	Name string `json:"name"`
}

type Schema struct {
	ID           primitive.ObjectID `json:"-" bson:"_id" yaml:"-"`
	Name         string             `json:"name" bson:"name"`
	DataSourceId primitive.ObjectID `json:"data_source_id" bson:"data_source_id" yaml:"-"`
	Description  string             `json:"description" bson:"description" yaml:"description,omitempty"`
	Tables       []Table            `json:"tables" bson:"tables"`
}

type UpdateSchema struct {
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description" yaml:"description,omitempty"`
	Tables      []Table `json:"tables" bson:"tables"`
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
