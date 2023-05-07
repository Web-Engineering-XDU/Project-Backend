// Code generated by swaggo/swag. DO NOT EDIT.

package swaggo

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/agent": {
            "get": {
                "description": "get agents",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "summary": "List agents",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "agent id. Don't include in request if you don't specify a id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of agents in a page",
                        "name": "number",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "page sequence number",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swaggo.GetAgentListResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Create new agents. Also support application/json",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "summary": "New agents",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "enable the agent",
                        "name": "enable",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "agent type id",
                        "name": "typeId",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name of the agent",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "description",
                        "name": "description",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "whether keep the event forever",
                        "name": "eventForever",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "event max age in nanosecond count",
                        "name": "eventMaxAge",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "props used by specific agent type in json",
                        "name": "propJsonStr",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swaggo.NewAgentResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Update agents with specific id. Also support application/json",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "summary": "Update agents",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "agent id",
                        "name": "id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "enable the agent",
                        "name": "enable",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name of the agent",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "description",
                        "name": "description",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "whether keep the event forever",
                        "name": "eventForever",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "event max age in timestamp",
                        "name": "eventMaxAge",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "props used by specific agent type in json",
                        "name": "propJsonStr",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swaggo.StateInfo"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete agents",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "agents"
                ],
                "summary": "Delete agents",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "agent id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swaggo.StateInfo"
                        }
                    }
                }
            }
        },
        "/agent-relation": {
            "get": {
                "description": "get all relations",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "relations"
                ],
                "summary": "List Relations",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swaggo.GetRelationsResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Set relations of an agent with specific id. Also support x-www-form-urlencoded",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "relations"
                ],
                "summary": "Set Agent Relations",
                "parameters": [
                    {
                        "description": "all relations of an agent",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app_service.setAgentRelationParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swaggo.GetEventListResponse"
                        }
                    }
                }
            }
        },
        "/event": {
            "get": {
                "description": "get events",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "List Events",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "src agent id. Don't include in request if you don't specify it",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of events in a page",
                        "name": "number",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "page sequence number",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swaggo.GetEventListResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app_service.setAgentRelationParams": {
            "type": "object",
            "properties": {
                "agentId": {
                    "type": "integer"
                },
                "dsts": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "srcs": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "github_com_Web-Engineering-XDU_Project-Backend_app_models.AgentDetail": {
            "type": "object",
            "properties": {
                "allowInput": {
                    "type": "boolean",
                    "example": false
                },
                "allowOutput": {
                    "type": "boolean",
                    "example": true
                },
                "createAt": {
                    "type": "string",
                    "example": "2023-04-11T05:07:53+08:00"
                },
                "description": {
                    "type": "string",
                    "example": "I'm a schedule agent"
                },
                "enable": {
                    "type": "boolean",
                    "example": true
                },
                "eventForever": {
                    "type": "boolean",
                    "example": false
                },
                "eventMaxAge": {
                    "type": "integer",
                    "example": 0
                },
                "id": {
                    "type": "integer",
                    "example": 8848
                },
                "name": {
                    "type": "string",
                    "example": "5s timer"
                },
                "propJsonStr": {
                    "type": "string",
                    "example": "{\"cron\":\"*/5 * * * * *\"}"
                },
                "typeId": {
                    "type": "integer",
                    "example": 1
                },
                "typeName": {
                    "type": "string",
                    "example": "Schedule Agent"
                }
            }
        },
        "github_com_Web-Engineering-XDU_Project-Backend_app_models.AgentRelation": {
            "type": "object",
            "properties": {
                "dstAgentId": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "srcAgentId": {
                    "type": "integer"
                }
            }
        },
        "github_com_Web-Engineering-XDU_Project-Backend_app_models.Event": {
            "type": "object",
            "properties": {
                "contentHash": {
                    "type": "string",
                    "example": "47ed5d26b86a7519"
                },
                "createAt": {
                    "type": "string",
                    "example": "2023-04-11T05:07:53+08:00"
                },
                "deleteAt": {
                    "type": "string",
                    "example": "2023-04-11T06:07:53+08:00"
                },
                "error": {
                    "type": "boolean",
                    "example": false
                },
                "id": {
                    "type": "integer",
                    "example": 15
                },
                "jsonStr": {
                    "type": "string",
                    "example": "{\"ip_without_dot\":\"11125117081\",\"ip\":\"111.251.170.81\",\"latitude\":\"25.0504\"}"
                },
                "log": {
                    "type": "string",
                    "example": ""
                },
                "srcAgentId": {
                    "type": "integer",
                    "example": 1145
                }
            }
        },
        "swaggo.GetAgentListResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "msg": {
                    "type": "string",
                    "example": "ok"
                },
                "result": {
                    "$ref": "#/definitions/swaggo.GetAgentListResponseResult"
                }
            }
        },
        "swaggo.GetAgentListResponseResult": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_Web-Engineering-XDU_Project-Backend_app_models.AgentDetail"
                    }
                },
                "count": {
                    "type": "integer"
                }
            }
        },
        "swaggo.GetEventListResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "msg": {
                    "type": "string",
                    "example": "ok"
                },
                "result": {
                    "$ref": "#/definitions/swaggo.GetEventListResponseResult"
                }
            }
        },
        "swaggo.GetEventListResponseResult": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_Web-Engineering-XDU_Project-Backend_app_models.Event"
                    }
                },
                "count": {
                    "type": "integer"
                }
            }
        },
        "swaggo.GetRelationsResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "msg": {
                    "type": "string",
                    "example": "ok"
                },
                "result": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_Web-Engineering-XDU_Project-Backend_app_models.AgentRelation"
                    }
                }
            }
        },
        "swaggo.NewAgentResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "msg": {
                    "type": "string",
                    "example": "ok"
                },
                "result": {
                    "$ref": "#/definitions/swaggo.NewAgentResponseResult"
                }
            }
        },
        "swaggo.NewAgentResponseResult": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "swaggo.StateInfo": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "msg": {
                    "type": "string",
                    "example": "ok"
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
	Version:          "1.0",
	Host:             "43.142.105.98:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
