package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	departmenthandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/department_handler"
	departmentPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/department"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/department"
	departmentService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/department"


	permissionhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/permission_handler"
	permissionPersistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/permission"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/permission"
	permissionService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/permission"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func main() {
	logger := utils.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	mongoClient, err := config.ConnectToMongoDB(cfg.MongoDBURI)
	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	departmentPersistence := departmentPersistence.InitDepartment(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
	departmentDomain := departmentService.InitDepartmentDomain(departmentPersistence, departmentPersistence, logger)
	departmentApp := department.InitDepartmentHandler(departmentDomain, logger)
	departmentRoutes := departmenthandler.NewDepartmentHTTPHandler(departmentApp, logger)
	departmenthandler.InitDepartmentRoutes(r, departmentRoutes)

    permissionPersistence := permissionPersistence.InitPermission(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
	permissionDomain := permissionService.InitPermissionDomain(permissionPersistence, permissionPersistence, permissionPersistence, logger)
	permissionApp := permission.InitPermissionHandler(permissionDomain, logger)
	permissionRoutes := permissionhandler.NewPermissionHTTPHandler(permissionApp, logger)
	permissionhandler.InitPermissionRoutes(r, permissionRoutes)

	server := http.Server{
		Addr:   ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		log.Println("🚀 Server started on", viper.GetString("Port"))
		log.Printf("Server stopped with error: %v\n", server.ListenAndServe())
	}()

	sig := <-quit

	log.Printf("server shutting down with signal: %v\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown gracefully with error %v", err)
	}

	log.Println("Server shutdown successfully")
}
