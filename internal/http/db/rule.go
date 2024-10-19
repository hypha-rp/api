package db

import (
	"hypha/api/internal/db"

	"github.com/gin-gonic/gin"
)

func InitRuleRoutes(router *gin.RouterGroup, dbOperations db.DatabaseOperations) {
	router.POST("/results-rule", func(context *gin.Context) {
		CreateResultsRule(dbOperations, context)
	})
	router.GET("/results-rule/:id", func(context *gin.Context) {
		GetResultsRule(dbOperations, context)
	})
}

func CreateResultsRule(dbOperations db.DatabaseOperations, context *gin.Context) {
	var newRule db.ResultsRule
	newRule.ID = db.GenerateUniqueID()
	db.CreateResource(dbOperations, context, &newRule)
}

func GetResultsRule(dbOperations db.DatabaseOperations, context *gin.Context) {
	var existingRule db.ResultsRule
	db.GetResource(dbOperations, context, &existingRule, "id", "ResultsRule")
}

func GetResultsRuleByRelationID(dbOperations db.DatabaseOperations, relationID string) (*db.ResultsRule, error) {
	var rule db.ResultsRule
	err := dbOperations.Connection().Where("relation_id = ?", relationID).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}
