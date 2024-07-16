// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://swagger.io/terms/",
        "contact": {
            "name": "Poligono Support",
            "url": "https://www.swagger.io/support",
            "email": "dev@poligono.xyz"
        },
        "license": {
            "name": "GNU Affero General Public License version 3",
            "url": "https://www.gnu.org/licenses/agpl-3.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/datasources": {
            "post": {
                "description": "Create a new data source configuration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data_source"
                ],
                "summary": "Add Data Source",
                "parameters": [
                    {
                        "description": "Data Source",
                        "name": "Data_Source",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.DataSource"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPSuccess"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            }
        },
        "/datasources/all": {
            "get": {
                "description": "Retrieve all data sources associated with the specified organization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data_source"
                ],
                "summary": "List Data Sources",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.DataSource"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            }
        },
        "/datasources/{name}": {
            "get": {
                "description": "Retrieve data source configuration with the specified name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data_source"
                ],
                "summary": "Retrieve Data Source by Name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Data Source Name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.DataSource"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update the configuration of the specified data source",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data_source"
                ],
                "summary": "Modify Data Source",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Data Source Name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Data Source",
                        "name": "Data_Source",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdateRequestDataSourceBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPSuccess"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Permanently deletes the specified data source",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data_source"
                ],
                "summary": "Remove Data Source",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Data Source Name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPSuccess"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            }
        },
        "/prompt": {
            "post": {
                "description": "Create an SQL query based on a natural language prompt",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "query"
                ],
                "summary": "Generate SQL Query",
                "parameters": [
                    {
                        "description": "Prompt Object",
                        "name": "Prompt",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Prompt"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Activity"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Activity": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "additionalProperties": true
                    }
                },
                "data_source_name": {
                    "type": "string"
                },
                "organization_id": {
                    "type": "string"
                },
                "prompt": {
                    "type": "string"
                },
                "query": {
                    "type": "string"
                }
            }
        },
        "models.DataSource": {
            "type": "object",
            "required": [
                "name",
                "secret",
                "type"
            ],
            "properties": {
                "name": {
                    "type": "string"
                },
                "secret": {
                    "type": "string"
                },
                "type": {
                    "enum": [
                        "PostgreSQL"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.DataSourceType"
                        }
                    ]
                }
            }
        },
        "models.DataSourceType": {
            "type": "string",
            "enum": [
                "PostgreSQL"
            ],
            "x-enum-varnames": [
                "PostgreSQL"
            ]
        },
        "models.HTTPError": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "models.HTTPSuccess": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Prompt": {
            "type": "object",
            "properties": {
                "data_source_name": {
                    "type": "string"
                },
                "prompt": {
                    "type": "string"
                }
            }
        },
        "models.UpdateRequestDataSourceBody": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "secret": {
                    "type": "string"
                },
                "type": {
                    "enum": [
                        "PostgreSQL"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.DataSourceType"
                        }
                    ]
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Poligono API",
	Description:      "Democratizing data access through plain English.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
