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
        "/": {
            "get": {
                "description": "Инициализирует логирование, подключается к базе данных, настраивает маршруты и запускает HTTP-сервер.",
                "tags": [
                    "main"
                ],
                "summary": "Основная точка входа приложения.",
                "responses": {
                    "200": {
                        "description": "Сервер успешно запущен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Не удалось запустить сервер",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/logger/error": {
            "post": {
                "description": "Записывает сообщение об ошибке в лог",
                "tags": [
                    "logger"
                ],
                "summary": "Запись сообщения об ошибке",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Сообщение для логирования",
                        "name": "msg",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Message logged successfully"
                    },
                    "500": {
                        "description": "Failed to log message"
                    }
                }
            }
        },
        "/logger/info": {
            "post": {
                "description": "Записывает информационное сообщение в лог",
                "tags": [
                    "logger"
                ],
                "summary": "Запись информационного сообщения",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Сообщение для логирования",
                        "name": "msg",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Message logged successfully"
                    },
                    "500": {
                        "description": "Failed to log message"
                    }
                }
            }
        },
        "/logger/init": {
            "post": {
                "description": "Открывает файл для записи логов и инициализирует логгеры для разных уровней",
                "tags": [
                    "logger"
                ],
                "summary": "Инициализация логгирования",
                "responses": {
                    "200": {
                        "description": "Logging initialized successfully"
                    },
                    "500": {
                        "description": "Failed to initialize logging"
                    }
                }
            }
        },
        "/logger/warning": {
            "post": {
                "description": "Записывает предупреждающее сообщение в лог",
                "tags": [
                    "logger"
                ],
                "summary": "Запись предупреждающего сообщения",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Сообщение для логирования",
                        "name": "msg",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Message logged successfully"
                    },
                    "500": {
                        "description": "Failed to log message"
                    }
                }
            }
        },
        "/user": {
            "put": {
                "description": "Updates the user details in the database.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "description": "Updated user details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dataBase.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Проверяет, что имя пользователя, номер телефона, пароль и email заполнены и соответствуют требованиям.",
                "summary": "Валидация данных пользователя для создания",
                "parameters": [
                    {
                        "description": "Данные пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Данные пользователя валидны",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Ошибка валидации данных",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/{id}": {
            "get": {
                "description": "Retrieves a user from the database by their ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get user by ID",
                "parameters": [
                    {
                        "type": "integer",
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
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dataBase.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Проверяет, что имя пользователя, номер телефона, пароль и email соответствуют требованиям.",
                "summary": "Валидация данных пользователя для обновления",
                "parameters": [
                    {
                        "description": "Данные пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Данные пользователя валидны",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Ошибка валидации данных",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a user from the database by their ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dataBase.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Retrieve a list of users with optional filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get all users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of users per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by email",
                        "name": "email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by phone",
                        "name": "phone",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of users",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new user with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Retrieve a user using their ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User data",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Invalid ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a user's details by their ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update a user's information",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "User updated successfully"
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove a user using their ID",
                "tags": [
                    "users"
                ],
                "summary": "Delete a user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "User deleted successfully"
                    },
                    "400": {
                        "description": "Invalid ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dataBase.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.User": {
            "description": "Модель пользователя с основными полями.",
            "type": "object",
            "properties": {
                "email": {
                    "description": "@Description Электронная почта пользователя\n@Example \"johndoe@example.com\"",
                    "type": "string"
                },
                "fromDateCreate": {
                    "description": "@Description Дата и время создания пользователя\n@Example \"2024-10-19T12:00:00Z\"",
                    "type": "string"
                },
                "fromDateUpdate": {
                    "description": "@Description Дата и время последнего обновления пользователя\n@Example \"2024-10-19T12:00:00Z\"",
                    "type": "string"
                },
                "id": {
                    "description": "@Description Уникальный идентификатор пользователя\n@Example 1",
                    "type": "integer"
                },
                "isBanned": {
                    "description": "@Description Флаг, указывающий, заблокирован ли пользователь\n@Example false",
                    "type": "boolean"
                },
                "isDeleted": {
                    "description": "@Description Флаг, указывающий, удален ли пользователь\n@Example false",
                    "type": "boolean"
                },
                "name": {
                    "description": "@Description Имя пользователя\n@Example \"John Doe\"",
                    "type": "string"
                },
                "password": {
                    "description": "@Description Пароль пользователя (сохраняется в зашифрованном виде)\n@Example \"password123\"",
                    "type": "string"
                },
                "phone": {
                    "description": "@Description Номер телефона пользователя\n@Example \"+1234567890\"",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Cloud Application API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}