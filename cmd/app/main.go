package main

import (
	"context"
	"database/sql"
	"github.com/calebikhuohon/hasty-test/cmd/app/internal/handler"
	"github.com/calebikhuohon/hasty-test/internal/service"
	"github.com/calebikhuohon/hasty-test/internal/storage"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main()  {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		systemCall := <-exit
		log.Printf("System call: %+v", systemCall)
		cancel()
	}()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to parse app config: %s", err)
	}
	if err := config.Validate(); err != nil {
		log.Fatalf("Failed to validate config: %s", err)
	}


	db, err := newDB(config)
	if err != nil {
		log.Fatalf("could no initialize db: %s", err)
	}

	rdb := newRedisCache(config)
	storage := storage.New(db)
	jobService := service.NewJobService(storage, config.JobMaxRetries, rdb)
	r := handler.New(jobService, handler.Config{Timeout: calculateHTTPTimeout("")})
	srv := NewServer(config.Ports.HTTP, r)

	go func() {
		log.Print("Starting http server..")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Could not listen and serve: %s", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	log.Print("Stopping the http server..")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Could not shut down gracefully: %s", err)
		os.Exit(1)
	}

	defer cancel()
}

func newRedisCache(config Config) *redis.Client  {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Redis.Address,
		Password: config.Redis.Password,
		DB: 0,
	})

	return rdb
}

func calculateHTTPTimeout(env string) time.Duration {
	switch env {
	case "dev":
		return 5 * time.Minute
	default:
		return 10 * time.Second
	}
}

func newDB(config Config) (*sql.DB, error) {
	conn := mysql.Config{
		User:                 config.MySQL.User,
		Passwd:               config.MySQL.Passwd,
		Addr:                 config.MySQL.Host,
		DBName:               config.MySQL.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
		Params:               map[string]string{"charset": "utf8"},
		MultiStatements:      true,
		Timeout:              5 * time.Second,
		Net:                  "tcp",
	}
	db, err := sql.Open("mysql", conn.FormatDSN())
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	return db, nil
}

func NewServer(addr string, h http.Handler) *http.Server {
	s := &http.Server{Addr: addr, Handler: h}
	return s
}
