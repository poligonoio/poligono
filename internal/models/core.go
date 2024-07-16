package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Prompt struct {
	Text           string `json:"prompt"`
	DataSourceName string `json:"data_source_name"`
}

type GeminiQueryResult struct {
	QueryMarkdown string `json:"query"`
}

type Activity struct {
	ID             primitive.ObjectID       `json:"-" bson:"_id,omitempty"`
	Prompt         string                   `json:"prompt" bson:"prompt"`
	Query          string                   `json:"query" bson:"query"`
	Data           []map[string]interface{} `json:"data" bson:"-"`
	MergedPrompt   string                   `json:"-" bson:"data_source_name"`
	DataSourceId   primitive.ObjectID       `json:"-" bson:"data_source_id,omitempty"`
	DataSourceName string                   `json:"data_source_name" bson:"data_source_name"`
	OrganizationId string                   `json:"organization_id" bson:"organization_id"`
	UserId         string                   `json:"-" bson:"user_id"`
}

type HTTPError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

type HTTPSuccess struct {
	Message string `json:"message"`
}
