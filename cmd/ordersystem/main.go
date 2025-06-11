package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/configs"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/event/handler"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/infra/graph"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/infra/grpc/pb"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/infra/grpc/service"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/infra/web"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/infra/web/webserver"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/internal/usecase"
	"github.com/mobenaus/fc-pos-go-desafio-cleancode/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs := loadConfigurations()

	db := getDbConnection(configs)
	defer db.Close()

	eventDispatcher := createEventDispatcher(configs)

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	go startWebServer(createOrderUseCase, listOrdersUseCase, configs)

	go startGRPCServer(createOrderUseCase, listOrdersUseCase, configs)

	go startGraphQLSever(createOrderUseCase, listOrdersUseCase, configs)

	forever := make(chan bool)

	<-forever
}

func startGraphQLSever(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase, configs *configs.Conf) {
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func startGRPCServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase, configs *configs.Conf) {
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	grpcServer.Serve(lis)
}

func startWebServer(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase, configs *configs.Conf) {
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(*createOrderUseCase, *listOrdersUseCase)
	webserver.AddPOSTHandler("/order", webOrderHandler.Create)
	webserver.AddGETHandler("/orders", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start()
}

func createEventDispatcher(configs *configs.Conf) *events.EventDispatcher {
	rabbitMQChannel := getRabbitMQChannel(configs)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	return eventDispatcher
}

func getDbConnection(configs *configs.Conf) *sql.DB {
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	return db
}

func loadConfigurations() *configs.Conf {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	return configs
}

func getRabbitMQChannel(configs *configs.Conf) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", configs.AMQPUser, configs.AMQPPassword, configs.AMQPHost, configs.AMQPPort))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
