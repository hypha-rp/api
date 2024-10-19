package logging

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// logErrorAndRespond logs an error message and sends a JSON response with an internal server error status.
// Parameters:
// - context: The Gin context to use for sending the response.
// - message: The error message to log.
// - err: The error to log.
func HttpLogErrorAndRespond(context *gin.Context, log zerolog.Logger, message string, err error) {
	log.Error().Msgf("%s: %v", message, err)
	context.JSON(http.StatusInternalServerError, gin.H{"error": "There was a problem processing your request"})
}
