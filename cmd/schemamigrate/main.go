package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

var directionFlag = flag.String("direction", "up", "The migration direction: up|down. Defaults to 'up'")
var directoryFlag = flag.String("directory", "/app/migrations/", "Directory where migration schemas as present. Defaults to '/app/migrations/'")

type logger struct {
	*log.Logger
}

func (l logger) Verbose() bool {
	return false
}

func main() {
	flag.Parse()

	dbConnStr := os.Getenv("MYSQL_CONNECTION_STRING")
	if dbConnStr == "" {
		log.Fatalf("MYSQL_CONNECTION_STRING env variable not set")
	}
	databaseURL := fmt.Sprintf("mysql://%s", dbConnStr)

	log.Info("Sleeping for 20 to wait for db to be reachable")

	retries := 3
	for i := 0; i < retries; i++ {
		err := run(databaseURL)
		if err != nil && strings.Contains(err.Error(), "i/o timeout") {
			log.Info("Retrying after 20s..")
			time.Sleep(20 * time.Second)
		} else if err != nil {
			log.Fatalf("Cannot process job")
			return
		} else {
			return
		}
	}
	log.Fatal("Retries exhausted")
}

func run(databaseURL string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", *directoryFlag), databaseURL)
	if err != nil {
		log.WithError(err).Errorf("Failed to initialize migrate client")
		return err
	}
	m.Log = &logger{
		Logger: log.New(),
	}

	switch *directionFlag {
	case "up":
		err := m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.WithError(err).Fatalf("Failed to run up migrations")
		}
	case "down":
		err := m.Down()
		if err != nil {
			log.WithError(err).Fatalf("Failed to run up migrations")
		}
	default:
		log.WithField("direction", directionFlag).Fatalf("Unexpected value of direction flag")
	}

	return nil
}
