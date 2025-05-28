package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vadim8q258475/store-cart-microservice/app"
	"github.com/vadim8q258475/store-cart-microservice/config"
	"github.com/vadim8q258475/store-cart-microservice/consumer"
	grpcService "github.com/vadim8q258475/store-cart-microservice/iternal/grpc"
	"github.com/vadim8q258475/store-cart-microservice/iternal/interceptor"
	"github.com/vadim8q258475/store-cart-microservice/iternal/repo"
	service "github.com/vadim8q258475/store-cart-microservice/iternal/service/cart"
	productService "github.com/vadim8q258475/store-cart-microservice/iternal/service/product"
	userService "github.com/vadim8q258475/store-cart-microservice/iternal/service/user"
	productpbv1 "github.com/vadim8q258475/store-product-microservice/gen/v1"
	userpbv1 "github.com/vadim8q258475/store-user-microservice/gen/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO
// add cacher

func main() {
	// logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// interceptor
	intterceptor := interceptor.NewInterceptor(logger)

	// load config
	cfg := config.MustLoadConfig()
	fmt.Println(cfg.String())

	// user grpc client
	userConn, err := grpc.NewClient(cfg.UserHost+":"+cfg.UserPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer userConn.Close()

	userClient := userpbv1.NewUserServiceClient(userConn)

	// product grpc client
	productConn, err := grpc.NewClient(cfg.ProductHost+":"+cfg.ProductPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer productConn.Close()

	productClient := productpbv1.NewProductServiceClient(productConn)

	// db
	db, err := repo.InitDB(cfg)
	if err != nil {
		panic(err)
	}

	// repo
	cartProductRepo := repo.NewCartProductRepo(db)
	cartRepo := repo.NewCartRepo(db)

	// service
	productService := productService.NewproductService(productClient)
	userService := userService.NewUserService(userClient)
	service := service.NewCartService(productService, userService, cartRepo, cartProductRepo)

	// rabbitmq consumer
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	))

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(err)
	}
	consumer := consumer.NewRabbitMQConsumer(ch, service, cfg, logger)

	// grpc service
	grpcService := grpcService.NewGrpcService(service)

	// grpc server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			intterceptor.UnaryServerInterceptor,
		),
	)

	// app
	app := app.NewApp(grpcService, server, logger, cfg, consumer)

	if err = app.Run(); err != nil {
		panic(err)
	}
}
