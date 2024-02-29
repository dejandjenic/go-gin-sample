package middleware

import (
	"github.com/dejandjenic/go-gin-sample/application/tracing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func TraceHandler(c *gin.Context) {

	ctx := c.Request.Context()
	ctx, span := tracing.Tracer.Start(ctx, "root")
	defer span.End()

	log.
		Info().
		Str("traceid", span.SpanContext().TraceID().String()).
		Str("spanid", span.SpanContext().SpanID().String()).
		Msg("trace middleware")

	c.Writer.Header().Set("traceid", span.SpanContext().TraceID().String())

}
