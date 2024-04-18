package testcasefailure

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupTestCaseFailureRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/test_case_failure", func(context *gin.Context) {
		CreateTestCaseFailure(dbOperations, context)
	})
	router.GET("/test_case_failure/:id", func(context *gin.Context) {
		GetTestCaseFailure(dbOperations, context)
	})
}

func CreateTestCaseFailure(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newTestCaseFailure tables.TestCaseFailure
	utils.CreateResource(dbOperations, context, &newTestCaseFailure)
}

func GetTestCaseFailure(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingTestCaseFailure tables.TestCaseFailure
	utils.GetResource(dbOperations, context, &existingTestCaseFailure, "id", "TestCaseFailure")
}
