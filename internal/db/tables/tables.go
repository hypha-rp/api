package tables

import (
	"github.com/go-orm/gorm"
)

type Product struct {
	gorm.Model
	FullName     string `json:"fullName"`
	ShortName    string `json:"shortName"`
	ContactEmail string `json:"contactEmail"`
}

type Test struct {
	gorm.Model
	Name    string `json:"name"`
	RepoURL string `json:"repoURL"`
}
