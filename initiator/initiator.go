package initiator

import (
	"cbe-super-app-budget/cmd/server"
	local "cbe-super-app-budget/config"
	"cbe-super-app-budget/platform/logger"
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Init(ctx context.Context) {
	config := zap.NewProductionConfig()
	build, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	log := logger.New(build, logger.Options{})

	env, err := local.LoadVault(log)
	if err != nil {
		log.Fatal(ctx, "failed to load configuration", zap.Error(err))
	}

	client, db, err := local.ConnectMongo(log, env)
	if err != nil {
		log.Fatal(ctx, "failed to connect to db", zap.Error(err))
	}

	defer local.DisconnectMongo(ctx, client, log)

	log.Info(ctx, "initialize persistance layer")
	persistanceLayer := InitPersistanceLayer(db, log)

	log.Info(ctx, "initialize service layer")
	serviceLayer := InitServiceLayer(persistanceLayer, log)

	log.Info(ctx, "initialize handler layer")
	handlerLayer := InitHandlerLayer(serviceLayer, log)

	r := chi.NewRouter()
	InitRoute(ctx, r, handlerLayer, log)

	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	time.Sleep(5 * time.Second)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Goroutines: ", runtime.NumGoroutine())
	}()
	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	srv := server.NewHTTPServer(env, r)

	go func() {
		srv.HTTPServerStart(ctx, log)
	}()
	srv.HTTPServerStop(ctx, log)
}
