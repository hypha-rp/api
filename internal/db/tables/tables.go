package tables

import (
	"github.com/go-orm/gorm"
)

type Product struct {
	gorm.Model
	FullName      string         `json:"fullName"`
	ShortName     string         `json:"shortName"`
	ContactEmail  string         `json:"contactEmail"`
	ProductsRepos []ProductsRepo `gorm:"foreignKey:ProductID"`
}

type Repo struct {
	gorm.Model
	Name           string           `json:"name"`
	Url            string           `json:"url"`
	ProductsRepos  []ProductsRepo   `gorm:"foreignKey:RepoID"`
	RepoConfig     []RepoConfig     `gorm:"foreignKey:RepoID"`
	TestCaseResult []TestCaseResult `gorm:"foreignKey:RepoID"`
}

type ProductsRepo struct {
	gorm.Model
	ProductID      uint `json:"productId"`
	RepoID         uint `json:"repoId"`
	PrimaryProduct bool `json:"primaryProduct"`
	Product        Product
	Repo           Repo
}

type RepoConfig struct {
	gorm.Model
	RepoID         uint             `json:"repoId"`
	TestType       *string          `json:"testType"`
	RepoConfigRule []RepoConfigRule `gorm:"foreignKey:RepoConfigID"`
	Repo           Repo
}

type RepoConfigRule struct {
	gorm.Model
	RepoConfigID uint   `json:"repoConfigId"`
	Rule         string `json:"rule"`
	RepoConfig   RepoConfig
}

type TestCaseResult struct {
	gorm.Model
	RepoID          uint    `json:"repoId"`
	TestSuiteName   string  `json:"testSuiteName"`
	TestCaseName    string  `json:"testCaseName"`
	Failure         bool    `json:"failure"`
	Skipped         bool    `json:"skipped"`
	ExecutionTime   float64 `json:"executionTime"`
	Repo            Repo
	TestCaseFailure []TestCaseFailure `gorm:"foreignKey:TestCaseResultID"`
}

type TestCaseFailure struct {
	gorm.Model
	TestCaseResultID uint   `json:"testCaseResultId"`
	FailureMessage   string `json:"failureMessage"`
	FailureType      string `json:"failureType"`
	TestCaseFailure  TestCaseResult
}
