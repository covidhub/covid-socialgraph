package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/covidhub/covid-socialgraph/pkg/database"
	"github.com/covidhub/covid-socialgraph/pkg/server"
	"go.uber.org/zap"
)

func main() {
	dbUser := os.Getenv("COVIDHUB_DB_USER")
	dbPassword := os.Getenv("COVIDHUB_DB_PASSWORD")

	tls := flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile := flag.String("cert_file", "", "The TLS cert file")
	keyFile := flag.String("key_file", "", "The TLS key file")

	grpcPort := flag.Int("port", 8080, "Port the gRPC server listens to")
	grpcHost := flag.String("host", "", "gRPC server host address")

	dbURL := flag.String("db", "bolt://localhost:7678", "Neo4j Bolt URL")

	flag.Parse()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		zap.L().Sugar().Fatalw(
			"Failed to initialize logger",
			"err", err,
		)
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	var db *database.Client

	logger.Info("Connecting to database")
	db, err = database.New(logger, *dbURL, dbUser, dbPassword)
	if err != nil {
		logger.Fatalw(
			"Failed to connect to database",
			"err", err,
		)
	}

	logger.Info("Starting grpc server")
	url := fmt.Sprintf("%s:%d", *grpcHost, *grpcPort)
	server, err := server.New(url, *tls, *certFile, *keyFile, db, logger)
	if err != nil {
		logger.Fatalw(
			"Failed to start grpc server",
			"err", err,
		)
	}

	go server.Serve()

	logger.Info("Server ready")

	<-signalChan
	logger.Info("Received shutdown signal, attempting graceful shutdown")

	server.Close()

	err = db.Close()
	if err != nil {
		logger.Fatalw(
			"Failed to close database connection",
			"err", err,
		)
	}
}
