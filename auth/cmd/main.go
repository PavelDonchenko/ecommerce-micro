package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	service_http "github.com/PavelDonchenko/ecommerce-micro/auth/api/http"
	"github.com/PavelDonchenko/ecommerce-micro/auth/repositories"
	"github.com/PavelDonchenko/ecommerce-micro/auth/security"
	"github.com/PavelDonchenko/ecommerce-micro/auth/service"
	common_cert "github.com/PavelDonchenko/ecommerce-micro/common/certificates"
	"github.com/PavelDonchenko/ecommerce-micro/common/config"
	"github.com/PavelDonchenko/ecommerce-micro/common/consul"
	"github.com/PavelDonchenko/ecommerce-micro/common/database/sqldb"
	"github.com/PavelDonchenko/ecommerce-micro/common/helpers"
	"github.com/PavelDonchenko/ecommerce-micro/common/httputils"
	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
	"github.com/PavelDonchenko/ecommerce-micro/common/metric"
	"github.com/PavelDonchenko/ecommerce-micro/common/tasks"
	"github.com/PavelDonchenko/ecommerce-micro/common/trace/otel/jaeger"
	"github.com/PavelDonchenko/ecommerce-micro/common/validation"
	"github.com/PavelDonchenko/ecommerce-micro/common/validator"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	production = flag.Bool("prod", false, "use -prod=true to run in production mode")
	disableTrace = flag.Bool("disable-trace", false, "use disable-trace=true if you want to disable tracing completely")
	flag.Parse()

	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return logger.GetTraceID(ctx)
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "AUTH", traceIDFn)
	// -------------------------------------------------------------------------

	cfg := config.MustLoadConfig(*production, "./auth/config/")

	log.Info(ctx, "Init config", slog.Bool("production mode", cfg.Production))

	validator.NewValidator("en")
	helpers.CreateFolders(cfg.Folders)

	if err := run(ctx, log, cfg); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger, cfg *config.Config) error {
	validation.New()
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
	log.Info(ctx, "Init consul", slog.String("ID", serviceID))
	//emailsServiceNameDone := make(chan bool)
	//go tasks.ReloadServiceName(
	//	ctx,
	//	cfg,
	//	consulClient,
	//	cfg.EmailService.ServiceName,
	//	consul.EmailService,
	//	emailsServiceNameDone,
	//)
	//<-emailsServiceNameDone

	//-----------------------------Init metric service----------------------
	metricSrv, err := metric.NewMetricsService(cfg)
	if err != nil {
		return fmt.Errorf("failed init metrics: %v", err)
	}
	log.Info(ctx, "Init metric")
	//---------------------------------Init database------------------------
	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.Postgres.User,
		Password:     cfg.Postgres.Password,
		Host:         cfg.Postgres.Host,
		Name:         cfg.Postgres.Name,
		MaxIdleConns: cfg.Postgres.MaxIdleConns,
		MaxOpenConns: cfg.Postgres.MaxOpenConns,
		DisableTLS:   cfg.Postgres.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("failed init database: %v", err)
	}
	defer db.Close()
	log.Info(ctx, "Init database")

	certificatesServiceCommon := common_cert.NewCertificatesService(cfg)
	certificatesService := service.NewCertificatesServices(cfg, certificatesServiceCommon)
	certificatesManager := security.NewManagerCertificates(cfg, certificatesService, certificatesServiceCommon)

	checkCertificates := tasks.NewCheckCertificatesTask(cfg, *certificatesManager)
	certsDone := make(chan bool)
	go checkCertificates.Start(ctx, certsDone)
	<-certsDone

	authRepo := repositories.NewUser(db)
	authSrv := service.NewAuth(authRepo, log)
	authHTTP := service_http.NewAuth(log, *authSrv)
	r := service_http.NewRouter(log, cfg, authHTTP, db, metricSrv)
	server := httputils.NewHttpServer(cfg, r.SetupRouter(*production), certificatesServiceCommon)

	serverErrors := make(chan error, 1)
	var svr *http.Server
	if *production {
		svr = server.RunTLSServer(serverErrors)
	} else {
		svr = server.RunUnsecuredServer(serverErrors)
	}

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error: %v", err)
	case <-ctx.Done():
		log.Info(ctx, "shutdown", "status", "shutdown started")
		defer log.Info(ctx, "shutdown", "status", "shutdown complete")

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err = consulClient.Agent().ServiceDeregister(serviceID); err != nil {
			log.Error(ctx, "failed to deregister auth service", "error", err)
		}

		if err = svr.Shutdown(ctx); err != nil {
			svr.Close()
			return fmt.Errorf("could not stop auth server gracefully: %w", err)
		}
	}

	return nil
}
