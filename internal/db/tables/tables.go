package tables

import "time"

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

type Result struct {
	ID           string      `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID    string      `json:"productID"`
	TestSuites   []TestSuite `gorm:"foreignKey:ResultID"`
	DateReported time.Time   `json:"dateReported"`
}

type TestSuite struct {
	ID         string     `gorm:"type:uuid;primaryKey" json:"id"`
	ResultID   string     `json:"resultID"`
	Name       string     `json:"name"`
	Tests      int        `json:"tests"`
	Failures   int        `json:"failures"`
	Errors     int        `json:"errors"`
	Skipped    int        `json:"skipped"`
	Time       float64    `json:"time"`
	TestCases  []TestCase `gorm:"foreignKey:TestSuiteID"`
	Properties []Property `gorm:"foreignKey:TestSuiteID"`
}

type TestCase struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"id"`
	TestSuiteID string     `json:"testSuiteID"`
	ClassName   string     `json:"className"`
	Name        string     `json:"name"`
	Time        float64    `json:"time"`
	Status      string     `json:"status"`
	Message     *string    `json:"message"`
	Properties  []Property `gorm:"foreignKey:TestCaseID"`
}

type Property struct {
	ID          string  `gorm:"type:uuid;primaryKey" json:"id"`
	TestSuiteID *string `json:"testSuiteID"`
	TestCaseID  *string `json:"testCaseID"`
	Name        string  `json:"name"`
	Value       string  `json:"value"`
}
