package initiator

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/cmd/server"
	local "github.com/CBE-Super-App/cbe-super-app-member-auth/config"

	// "github.com/CBE-Super-App/cbe-super-app-member-auth/platform/logger"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

func Init(ctx context.Context) {
	logger := utils.NewLogger()
	logger.Infof("Initializing configuration...")
	cfg := InitConfig(logger)
	logger.Infof("Configuration initialized")

	logger.Infof("Initializing MongoDB client...")
	mongoClient := InitMongo(cfg.MongoDBURI, logger)
	logger.Infof("MongoDB client initialized")

	logger.Infof("Initializing Minio client...")
	minioClient := InitMinio(cfg.MinioEndPoint, cfg.MinioAccessKey, cfg.MinioSecretKey, logger)
	logger.Infof("Minio client initialized")

	logger.Infof("Initializing persistence...")
	persitence := InitPersistanceLayer(mongoClient, cfg.MongoDBDatabase, logger)
	logger.Infof("Persistence initialized")

	defer local.DisconnectMongo(ctx, mongoClient, logger)

	logger.Infof("initialize service layer")
	serviceLayer := InitServiceLayer(persitence, logger, cfg, minioClient)

	logger.Infof("initialize handler layer")
	handlerLayer := InitHandler(serviceLayer.UserService, logger)

	r := chi.NewRouter()
	InitRoute(ctx, r, handlerLayer, logger)

	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	time.Sleep(5 * time.Second)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Goroutines: ", runtime.NumGoroutine())
	}()
	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	srv := server.NewHTTPServer(cfg, r)

	go func() {
		srv.HTTPServerStart(ctx, logger)
	}()
	srv.HTTPServerStop(ctx, logger)
}
