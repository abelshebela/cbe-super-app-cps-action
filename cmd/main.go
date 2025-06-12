package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/pkg/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/pkg/utils"

	branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/branch_handler"
	branch_repo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/catch"
	branch_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/services"
)

func main() {
	logger := utils.NewLogger()
	defer logger.Sync()

if _, err := config.Load(); err != nil {
    logger.Fatalf("failed to load config %v", err)
}
	mongoURI := "mongodb://localhost:27017"

	mongoClient, err := config.ConnectToMongoDB(mongoURI)

	if err != nil {
		logger.Fatalf("failed to connect to mongo %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()
	log.Println("Connected to MongoDB!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

branchPersistence := branch_repo.NewBranchPersistence(mongoClient, "member", 5*time.Second)
	branchService := branch_domain.NewBranchService(branchPersistence)
	branchHandler := branch_handler.NewBranchHandler(branchService, logger)
	branch_handler.RegisterBranchRoutes(r, branchHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("🚀 Server started on http://localhost:8080")
		log.Printf("Server stopped with error: %v\n", server.ListenAndServe())
	}()

	sig := <-quit
	log.Printf("server shutting down with signal: %v\n", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown gracefully with error %v", err)
	}

	log.Println("Server shutdown successfully")
}
