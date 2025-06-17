package main

import (
	"context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/amount_based_auth_handler"
	persistence "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/amount_based_auth"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/amount_based_auth_app"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/amount_based_auth"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	branch_handler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/branch_handler"
	customerhandler "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/inbound/http/customer_handler"
	branch_repo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/branch"
	customerPersistance "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/adapter/outbound/persistence/customer"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/customer"
	branch_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/bulkcustomer/services"
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

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

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

	amountBasedAuthRepo := persistence.InitAmountBasedAuth(mongoClient, cfg.MongoDBDatabase, "authTier")
	amountBasedAuthService := amount_based_auth_domain.NewAmountBasedAuthService(amountBasedAuthRepo, ctx)
	amountBasedAuthApplication := amount_based_auth_app.AmountBasedAuthHandler(amountBasedAuthService)
	amountBasedAuthHandler := amount_based_auth_handler.NewAmountBasedAuthHandler(amountBasedAuthApplication)
	amount_based_auth_handler.InitAmountBasedAuthHandler(r, amountBasedAuthHandler)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerPort),
		Handler: r,
	}

	logger.Infof("Starting server at %v", server.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("🚀 Server started on", server.Addr)
		log.Printf("Server stopped with error: %v\n", server.ListenAndServe())
	}()

	sig := <-quit
	log.Printf("server shutting down with signal: %v\n", sig)

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown gracefully with error %v", err)
	}

	log.Println("Server shutdown successfully")
}
