package handlers

import (
	"github.com/dejandjenic/go-gin-sample/application/counters"
	"github.com/dejandjenic/go-gin-sample/application/tracing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Test struct {
	Message string `json:"message"`
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} handlers.Test
// @Router /ping [get]
func (h *Handler) Ping(g *gin.Context) {
	counters.PingCounter.Inc()

	ctx := g.Request.Context()
	ctx, span := tracing.Tracer.Start(ctx, "ping")
	defer span.End()
	log.Info().
		Str("traceid", span.SpanContext().TraceID().String()).
		Str("spanid", span.SpanContext().SpanID().String()).
		Msg("ping")

	g.JSON(200,
		Test{
			Message: "pong",
		},
	)
}
