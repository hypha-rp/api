package repoconfig

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupRepoConfigRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/repo_config", func(context *gin.Context) {
		CreateRepoConfig(dbOperations, context)
	})
	router.GET("/repo_config/:id", func(context *gin.Context) {
		GetRepoConfig(dbOperations, context)
	})
}

func CreateRepoConfig(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newRepoConfig tables.RepoConfig
	utils.CreateResource(dbOperations, context, &newRepoConfig)
}

func GetRepoConfig(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingRepoConfig tables.RepoConfig
	utils.GetResource(dbOperations, context, &existingRepoConfig, "id", "RepoConfig")
}
