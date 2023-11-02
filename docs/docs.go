// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "So Nim Wai"
        },
        "license": {
            "name": "LSCM"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/addDeliveryOrder": {
            "post": {
                "description": "Create adhoc delivery order.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order Management"
                ],
                "summary": "Add Delivery Order.",
                "parameters": [
                    {
                        "description": "Add Delivery Order Parameters",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AddDeliveryOrderDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderManagement.AddDeliveryOrderResponse"
                        }
                    }
                }
            }
        },
        "/cancelDeliveryOrder": {
            "post": {
                "description": "Update Non Started Delivery Order .",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order Management"
                ],
                "summary": "Cancel Delivery Order.",
                "parameters": [
                    {
                        "description": "Cancel Delivery Parameters",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.CancelDeliveryOrderDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderManagement.CancelDeliveryOrderResponse"
                        }
                    }
                }
            }
        },
        "/getDeliveryOrder": {
            "get": {
                "description": "Get the list of delivery order by order status .",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order Management"
                ],
                "summary": "Get Delivery Order.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderManagement.OrderListBody"
                        }
                    }
                }
            }
        },
        "/getDutyRooms": {
            "get": {
                "description": "Get the list of location.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Map Handling"
                ],
                "summary": "Get Duty Rooms.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/mapHandling.GetDutyRoomsResponse"
                        }
                    }
                }
            }
        },
        "/getFloorPlan": {
            "get": {
                "description": "Get UI Floor Plan.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Map Handling"
                ],
                "summary": "Get Floor Plan.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/mapHandling.GetFloorPlanResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/loginAdmin": {
            "post": {
                "description": "Login to OMS.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Login Auth"
                ],
                "summary": "Login to OMS.",
                "parameters": [
                    {
                        "description": "Login Parameters",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginAdminDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/loginAuth.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.FailResponse"
                        }
                    }
                }
            }
        },
        "/loginStaff": {
            "post": {
                "description": "Login to OMS.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Login Auth"
                ],
                "summary": "Login to OMS.",
                "parameters": [
                    {
                        "description": "Login Parameters",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginStaffDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/loginAuth.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.FailResponse"
                        }
                    }
                }
            }
        },
        "/logout": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Logout from OMS.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Login Auth"
                ],
                "summary": "Logout from OMS.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/loginAuth.LogoutResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.FailResponse"
                        }
                    }
                }
            }
        },
        "/renewToken": {
            "get": {
                "description": "Using Valid Token to renew token before expired",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Login Auth"
                ],
                "summary": "Websocket Connection.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/loginAuth.LogoutResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.FailResponse"
                        }
                    }
                }
            }
        },
        "/testAW2": {
            "get": {
                "description": "Get the response of AW2 (Server notify the user which location selected).",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Test"
                ],
                "summary": "Test AW2 websocket response.",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/testHW1": {
            "get": {
                "description": "Get the response of MW1 (Server report robot status and location (every 1s) ).",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Test"
                ],
                "summary": "Test MW1 websocket response.",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/testOW1": {
            "get": {
                "description": "Get the response of OW2 (Server notify any of created order status changed).",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Test"
                ],
                "summary": "Test OW2 websocket response.",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/testSW1": {
            "get": {
                "description": "Get the response of SW1 (Server report robot status and location (every 1s) ).",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Test"
                ],
                "summary": "Test SW1 websocket response.",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/triggerHandlingOrder": {
            "get": {
                "description": "Notify system users are ready to handle the current order.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order Management"
                ],
                "summary": "Trigger Delivery Order.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Order IDs",
                        "name": "orderIds",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Schedule IDs",
                        "name": "scheduleId",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderManagement.TriggerHandlingOrderResponse"
                        }
                    }
                }
            }
        },
        "/updateDeliveryOrder": {
            "get": {
                "description": "Update Non Started Delivery Order .",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order Management"
                ],
                "summary": "Update Delivery Order.",
                "parameters": [
                    {
                        "description": "Update Delivery Order Parameters",
                        "name": "todo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.UpdateDeliveryOrderDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderManagement.UpdateDeliveryOrderResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.AddDeliveryOrderDTO": {
            "type": "object",
            "properties": {
                "endLocationId": {
                    "type": "integer"
                },
                "endLocationName": {
                    "type": "string"
                },
                "expectingDeliveryTime": {
                    "type": "string"
                },
                "expectingStartTime": {
                    "type": "string"
                },
                "numberOfAmrRequire": {
                    "type": "integer"
                },
                "orderType": {
                    "type": "string"
                },
                "startLocationId": {
                    "type": "integer"
                },
                "startLocationName": {
                    "type": "string"
                }
            }
        },
        "handlers.CancelDeliveryOrderDTO": {
            "type": "object",
            "properties": {
                "scheduleId": {
                    "type": "integer"
                }
            }
        },
        "handlers.LoginAdminDTO": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handlers.LoginStaffDTO": {
            "type": "object",
            "properties": {
                "dutyLocationId": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handlers.UpdateDeliveryOrderDTO": {
            "type": "object",
            "properties": {
                "endLocationId": {
                    "type": "integer"
                },
                "endLocationName": {
                    "type": "string"
                },
                "expectingDeliveryTime": {
                    "type": "string"
                },
                "expectingStartTime": {
                    "type": "string"
                },
                "numberOfAmrRequire": {
                    "type": "integer"
                },
                "scheduleId": {
                    "type": "integer"
                },
                "startLocationId": {
                    "type": "integer"
                },
                "startLocationName": {
                    "type": "string"
                }
            }
        },
        "loginAuth.LoginResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/loginAuth.LoginResponseBody"
                },
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        },
        "loginAuth.LoginResponseBody": {
            "type": "object",
            "properties": {
                "authToken": {
                    "type": "string"
                },
                "dutyLocationId": {
                    "type": "integer"
                },
                "dutyLocationName": {
                    "type": "string"
                },
                "loginDateTime": {
                    "type": "string"
                },
                "tokenExpiryDateTime": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                },
                "userType": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "loginAuth.LogoutResponse": {
            "type": "object",
            "properties": {
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        },
        "mapHandling.GetDutyRoomsResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/mapHandling.LocationListBody"
                },
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        },
        "mapHandling.GetFloorPlanResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/mapHandling.MapListBody"
                },
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        },
        "mapHandling.LocationListBody": {
            "type": "object",
            "properties": {
                "locationList": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "locationId": {
                                "type": "integer"
                            },
                            "locationName": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "mapHandling.MapListBody": {
            "type": "object",
            "properties": {
                "mapList": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "floorId": {
                                "type": "integer"
                            },
                            "floorImage": {
                                "type": "string"
                            },
                            "floorName": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "models.FailResponse": {
            "type": "object",
            "properties": {
                "header": {
                    "$ref": "#/definitions/models.FailResponseHeader"
                }
            }
        },
        "models.FailResponseHeader": {
            "type": "object",
            "properties": {
                "failedReason": {
                    "type": "string"
                },
                "responseCode": {
                    "type": "integer"
                },
                "responseMessage": {
                    "type": "string"
                }
            }
        },
        "models.ResponseHeader": {
            "type": "object",
            "properties": {
                "responseCode": {
                    "type": "integer"
                },
                "responseMessage": {
                    "type": "string"
                }
            }
        },
        "orderManagement.AddDeliveryOrderResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/orderManagement.OrderListBody"
                },
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        },
        "orderManagement.CancelDeliveryOrderResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/orderManagement.OrderListBody"
                },
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        },
        "orderManagement.OrderListBody": {
            "type": "object",
            "properties": {
                "orderList": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "endLocationId": {
                                "type": "integer"
                            },
                            "endLocationName": {
                                "type": "string"
                            },
                            "endTime": {
                                "type": "string"
                            },
                            "expectingDeliveryTime": {
                                "type": "string"
                            },
                            "expectingStartTime": {
                                "type": "string"
                            },
                            "orderCreatedBy": {
                                "type": "integer"
                            },
                            "orderCreatedType": {
                                "type": "string"
                            },
                            "orderId": {
                                "type": "integer"
                            },
                            "orderStatus": {
                                "type": "string"
                            },
                            "orderType": {
                                "type": "string"
                            },
                            "processingStatus": {
                                "type": "string"
                            },
                            "scheduleId": {
                                "type": "integer"
                            },
                            "startLocationId": {
                                "type": "integer"
                            },
                            "startLocationName": {
                                "type": "string"
                            },
                            "startTime": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "orderManagement.TriggerHandlingOrderResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/orderManagement.OrderListBody"
                },
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        },
        "orderManagement.UpdateDeliveryOrderResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/orderManagement.OrderListBody"
                },
                "header": {
                    "$ref": "#/definitions/models.ResponseHeader"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "ams.lscm.hk",
	BasePath:         "/oms/",
	Schemes:          []string{},
	Title:            "TKOH OMS",
	Description:      "TKOH OMS Backend Server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
