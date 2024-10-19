package logging

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

// init initializes the global Logger with zerolog settings.
// The logger is configured to output to stderr with a timestamp and trace level logging.
func init() {
	Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Logger()
}

// GinLogger returns a gin.HandlerFunc (middleware) that logs HTTP requests using zerolog.
// It logs the status code, client IP, latency, method, and path of each request.
//
// Returns:
//   - gin.HandlerFunc: A middleware handler function for logging HTTP requests.
func GinLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now()

		context.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		statusCode := context.Writer.Status()
		clientIP := context.ClientIP()
		method := context.Request.Method
		path := context.Request.URL.Path

		Logger.Info().
			Int("status", statusCode).
			Str("client_ip", clientIP).
			Str("latency", latency.String()).
			Str("method", method).
			Str("path", path).
			Msg("REQUEST")
	}
}
