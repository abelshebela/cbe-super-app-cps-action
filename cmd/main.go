package main

import (
	"context"

	"cbe-super-app-cps-action/initiator"
)
// @title CPS Action API
// @version 1.0.0
// @description API documentation for CPS Action service

// @contact.name API Support
// @contact.email contact@eaglelionsystems.com

// @host localhost:8080
// @BasePath /api/v1/cbesuperapp/cps_action

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	initiator.Init(context.Background())
}
