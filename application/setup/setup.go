package setup

import (
	"context"
	"strings"

	"github.com/dejandjenic/go-gin-sample/application"
	"github.com/dejandjenic/go-gin-sample/application/counters"
	"github.com/dejandjenic/go-gin-sample/application/handlers"
	"github.com/dejandjenic/go-gin-sample/application/middleware"
	"github.com/dejandjenic/go-gin-sample/application/tracing"
	docs "github.com/dejandjenic/go-gin-sample/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func GinHandler(oidcUrl string) (*gin.Engine, func(context.Context) error) {
	r := gin.New()
	cleanup := tracing.InitTracing()
	prometheus.MustRegister(counters.PingCounter)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Initialize Gin
	//gin.SetMode(gin.ReleaseMode)

	var app application.Application = application.NewApplication(context.Background(), oidcUrl)
	log.Info().Msg("zerolog start")

	r.Use(otelgin.Middleware(tracing.ServiceName))
	p := ginprometheus.NewPrometheus("gin")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "id" {
				url = strings.Replace(url, p.Value, ":id", 1)
				break
			}
		}
		return url
	}
	p.Use(r)
	r.Use(middleware.TraceHandler)
	r.Use(middleware.AuthHandler(&app.Configuration))
	r.Use(middleware.DefaultStructuredLogger()) // adds our new middleware
	r.Use(gin.Recovery())                       // adds the default recovery middleware

	docs.SwaggerInfo.BasePath = "/api/v1"
	eg := r.Group("/api/v1")
	{
		h := handlers.ToHandler(&app)
		eg.GET("/ping", h.Ping)
		eg.POST("/todos", h.CreateTodo)
		eg.GET("/todos", h.ListTodos)
		eg.DELETE("/todos/:id", h.DeleteTodo)
		eg.PUT("/todos/:id", h.UpdateTodo)
		eg.GET("/todos/:id", h.TodoDetail)

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r, cleanup
}
