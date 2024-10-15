package integration

import (
	"hypha/api/internal/db"

	"github.com/gin-gonic/gin"
)

func InitIntegrationRoutes(router *gin.RouterGroup, dbOperations db.DatabaseOperations) {
	router.POST("/integration", func(context *gin.Context) {
		CreateIntegration(dbOperations, context)
	})
	router.GET("/integration/:id", func(context *gin.Context) {
		GetIntegration(dbOperations, context)
	})
}

func CreateIntegration(dbOperations db.DatabaseOperations, context *gin.Context) {
	var newIntegration db.Integration
	newIntegration.ID = db.GenerateUniqueID()
	db.CreateResource(dbOperations, context, &newIntegration)
}

func GetIntegration(dbOperations db.DatabaseOperations, context *gin.Context) {
	var existingIntegration db.Integration
	if err := dbOperations.Connection().
		Preload("Product1").
		Preload("Product2").
		Where("id = ?", context.Param("id")).
		First(&existingIntegration).Error; err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, existingIntegration)
}
