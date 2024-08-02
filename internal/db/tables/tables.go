package tables

import (
	"github.com/go-orm/gorm"
)

type Product struct {
	ID           string        `gorm:"type:uuid;primaryKey" json:"id"`
	FullName     string        `json:"fullName"`
	ShortName    string        `json:"shortName"`
	ContactEmail string        `json:"contactEmail"`
	Integrations []Integration `gorm:"foreignKey:ProductID1;foreignKey:ProductID2"`
}

type Integration struct {
	ID         string  `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID1 string  `gorm:"type:uuid;" json:"productID1"`
	ProductID2 string  `gorm:"type:uuid" json:"productID2"`
	Product1   Product `gorm:"foreignKey:ProductID1"`
	Product2   Product `gorm:"foreignKey:ProductID2"`
}

type Assembly struct {
	ID            string       `gorm:"primaryKey" json:"id"`
	Name          string       `json:"name"`
	TestFramework string       `json:"testFramework"`
	RunDate       string       `json:"runDate"`
	RunTime       string       `json:"runTime"`
	Total         int          `json:"total"`
	Passed        int          `json:"passed"`
	Failed        int          `json:"failed"`
	Skipped       int          `json:"skipped"`
	Time          float64      `json:"time"`
	ProductID     uint         `json:"productId"`
	Product       Product      `gorm:"foreignKey:ProductID"`
	Collections   []Collection `gorm:"foreignKey:AssemblyID"`
}

type Collection struct {
	ID         string `gorm:"primaryKey" json:"id"`
	AssemblyID string `json:"assemblyID"`
	Total      int    `json:"total"`
	Passed     int    `json:"passed"`
	Failed     int    `json:"failed"`
	Skipped    int    `json:"skipped"`
	Name       string `json:"name"`
	Tests      []Test `gorm:"foreignKey:CollectionID"`
}

type Test struct {
	ID           string  `gorm:"primaryKey" json:"id"`
	CollectionID string  `json:"collectionID"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Method       string  `json:"method"`
	Time         float64 `json:"time"`
	Result       string  `json:"result"`
	Traits       []Trait `gorm:"foreignKey:TestID"`
}

type Trait struct {
	gorm.Model
	TestID string `json:"testID"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}
