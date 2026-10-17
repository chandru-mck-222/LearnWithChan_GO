package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RequestAnalyser struct {
	Route        string
	RequestId    string
	Method       string
	Status       int
	StartTime    time.Time
	Latency      time.Duration
	ErrorMessage []string
}

func Analystic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		errors := []string{}
		latency := time.Since(start)
		if len(ctx.Errors) > 0 {
			errors = append(errors, ctx.Errors.Errors()...)
		}
		body := RequestAnalyser{
			Route:        ctx.FullPath(),
			RequestId:    ctx.GetHeader("X-Request-ID"),
			Status:       ctx.Writer.Status(),
			ErrorMessage: errors,
			Method:       ctx.Request.Method,
			StartTime:    start,
			Latency:      latency,
		}

		Condlog(
			"Endpoint analytics: Route=%s RequestID=%s Status=%d Method=%s Start=%s Latency=%s Errors=%v",
			body.Route,
			body.RequestId,
			body.Status,
			body.Method,
			body.StartTime.Format(time.RFC3339),
			body.Latency,
			body.ErrorMessage,
		)
	}
}
