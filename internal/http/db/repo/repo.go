package repo

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupRepoRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/repo", func(context *gin.Context) {
		CreateRepo(dbOperations, context)
	})
	router.GET("/repo/:id", func(context *gin.Context) {
		GetRepo(dbOperations, context)
	})
}

func CreateRepo(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newRepo tables.Repo
	utils.CreateResource(dbOperations, context, &newRepo)
}

func GetRepo(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingRepo tables.Repo
	utils.GetResource(dbOperations, context, &existingRepo, "id", "Repo")
}
