// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/db"
	broker2 "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/adapter/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/adapter/repository"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/application"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/broker"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc/client"
	product2 "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc/product"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/middleware"
	product3 "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/order"
	product4 "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/payment"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http/product"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/observe"
)

// Injectors from wire.go:

func InitApplicationConfig(path string) *config.ApplicationConfig {
	applicationConfig := config.LoadApplicationConfig(path)
	return applicationConfig
}

func InitBootstrapConfig(path string) *bootstrap.BootstrapConfig {
	bootstrapConfig := bootstrap.LoadBootstrapConfig(path)
	return bootstrapConfig
}

func InitializeMigrator(app string, appCfg *config.ApplicationConfig) (*db.Migrator, error) {
	gormDB := db.NewDatabase(appCfg)
	migrator := db.NewMigrator(app, gormDB)
	return migrator, nil
}

func InitializeProductServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.ProductServer {
	engine := product.NewGinEngine(bootCfg)
	gormDB := db.NewDatabase(appCfg)
	productRepository := repository.NewGormProductRepository(gormDB)
	productUseCase := application.NewProductService(productRepository)
	productApplication := application.NewProductApplication(productUseCase)
	authConn := client.NewAuthConn(appCfg)
	authUseCase := application.NewAuthService(authConn)
	jwtAuthenticator := middleware.NewJwtAuthenticator(authUseCase)
	router := product.NewRouter(engine, productApplication, jwtAuthenticator)
	httpServer := product.New(bootCfg, engine, router)
	grpcProductServer := product2.NewGrpcProductServer(bootCfg, productUseCase)
	messageRouter := broker.InitializeRouter(bootCfg)
	natsPublisher := broker.NewNATSPublisher(bootCfg, appCfg)
	natsSubscriber := broker.NewNATSSubscriber(bootCfg, appCfg)
	sagaProductUseCase := application.NewSagaProductService(productRepository)
	sagaProductController := broker2.NewSagaProductController(sagaProductUseCase)
	eventRouter := broker2.NewProductEventRouter(messageRouter, natsPublisher, natsSubscriber, sagaProductController)
	tracerProvider := observe.NewTracer(bootCfg, appCfg)
	productServer := infrastructure.NewProductServer(httpServer, grpcProductServer, eventRouter, tracerProvider)
	return productServer
}

func InitializeOrderServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.OrderServer {
	engine := product3.NewGinEngine(bootCfg)
	gormDB := db.NewDatabase(appCfg)
	productConn := client.NewProductConn(appCfg)
	orderRepository := repository.NewGormOrderRepository(gormDB, productConn)
	orderUseCase := application.NewOrderService(orderRepository)
	orderApplication := application.NewOrderApplication(orderUseCase)
	authConn := client.NewAuthConn(appCfg)
	authUseCase := application.NewAuthService(authConn)
	jwtAuthenticator := middleware.NewJwtAuthenticator(authUseCase)
	router := product3.NewRouter(engine, orderApplication, jwtAuthenticator)
	httpServer := product3.New(bootCfg, engine, router)
	messageRouter := broker.InitializeRouter(bootCfg)
	natsPublisher := broker.NewNATSPublisher(bootCfg, appCfg)
	natsSubscriber := broker.NewNATSSubscriber(bootCfg, appCfg)
	sagaOrderUseCase := application.NewSagaOrderService(orderRepository)
	sagaOrderController := broker2.NewSagaOrderController(sagaOrderUseCase)
	eventRouter := broker2.NewOrderEventRouter(messageRouter, natsPublisher, natsSubscriber, sagaOrderController)
	tracerProvider := observe.NewTracer(bootCfg, appCfg)
	orderServer := infrastructure.NewOrderServer(httpServer, eventRouter, tracerProvider)
	return orderServer
}

func InitializePaymentServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.PaymentServer {
	engine := product4.NewGinEngine(bootCfg)
	gormDB := db.NewDatabase(appCfg)
	paymentRepository := repository.NewGormPaymentRepository(gormDB)
	paymentUseCase := application.NewPaymentService(paymentRepository)
	paymentApplication := application.NewPaymentApplication(paymentUseCase)
	authConn := client.NewAuthConn(appCfg)
	authUseCase := application.NewAuthService(authConn)
	jwtAuthenticator := middleware.NewJwtAuthenticator(authUseCase)
	router := product4.NewRouter(engine, paymentApplication, jwtAuthenticator)
	httpServer := product4.New(bootCfg, engine, router)
	messageRouter := broker.InitializeRouter(bootCfg)
	natsPublisher := broker.NewNATSPublisher(bootCfg, appCfg)
	natsSubscriber := broker.NewNATSSubscriber(bootCfg, appCfg)
	sagaPaymentUseCase := application.NewSagaPaymentService(paymentRepository)
	sagaPaymentController := broker2.NewSagaPaymentController(sagaPaymentUseCase)
	eventRouter := broker2.NewPaymentEventRouter(messageRouter, natsPublisher, natsSubscriber, sagaPaymentController)
	tracerProvider := observe.NewTracer(bootCfg, appCfg)
	paymentServer := infrastructure.NewPaymentServer(httpServer, eventRouter, tracerProvider)
	return paymentServer
}

func InitializeOrchestratorServer(appCfg *config.ApplicationConfig, bootCfg *bootstrap.BootstrapConfig) *infrastructure.OrchestratorServer {
	router := broker.InitializeRouter(bootCfg)
	natsPublisher := broker.NewNATSPublisher(bootCfg, appCfg)
	natsSubscriber := broker.NewNATSSubscriber(bootCfg, appCfg)
	redisPublisher := broker.NewRedisPublisher(bootCfg, appCfg)
	purchaseResultRepository := broker2.NewPurchaseResultPublisher(redisPublisher)
	orchestratorUseCase := application.NewOrchestratorService(natsPublisher, purchaseResultRepository)
	sagaOrchestratorController := broker2.NewSagaOrchestratorController(orchestratorUseCase)
	eventRouter := broker2.NewOrchestratorEventRouter(router, natsPublisher, natsSubscriber, sagaOrchestratorController)
	tracerProvider := observe.NewTracer(bootCfg, appCfg)
	orchestratorServer := infrastructure.NewOrchestratorServer(eventRouter, tracerProvider)
	return orchestratorServer
}
