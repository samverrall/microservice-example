package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/samverrall/microservice-example/internal/app/user"
	"github.com/samverrall/microservice-example/internal/grpc"
	"github.com/samverrall/microservice-example/internal/grpc/handler"
	"github.com/samverrall/microservice-example/internal/grpc/health"
	userhandler "github.com/samverrall/microservice-example/internal/grpc/user"
	"github.com/samverrall/microservice-example/internal/postgres"
	"github.com/samverrall/microservice-example/pkg/config"
	"github.com/samverrall/microservice-example/pkg/proto"
	"google.golang.org/grpc/reflection"
)

var opts struct {
	grpc struct {
		port int
	}
}

func main() {
	flag.IntVar(&opts.grpc.port, "grpc-port", 8000, "grpc port")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger) // Updates slogs default instance of slog with our own handler.

	conf := config.NewConfig()
	if err := conf.Load(ctx); err != nil {
		logger.Error("failed to load service config", "error", err.Error())

		os.Exit(1)
	}
	// dbConf := conf.Database()

	// pgConn, err := postgresutil.Connect(ctx, postgresutil.BuildDSN(dbConf.Host, dbConf.Username, dbConf.Password, dbConf.DatabaseName, dbConf.Port))
	// if err != nil {
	// 	logger.Error("failed to connect to postgres", "error", err.Error())
	//
	// 	os.Exit(1)
	// }

	userRepo := postgres.NewUserRepository()
	userSvc := user.NewService(userRepo)

	appSvcs := handler.Service{
		User: userSvc,
	}
	appRepos := handler.Repo{
		User: userRepo,
	}

	baseHandler := handler.NewHandler(&appSvcs, &appRepos)

	grpcServer, err := grpc.NewServer(baseHandler, grpc.WithPort(opts.grpc.port))
	if err != nil {
		logger.Error("failed to init new grpc server", "error", err)

		os.Exit(1)
	}

	proto.RegisterHealthServer(grpcServer.Server(), health.NewHandler())
	proto.RegisterUserServer(grpcServer.Server(), userhandler.NewHandler(baseHandler))

	reflection.Register(grpcServer.Server())

	go func() {
		err := grpcServer.Serve()
		if err != nil {
			logger.Error("error serving grpc server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		logger.Error(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.Error(fmt.Sprintf("ctx.Done: %v", done))
	}
}
