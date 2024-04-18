package repoconfigrule

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupRepoConfigRuleRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/repo_config_rule", func(context *gin.Context) {
		CreateRepoConfigRule(dbOperations, context)
	})
	router.GET("/repo_config_rule/:id", func(context *gin.Context) {
		GetRepoConfigRule(dbOperations, context)
	})
}

func CreateRepoConfigRule(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newRepoConfigRule tables.RepoConfigRule
	utils.CreateResource(dbOperations, context, &newRepoConfigRule)
}

func GetRepoConfigRule(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingRepoConfigRule tables.RepoConfigRule
	utils.GetResource(dbOperations, context, &existingRepoConfigRule, "id", "RepoConfigRule")
}
