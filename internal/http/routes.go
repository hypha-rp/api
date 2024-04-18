package http

import (
	"hypha/api/internal/http/db/product"
	"hypha/api/internal/http/db/productsrepo"
	"hypha/api/internal/http/db/repo"
	"hypha/api/internal/http/db/repoconfig"
	"hypha/api/internal/http/db/repoconfigrule"
	"hypha/api/internal/http/db/testcasefailure"
	"hypha/api/internal/http/db/testcaseresult"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, dbOperations utils.DatabaseOperations) {
	dbGroup := router.Group("/db")
	product.SetupProductRoutes(dbGroup, dbOperations)
	productsrepo.SetupProductsRepoRoutes(dbGroup, dbOperations)
	repo.SetupRepoRoutes(dbGroup, dbOperations)
	repoconfig.SetupRepoConfigRoutes(dbGroup, dbOperations)
	repoconfigrule.SetupRepoConfigRuleRoutes(dbGroup, dbOperations)
	testcasefailure.SetupTestCaseFailureRoutes(dbGroup, dbOperations)
	testcaseresult.SetupTestCaseResultRoutes(dbGroup, dbOperations)
}
