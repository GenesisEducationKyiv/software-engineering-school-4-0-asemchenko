// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/rate": {
            "get": {
                "description": "Request returns the current USD to UAH exchange rate using Monobank API",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rate"
                ],
                "summary": "Get the current USD to UAH exchange rate",
                "responses": {
                    "200": {
                        "description": "Current USD to UAH exchange rate",
                        "schema": {
                            "type": "number"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/subscribe": {
            "post": {
                "description": "Request adds a new email to receive USD to UAH exchange rate updates",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscription"
                ],
                "summary": "Subscribe to rate change notifications",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email address",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Email added",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Return if email already exists in the database",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "UAH currency application",
	Description:      "API for current USD-UAH exchange rate and for email-subscribing on the currency rate",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
