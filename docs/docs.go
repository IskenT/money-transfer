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
        "/api/transfers": {
            "get": {
                "description": "Get a list of all transfers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transfers"
                ],
                "summary": "List all transfers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.TransferResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Transfer money from one user to another",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transfers"
                ],
                "summary": "Create a new money transfer",
                "parameters": [
                    {
                        "description": "Transfer details",
                        "name": "transfer",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.TransferRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.TransferResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/transfers/{id}": {
            "get": {
                "description": "Get transfer details by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transfers"
                ],
                "summary": "Get a specific transfer",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Transfer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.TransferResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users": {
            "get": {
                "description": "Get a list of all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "List all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.UserResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/{id}": {
            "get": {
                "description": "Get user details by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a specific user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.UserResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_IskenT_money-transfer_internal_infra_http_model.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "insufficient funds"
                }
            }
        },
        "github_com_IskenT_money-transfer_internal_infra_http_model.TransactionResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 1000
                },
                "amount_formatted": {
                    "type": "string",
                    "example": "$10.00"
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-04-10T12:34:56Z"
                },
                "note": {
                    "type": "string",
                    "example": "Transfer to Jane"
                },
                "payment_source": {
                    "type": "string",
                    "example": "TRANSFER"
                },
                "stan": {
                    "type": "string",
                    "example": "TRX1647881234567"
                },
                "state": {
                    "type": "string",
                    "example": "COMPLETED"
                },
                "transaction_type": {
                    "type": "string",
                    "example": "DEBIT"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-04-10T12:34:56Z"
                }
            }
        },
        "github_com_IskenT_money-transfer_internal_infra_http_model.TransferRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 1000
                },
                "from_user_id": {
                    "type": "string",
                    "example": "1"
                },
                "to_user_id": {
                    "type": "string",
                    "example": "2"
                }
            }
        },
        "github_com_IskenT_money-transfer_internal_infra_http_model.TransferResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 1000
                },
                "amount_formatted": {
                    "type": "string",
                    "example": "$10.00"
                },
                "completed_at": {
                    "type": "string",
                    "example": "2023-04-10T12:34:56Z"
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-04-10T12:34:56Z"
                },
                "credit_tx": {
                    "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.TransactionResponse"
                },
                "debit_tx": {
                    "$ref": "#/definitions/github_com_IskenT_money-transfer_internal_infra_http_model.TransactionResponse"
                },
                "from_user_id": {
                    "type": "string",
                    "example": "1"
                },
                "id": {
                    "type": "string",
                    "example": "TRF1647881234567"
                },
                "state": {
                    "type": "string",
                    "example": "COMPLETED"
                },
                "to_user_id": {
                    "type": "string",
                    "example": "2"
                }
            }
        },
        "github_com_IskenT_money-transfer_internal_infra_http_model.UserResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer",
                    "example": 10000
                },
                "balance_formatted": {
                    "type": "string",
                    "example": "$100.00"
                },
                "id": {
                    "type": "string",
                    "example": "1"
                },
                "name": {
                    "type": "string",
                    "example": "Mark"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
