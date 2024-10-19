package http

import (
	"hypha/api/internal/db"
	"hypha/api/internal/http/routes"
	"hypha/api/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

var log = logging.Logger

// InitRoutes initializes all the routes for the given router engine.
// It sets up the database-related routes and results-related routes.
//
// Parameters:
// - router: The Gin engine to which the routes will be added.
// - dbOps: The database operations interface used for database interactions.
func InitRoutes(router *gin.Engine, dbOps db.DatabaseOperations) {
	log.Info().Msg("Initializing routes")

	dbGroup := router.Group("/db")
	routes.InitProductRoutes(dbGroup, dbOps)
	routes.InitRelationRoutes(dbGroup, dbOps)
	routes.InitRuleRoutes(dbGroup, dbOps)

	resultsGroup := router.Group("/results")
	routes.InitResultsRoutes(resultsGroup, dbOps)

	log.Info().Msg("Routes initialized")
}
