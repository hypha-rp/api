package integration

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func InitIntegrationRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/integration", func(context *gin.Context) {
		CreateIntegration(dbOperations, context)
	})
	router.GET("/integration/:id", func(context *gin.Context) {
		GetIntegration(dbOperations, context)
	})
}

func CreateIntegration(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newProduct tables.Integration
	utils.CreateResource(dbOperations, context, &newProduct)
}

func GetIntegration(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingProduct tables.Integration
	utils.GetResource(dbOperations, context, &existingProduct, "id", "Product")
}
