package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-gateway/internal/app/config"
	"payment-gateway/internal/app/repository/pgrepo"
	"payment-gateway/internal/app/services"
	"payment-gateway/internal/app/transport/httpserver"
	"payment-gateway/internal/pkg"
	"payment-gateway/internal/pkg/pg"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)

}

func run() error {
	cfg := config.Read()

	pgDB, err := pg.Dial(cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if pgDB != nil {
		log.Println("Running Postgres migrations")
		if err := runPgMigrations(cfg.DSN, cfg.MigrationsPath); err != nil {
			return fmt.Errorf("runPgMigrations failed: %w", err)
		}
	}
	redisClient := pkg.ConnectToRedis(cfg.RedisURL)
	failedTransactionRedisService := pkg.NewRedisClient(redisClient)

	// Set up the HTTP server and routes

	// Initialize repositories
	userRepo := pgrepo.NewUserRepo(pgDB)
	transactionRepo := pgrepo.NewTransactionRepo(pgDB)
	countryRepo := pgrepo.NewCountryRepo(pgDB)
	gatewayRepo := pgrepo.NewGatewayRepo(pgDB)
	mq, err := pkg.ConnectRabbitMQ(cfg.RabbitMQURL)
	encryptorService := services.NewDataEncryptorService()
	faultToleranceService := services.NewFaultToleranceService()
	gatewayService := services.NewGatewayService(gatewayRepo)
	failedTransactionPublisherService := services.NewTransactionPublisherService(faultToleranceService, mq, encryptorService)
	userService := services.NewUserService(
		userRepo,
		transactionRepo,
		countryRepo,
		gatewayRepo,
		faultToleranceService,
		failedTransactionRedisService,
	)
	failedTransactionsService := services.NewFailedTransactionsService(failedTransactionRedisService, mq, faultToleranceService, encryptorService, failedTransactionPublisherService)
	callbackService := services.NewCallbackService(mq, encryptorService, transactionRepo)
	callbackService.StartListeningJsonMessages()
	callbackService.StartListeningSoapMessages()
	failedTransactionsService.ProcessFailedTransactionsEveryTimePeriod(20 * time.Second)
	// Initialize HTTP server and routes

	server := httpserver.NewHttpServer(userService, gatewayService, failedTransactionPublisherService)
	router := mux.NewRouter()

	router.Handle("/deposit", http.HandlerFunc(server.DepositHandler)).Methods("POST")
	router.Handle("/withdraw", http.HandlerFunc(server.WithdrawalHandler)).Methods("POST")
	router.Handle("/callback", http.HandlerFunc(server.CallbackHandler)).Methods("POST")
	router.Handle("/gateway", http.HandlerFunc(server.UpdateGatewayPriority)).Methods("PUT")

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}
	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", cfg.HTTPAddr)

	// start HTTP server
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day!")
	return nil
}

func runPgMigrations(dsn, path string) error {
	if path == "" {
		return errors.New("no migrations path provided")
	}
	if dsn == "" {
		return errors.New("no DSN provided")
	}

	m, err := migrate.New(
		path,
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
