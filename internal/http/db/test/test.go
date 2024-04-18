package test

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupTestRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/test", func(context *gin.Context) {
		CreateTest(dbOperations, context)
	})
	router.GET("/test/:id", func(context *gin.Context) {
		GetTest(dbOperations, context)
	})
}

func CreateTest(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newTest tables.Test
	if err := context.ShouldBindJSON(&newTest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbOperations.Create(&newTest); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, newTest)
}

func GetTest(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingTest tables.Test
	testID := context.Param("id")

	if err := dbOperations.First(&existingTest, "id = ?", testID); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}

	context.JSON(http.StatusOK, existingTest)
}
