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
    "github.com/spf13/viper"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

    branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/branch_handler"
    branch_repo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/branch"
    branch_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/services"
    customerhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/customer_handler"
    customerPersistance "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/customer"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/customer"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/customer/service"

    faydaRoutes "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/fayda_account"
    faydaaccount "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/fayda_account"
    faydaHandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/fayda_account"
    faydaService "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/fayda_account/service"
)

func main() {
    logger := utils.NewLogger()
    defer logger.Sync()

    cfg, err := config.Load()
    if err != nil {
        logger.Fatalf("failed to load config %v", err)
    }

    mongoClient, err := config.ConnectToMongoDB(cfg.MongoDBURI)
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

    branchPersistence := branch_repo.NewBranchPersistence(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
    branchService := branch_domain.NewBranchService(branchPersistence)
    branchHandler := branch_handler.NewBranchHandler(branchService, logger)
    branch_handler.RegisterBranchRoutes(r, branchHandler)

    customerPersitance := customerPersistance.InitCustomerDetail(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
    customerDomain := service.IntiCustomerDomain(customerPersitance, logger)
    customerApp := customer.InitCustomerHandler(customerDomain, logger)
    customerRoutes := customerhandler.NewCustomerHTTPHandler(customerApp, logger)
    customerhandler.InitCustomerRoutes(r, customerRoutes)

    faydaPersistence := faydaaccount.InitFaydaAccountPersistence(mongoClient, cfg.MongoDBDatabase, viper.GetDuration("timeout"), logger)
    faydaDomin := faydaService.InitFaydaAccountDomain(faydaPersistence, logger)
    faydaApp := faydaHandler.InitFaydaHandler(faydaDomin, logger)
    faydaHandlers := faydaRoutes.InitFaydaAdapter(faydaApp, logger)
    faydaRoutes.InitFaydaRoutes(r, faydaHandlers)

    server := http.Server{
        Addr:    viper.GetString("Host") + ":" + viper.GetString("Port"),
        Handler: r,
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    go func() {
        log.Println("🚀 Server started on", server.Addr)
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