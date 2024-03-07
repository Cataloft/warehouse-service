package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"lamoda_task/internal/tools/duration"

	"github.com/gin-gonic/gin"
	requestid "github.com/sumit-tembe/gin-requestid"
)

func LogMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		dur := duration.GetDurationInMilliseconds(start)
		entry := log.With(
			"client_ip", c.ClientIP(),
			"duration", fmt.Sprintf("%.3f%s", dur, "ms"),
			"method", c.Request.Method,
			"path", c.Request.RequestURI,
			"status", c.Writer.Status(),
			"request_id", requestid.GetRequestIDFromHeaders(c),
		)

		if c.Writer.Status() >= http.StatusBadRequest {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("Success")
		}
	}
}
