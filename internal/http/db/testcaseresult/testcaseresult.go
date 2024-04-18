package testcaseresult

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupTestCaseResultRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/test_case_result", func(context *gin.Context) {
		CreateTestCaseResult(dbOperations, context)
	})
	router.GET("/test_case_result/:id", func(context *gin.Context) {
		GetTestCaseResult(dbOperations, context)
	})
}

func CreateTestCaseResult(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newTestCaseResult tables.TestCaseResult
	utils.CreateResource(dbOperations, context, &newTestCaseResult)
}

func GetTestCaseResult(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingTestCaseResult tables.TestCaseResult
	utils.GetResource(dbOperations, context, &existingTestCaseResult, "id", "TestCaseResult")
}
