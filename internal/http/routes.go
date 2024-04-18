package http

import (
	"hypha/api/internal/http/db/product"
	"hypha/api/internal/http/db/repo"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, dbOperations utils.DatabaseOperations) {
	dbGroup := router.Group("/db")
	product.SetupProductRoutes(dbGroup, dbOperations)
	repo.SetupTestRoutes(dbGroup, dbOperations)
}
