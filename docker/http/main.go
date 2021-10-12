package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gokitgrpc/user"

	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var httpAddr = flag.String("http", ":3000", "http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "user",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://godb:123456@mongodb:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		level.Error(logger).Log("exit", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		level.Error(logger).Log("exit", err)
	}

	var db *mongo.Collection
	db = client.Database("godb").Collection("users")

	fmt.Println("Connected to MongoDB!")

	flag.Parse()
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	var srv user.Service
	{
		repository := user.NewRepo(db, logger)

		srv = user.NewService(repository, logger)
	}

	endpoints := user.MakeEndpoints(srv)
	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := user.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()


	level.Error(logger).Log("exit", <-errs)
}