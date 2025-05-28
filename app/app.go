package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vadim8q258475/store-cart-microservice/config"
	"github.com/vadim8q258475/store-cart-microservice/consumer"
	gen "github.com/vadim8q258475/store-cart-microservice/gen/v1"
	grpcService "github.com/vadim8q258475/store-cart-microservice/iternal/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	service  *grpcService.GrpcService
	server   *grpc.Server
	logger   zap.Logger
	port     string
	consumer consumer.Consumer
}

func NewApp(service *grpcService.GrpcService, server *grpc.Server, logger *zap.Logger, cfg config.Config, consumer consumer.Consumer) *App {
	return &App{
		service:  service,
		port:     cfg.Port,
		logger:   *logger,
		server:   server,
		consumer: consumer,
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return err
	}
	gen.RegisterCartServiceServer(a.server, a.service)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := a.server.Serve(l); err != nil {
			a.logger.Error("Server error", zap.Error(err))
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		a.consumer.Listen(ctx)
	}()

	<-stop

	a.logger.Info("Shutting down gRPC server...")
	a.server.GracefulStop()
	a.logger.Info("gRPC server stopped gracefully")

	return nil
}
