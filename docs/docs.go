// Package docs provides Swagger documentation.
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
        "/health": {
            "get": {
                "description": "Check the health status of the service",
                "produces": ["application/json"],
                "tags": ["health"],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "Service is healthy",
                        "schema": {
                            "$ref": "#/definitions/HealthResponse"
                        }
                    },
                    "503": {
                        "description": "Service is unhealthy",
                        "schema": {
                            "$ref": "#/definitions/HealthResponse"
                        }
                    }
                }
            }
        },
        "/tracks/{id}": {
            "get": {
                "description": "Get track information by ID",
                "produces": ["application/json"],
                "tags": ["tracks"],
                "summary": "Get track",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Track ID (MD5 hash)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Track found",
                        "schema": {
                            "$ref": "#/definitions/GetTrackResponse"
                        }
                    },
                    "404": {
                        "description": "Track not found"
                    }
                }
            }
        },
        "/tracks": {
            "post": {
                "description": "Create or update a track",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["tracks"],
                "summary": "Upsert track",
                "parameters": [
                    {
                        "description": "Track data",
                        "name": "track",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateTrackRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Track created/updated",
                        "schema": {
                            "$ref": "#/definitions/TrackResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    }
                }
            }
        },
        "/tracks/{trackId}/like": {
            "post": {
                "description": "Add a like reaction to a track",
                "tags": ["reactions"],
                "summary": "Like track",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Track ID",
                        "name": "trackId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User-ID",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Reaction added"
                    },
                    "409": {
                        "description": "Already reacted"
                    }
                }
            }
        },
        "/tracks/{trackId}/dislike": {
            "post": {
                "description": "Add a dislike reaction to a track",
                "tags": ["reactions"],
                "summary": "Dislike track",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Track ID",
                        "name": "trackId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User-ID",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Reaction added"
                    },
                    "409": {
                        "description": "Already reacted"
                    }
                }
            }
        },
        "/tracks/{trackId}/reaction": {
            "get": {
                "description": "Check if user has reacted to a track",
                "produces": ["application/json"],
                "tags": ["reactions"],
                "summary": "Check reaction",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Track ID",
                        "name": "trackId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "X-User-ID",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reaction status",
                        "schema": {
                            "$ref": "#/definitions/CheckReactionResponse"
                        }
                    }
                }
            }
        },
        "/radio/info": {
            "get": {
                "description": "Get current radio stream information",
                "produces": ["application/json"],
                "tags": ["radio"],
                "summary": "Get radio info",
                "responses": {
                    "200": {
                        "description": "Radio info",
                        "schema": {
                            "$ref": "#/definitions/RadioInfoResponse"
                        }
                    }
                }
            }
        },
        "/radio/listeners": {
            "get": {
                "description": "Get current listener count",
                "produces": ["application/json"],
                "tags": ["radio"],
                "summary": "Get listeners",
                "responses": {
                    "200": {
                        "description": "Listener count",
                        "schema": {
                            "$ref": "#/definitions/ListenerResponse"
                        }
                    }
                }
            }
        },
        "/radio/statistics": {
            "get": {
                "description": "Get track statistics including history, top listened, top rotated, top likes and dislikes",
                "produces": ["application/json"],
                "tags": ["radio"],
                "summary": "Get statistics",
                "responses": {
                    "200": {
                        "description": "Statistics",
                        "schema": {
                            "$ref": "#/definitions/StatisticsResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "HealthResponse": {
            "type": "object",
            "properties": {
                "status": {"type": "string"},
                "checks": {"type": "object"}
            }
        },
        "CreateTrackRequest": {
            "type": "object",
            "required": ["Md5", "StreamTitle"],
            "properties": {
                "Md5": {"type": "string"},
                "StreamTitle": {"type": "string"},
                "StreamUrl": {"type": "string"}
            }
        },
        "TrackResponse": {
            "type": "object",
            "properties": {
                "rotate": {"type": "integer"}
            }
        },
        "GetTrackResponse": {
            "type": "object",
            "properties": {
                "id": {"type": "string"},
                "title": {"type": "string"},
                "cover": {"type": "string"},
                "rotate": {"type": "integer"},
                "likes": {"type": "integer"},
                "dislikes": {"type": "integer"},
                "listeners": {"type": "integer"}
            }
        },
        "CheckReactionResponse": {
            "type": "object",
            "properties": {
                "hasReacted": {"type": "boolean"},
                "reaction": {"type": "string"}
            }
        },
        "RadioInfoResponse": {
            "type": "object",
            "properties": {
                "name": {"type": "string"},
                "description": {"type": "string"},
                "streamUrl": {"type": "string"},
                "listener": {"$ref": "#/definitions/ListenerResponse"}
            }
        },
        "ListenerResponse": {
            "type": "object",
            "properties": {
                "current": {"type": "integer"},
                "peak": {"type": "integer"}
            }
        },
        "StatisticsResponse": {
            "type": "object",
            "properties": {
                "statistics": {"type": "array", "items": {"$ref": "#/definitions/Category"}}
            }
        },
        "Category": {
            "type": "object",
            "properties": {
                "key": {"type": "string"},
                "icon": {"type": "string"},
                "tracks": {"type": "array", "items": {"$ref": "#/definitions/TrackStats"}}
            }
        },
        "TrackStats": {
            "type": "object",
            "properties": {
                "title": {"type": "string"},
                "cover": {"type": "string"},
                "rotate": {"type": "integer"},
                "likes": {"type": "integer"},
                "dislikes": {"type": "integer"},
                "listeners": {"type": "integer"}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Radio Streaming API",
	Description:      "API for radio streaming platform with tracks, reactions, and statistics",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
