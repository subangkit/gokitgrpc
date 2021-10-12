package main

import (
	"context"
	"flag"
	"fmt"
	//"time"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"

	"net"
	"os"
	"os/signal"
	"syscall"

	"gokitgrpc/user"
	pb "gokitgrpc/user/pb"

	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
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
	//ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	var srv user.Service
	{
		repository := user.NewRepo(db, logger)

		srv = user.NewService(repository, logger)
	}

	endpoints := user.MakeEndpoints(srv)
	grpcServer := user.NewGRPCServer(endpoints, logger)
	
	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", ":3031")
	if err != nil {
		logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}
	go func() {
		baseServer := grpc.NewServer()
        pb.RegisterUserGrpcServer(baseServer, grpcServer)
        fmt.Println("Server gRPC started successfully on port 3031")
        baseServer.Serve(grpcListener)
	}()

	level.Error(logger).Log("exit", <-errs)
}