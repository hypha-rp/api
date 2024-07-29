package tables

import (
	"github.com/go-orm/gorm"
)

type Product struct {
	gorm.Model
	FullName     string        `json:"fullName"`
	ShortName    string        `json:"shortName"`
	ContactEmail string        `json:"contactEmail"`
	Integrations []Integration `gorm:"foreignKey:ProductID1;foreignKey:ProductID2"`
}

type Integration struct {
	gorm.Model
	ProductID1 uint    `json:"productID1"`
	ProductID2 uint    `json:"productID2"`
	Product1   Product `gorm:"foreignKey:ProductID1"`
	Product2   Product `gorm:"foreignKey:ProductID2"`
}
