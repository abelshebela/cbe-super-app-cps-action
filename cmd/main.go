package main

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

// 	"github.com/rs/zerolog/log"
// )

// func main() {
// 	var (
// 		ctx    = context.Background()
// 		logger = utils.NewLogger()
// 	)

// 	defer logger.Sync()

// 	client, err := config.ConnectToMongoDB()

// 	if err != nil {
// 		log.Fatal().Msgf("Failed to connect to the Database: : %v", err)
// 	}

// 	defer func() {
// 		if err := client.Disconnect(ctx); err != nil {
// 			log.Fatal().Msgf("Error Occured while Disconnecting from Database: %v", err)
// 		}
// 	}()
// 	shutdown := make(chan error)

// 	go func() {
// 		quit := make(chan os.Signal, 1)
// 		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 		s := <-quit

// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()

// 		log.Info().Msgf("Signal caught %v", s.String())
// 		shutdown <- srv.Shutdown(ctx)
// 	}()

// 	log.Info().Msg(" Server is running")

// 	err = srv.ListenAndServe()
// 	if !errors.Is(err, http.ErrServerClosed) {
// 		log.Err(err).Msg("Server Failed")
// 		os.Exit(1)
// 	}

// 	err = <-shutdown
// 	if err != nil {
// 		log.Err(err).Msg("Server Shutdown Error")
// 		os.Exit(1)
// 	}
// }
