// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/api/actor/films/": {
            "get": {
                "description": "Get all the films which a certain actor was shot in",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "actor"
                ],
                "summary": "Films actor took a part in",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Actor filter",
                        "name": "actor",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Film"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/add/": {
            "post": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Based on the body of POST request add film to the DB",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "film"
                ],
                "summary": "Add film to the DB",
                "parameters": [
                    {
                        "description": "film model",
                        "name": "film",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.Film"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.Film"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/all/": {
            "get": {
                "description": "Get every available film from the DB",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "films"
                ],
                "summary": "Get all the films from the DB",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Film"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/delete/": {
            "delete": {
                "security": [
                    {
                        "JWT": []
                    }
                ],
                "description": "Delete film from both cache and main repositories based on the user's provided filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "film"
                ],
                "summary": "Delete film",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Film title",
                        "name": "title",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Film genre",
                        "name": "genre",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Film release date",
                        "name": "released_year",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DeleteRequest"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/get/": {
            "get": {
                "description": "Get all films with the specified name.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "films"
                ],
                "summary": "Get films with the specified name (there can be more than 1 film with the same name)",
                "parameters": [
                    {
                        "type": "string",
                        "description": "film title",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Film"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/producer/films/": {
            "get": {
                "description": "Get all the films that'd been produced by a specified producer.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "producer"
                ],
                "summary": "Films, which were produced by the specified person",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Producer filter",
                        "name": "producer",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/api.Film"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/signin/": {
            "post": {
                "description": "Process authentication",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Sign in the service",
                "parameters": [
                    {
                        "description": "Sign in model",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/signup/": {
            "post": {
                "description": "Create an account in the service",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Sign up to the service",
                "parameters": [
                    {
                        "description": "user model",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "405": {
                        "description": "Method Not Allowed"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "api.DeleteRequest": {
            "description": "request, according to which a film is going to be deleted from the database. It has 3 filters: Genre, Title and ReleasedYear",
            "type": "object",
            "properties": {
                "genre": {
                    "description": "A genre of the show to be deleted",
                    "type": "string"
                },
                "released_year": {
                    "description": "A year of release of a show to be deleted",
                    "type": "integer"
                },
                "title": {
                    "description": "A title of the show to be deleted",
                    "type": "string"
                }
            }
        },
        "api.Film": {
            "description": "model's being operated on",
            "type": "object",
            "properties": {
                "crew": {
                    "description": "Crew that took a part in production",
                    "allOf": [
                        {
                            "$ref": "#/definitions/general.Crew"
                        }
                    ]
                },
                "genre": {
                    "description": "A genre of the show",
                    "type": "string"
                },
                "released_year": {
                    "description": "Year the show was released in",
                    "type": "integer"
                },
                "revenue": {
                    "description": "Revenue which was received by the show",
                    "type": "number"
                },
                "title": {
                    "description": "Title of the show",
                    "type": "string"
                }
            }
        },
        "auth.Role": {
            "description": "Value that represents user's right in the service.",
            "type": "string",
            "enum": [
                "admin",
                "guest"
            ],
            "x-enum-comments": {
                "ADMIN": "This role provides full access to the API",
                "GUEST": "This role provides restricted access to the API - client gets only methods for observation."
            },
            "x-enum-varnames": [
                "ADMIN",
                "GUEST"
            ]
        },
        "auth.SignInRequest": {
            "description": "Model that represents user's input to sign in the service provided",
            "type": "object",
            "properties": {
                "password": {
                    "description": "Password of the user",
                    "type": "string"
                },
                "username": {
                    "description": "Username of the user",
                    "type": "string"
                }
            }
        },
        "auth.User": {
            "description": "Model that represents user's model, also the content of a body in the request to be signed up.",
            "type": "object",
            "properties": {
                "id": {
                    "description": "id of the user that's server-side generated",
                    "type": "string"
                },
                "password": {
                    "description": "password of the client that is hashed on the server side",
                    "type": "string"
                },
                "role": {
                    "description": "role of the user. Either guest or admin.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/auth.Role"
                        }
                    ]
                },
                "username": {
                    "description": "username of the client. Must be unique",
                    "type": "string"
                }
            }
        },
        "general.Actor": {
            "description": "Actor that took a part in the show",
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "description": "Actor's role in the show",
                    "type": "string"
                }
            }
        },
        "general.Crew": {
            "description": "The model of the crew which took part in shooting the show",
            "type": "object",
            "properties": {
                "actors": {
                    "description": "All actors took a part in the show",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/general.Actor"
                    }
                },
                "producers": {
                    "description": "All producers that took a part in the show",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/general.Producer"
                    }
                }
            }
        },
        "general.Producer": {
            "description": "producer of the show",
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWT": {
            "type": "apiKey",
            "name": "jwt_token",
            "in": "cookie"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Films Aggregator",
	Description:      "API Server for Films Aggregator application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
