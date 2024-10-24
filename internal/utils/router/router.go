package router

import (
	"errors"
	"hypha/api/internal/config"
	"hypha/api/internal/utils/logging"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var log = logging.Logger

// validateCorsPolicy validates the CORS policy configuration.
// It checks that the configuration specifies at least one allowed origin, method, and header.
// It also ensures that origins contain '*' or include 'http://' or 'https://'.
//
// Parameters:
// - cfg: The configuration object containing the CORS policy.
//
// Returns:
// - error: An error if the CORS policy configuration is invalid.
func validateCorsPolicy(cfg *config.Config) error {
	if len(cfg.Http.CorsPolicy.AllowOrigins) == 0 {
		return errors.New("CORS policy must specify at least one allowed origin")
	}
	for _, origin := range cfg.Http.CorsPolicy.AllowOrigins {
		if origin != "*" && !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
			return errors.New("bad origin: origins must contain '*' or include http://,https://")
		}
	}
	if len(cfg.Http.CorsPolicy.AllowMethods) == 0 {
		return errors.New("CORS policy must specify at least one allowed method")
	}
	if len(cfg.Http.CorsPolicy.AllowHeaders) == 0 {
		return errors.New("CORS policy must specify at least one allowed header")
	}
	return nil
}

// InitRouter initializes the Gin router with the given configuration.
// It validates the CORS policy, sets up middleware, and configures CORS settings.
//
// Parameters:
// - cfg: The configuration object containing the CORS policy and other settings.
//
// Returns:
// - *gin.Engine: The initialized Gin engine.
// - error: An error if the CORS policy configuration is invalid.
func InitRouter(cfg *config.Config) (*gin.Engine, error) {
	log.Info().Msg("Initializing router")

	if err := validateCorsPolicy(cfg); err != nil {
		log.Error().Err(err).Msg("Invalid CORS policy configuration")
		return nil, err
	}
	log.Info().Msg("CORS policy configuration is valid")

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logging.GinLogger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Http.CorsPolicy.AllowOrigins,
		AllowMethods:     cfg.Http.CorsPolicy.AllowMethods,
		AllowHeaders:     cfg.Http.CorsPolicy.AllowHeaders,
		ExposeHeaders:    cfg.Http.CorsPolicy.ExposeHeaders,
		AllowCredentials: cfg.Http.CorsPolicy.AllowCredentials,
		MaxAge:           time.Duration(cfg.Http.CorsPolicy.MaxAge) * time.Second,
	}))

	log.Info().Msg("Router initialized successfully")
	return router, nil
}
