package http

import (
	"fmt"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/PavelDonchenko/ecommerce-micro/common/config"
	"github.com/PavelDonchenko/ecommerce-micro/common/database/sqldb"
	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
	"github.com/PavelDonchenko/ecommerce-micro/common/metric"
	middlewares "github.com/PavelDonchenko/ecommerce-micro/common/mid"
)

type Router struct {
	log            *logger.Logger
	cfg            *config.Config
	auth           *Auth
	db             *sqlx.DB
	serviceMetrics metric.Metrics
}

func NewRouter(log *logger.Logger, cfg *config.Config, auth *Auth, db *sqlx.DB, sm metric.Metrics) *Router {
	return &Router{
		log:            log,
		cfg:            cfg,
		auth:           auth,
		db:             db,
		serviceMetrics: sm,
	}
}

func (r *Router) SetupRouter(prod bool) *gin.Engine {
	transaction := middlewares.BeginCommitRollback(r.log, sqldb.NewBeginner(r.db))
	router := r.initRoute(prod)

	router.GET("/health", middlewares.Healthy())
	router.GET("/metrics", middlewares.MetricsHandler())

	v1 := router.Group(fmt.Sprintf("/api/%s", r.cfg.APIVersion))

	v1.Use(gin.Recovery())
	v1.Use(middlewares.CORS())
	v1.Use(location.Default())
	v1.Use(otelgin.Middleware(r.cfg.Jaeger.ServiceName))
	v1.Use(middlewares.Metrics(r.serviceMetrics))
	v1.Use(middlewares.Logger(r.log))

	v1.POST("/signup", transaction, r.auth.SignUp)

	if !r.cfg.Production {
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return router
}

func (r *Router) initRoute(prod bool) *gin.Engine {
	if prod {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	return gin.New()
}
