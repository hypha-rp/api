package handlers

import (
	"hypha/api/internal/db"

	"github.com/gin-gonic/gin"
)

func CreateResultsRule(dbOps db.DatabaseOperations, context *gin.Context) {
	var newRule db.ResultsRule
	newRule.ID = db.GenerateUniqueID()
	db.CreateResource(dbOps, context, &newRule)
}

func GetResultsRule(dbOps db.DatabaseOperations, context *gin.Context) {
	var existingRule db.ResultsRule
	db.GetResource(dbOps, context, &existingRule, "id", "ResultsRule")
}

func GetResultsRuleByRelationID(dbOps db.DatabaseOperations, relationID string) (*db.ResultsRule, error) {
	var rule db.ResultsRule
	err := dbOps.Connection().Where("relation_id = ?", relationID).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}
