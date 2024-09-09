package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/PavelDonchenko/template/MICROSERVICES/ecommerce-micro/common/config"
	"github.com/PavelDonchenko/template/MICROSERVICES/ecommerce-micro/common/consul"
	"github.com/PavelDonchenko/template/MICROSERVICES/ecommerce-micro/common/helpers"
	"github.com/PavelDonchenko/template/MICROSERVICES/ecommerce-micro/common/logger"
	"github.com/PavelDonchenko/template/MICROSERVICES/ecommerce-micro/common/trace/otel/jaeger"
	"github.com/PavelDonchenko/template/MICROSERVICES/ecommerce-micro/common/validator"
)

var (
	production   *bool
	disableTrace *bool
)

// @title           Go E-commerce micro
// @version         1.0
// @description     Authentication server.
// @termsOfService  http://swagger.io/terms/

// @contact.email  przmld033@gmail.com

// @BasePath  /api/v1
// @Schemes  https

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description		Type "Bearer" followed by a space and JWT token.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	production = flag.Bool("prod", false, "use -prod=true to run in production mode")
	disableTrace = flag.Bool("disable-trace", false, "use disable-trace=true if you want to disable tracing completely")
	flag.Parse()

	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		//return web.GetTraceID(ctx)
		return "00000000-0000-0000-0000-000000000000"
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "AUTH", traceIDFn, events)
	// -------------------------------------------------------------------------

	cfg := config.MustLoadConfig(*production, "./auth/config/")

	log.Info(ctx, "Init config", slog.Bool("production mode", cfg.Production))

	validator.NewValidator("en")
	helpers.CreateFolders(cfg.Folders)

	run(ctx, log, cfg)
}

func run(ctx context.Context, log *logger.Logger, cfg *config.Config) error {
	//------------------------Init Trace------------------------------------------
	providerTracer, err := jaeger.NewProvider(jaeger.ProviderConfig{
		JaegerEndpoint: cfg.Jaeger.JaegerEndpoint,
		ServiceName:    cfg.Jaeger.ServiceName,
		ServiceVersion: cfg.Jaeger.ServiceVersion,
		Production:     *production,
		Disabled:       *disableTrace,
	})
	if err != nil {
		return fmt.Errorf("failed to init Jaeger: %v", err)
	}
	defer providerTracer.Close(ctx)
	log.Info(ctx, "Connected to Jaeger")

	//-------------------------Init consul-----------------------------------------
	consulClient, serviceID, err := consul.NewConsulClient(cfg)
	if err != nil {
		return fmt.Errorf("failed init consul: %v", err)
	}
}
