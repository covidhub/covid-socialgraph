package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/covidhub/covid-socialgraph/pkg/database"
	"go.uber.org/zap"
)

func main() {
	var dbURL string
	var dbUser string

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	flag.StringVar(&dbURL, "db", "localhost:7678", "Neo4j Bolt URL")
	flag.StringVar(&dbUser, "user", "neo4j", "Neo4j Bolt URL")
	flag.Parse()

	dbUser = os.Getenv("COVIDHUB_DB_USER")
	dbPassword := os.Getenv("COVIDHUB_DB_PASSWORD")

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		zap.L().Sugar().Fatalw(
			"Failed to initialize logger",
			"err", err,
		)
	}
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()

	db, err := database.New(logger, dbURL, dbUser, dbPassword)
	if err != nil {
		logger.Fatalw(
			"Failed to connect to database",
			"err", err,
		)
	}

	<-signalChan
	logger.Info("Received shutdown signal, attempting graceful shutdown")
	err = db.Close()
	if err != nil {
		logger.Fatalw(
			"Failed to close database connection",
			"err", err,
		)
	}
}
