package parse

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils/results/structs"
)

func ParseJUnitResults(testSuites structs.JUnitTestSuites, dbOperations ops.DatabaseOperations, productId string) error {
	for _, suite := range testSuites.TestSuites {
		resultModel := tables.Result{
			ID:        ops.GenerateUniqueID(),
			ProductID: productId,
		}
		if err := dbOperations.Create(&resultModel); err != nil {
			return err
		}

		testSuiteModel := tables.TestSuite{
			ID:       ops.GenerateUniqueID(),
			ResultID: resultModel.ID,
			Name:     suite.Name,
			Tests:    suite.Tests,
			Failures: suite.Failures,
			Errors:   suite.Errors,
			Skipped:  suite.Skipped,
			Time:     suite.Time,
		}
		if err := dbOperations.Create(&testSuiteModel); err != nil {
			return err
		}

		for _, property := range suite.Properties {
			propertyModel := tables.Property{
				ID:          ops.GenerateUniqueID(),
				TestSuiteID: &testSuiteModel.ID,
				Name:        property.Name,
				Value:       property.Value,
			}
			if err := dbOperations.Create(&propertyModel); err != nil {
				return err
			}
		}

		for _, testCase := range suite.TestCases {
			testCaseModel := tables.TestCase{
				ID:          ops.GenerateUniqueID(),
				TestSuiteID: testSuiteModel.ID,
				ClassName:   testCase.ClassName,
				Name:        testCase.Name,
				Time:        testCase.Time,
			}
			if err := dbOperations.Create(&testCaseModel); err != nil {
				return err
			}

			for _, property := range testCase.Properties {
				propertyModel := tables.Property{
					ID:         ops.GenerateUniqueID(),
					TestCaseID: &testCaseModel.ID,
					Name:       property.Name,
					Value:      property.Value,
				}
				if err := dbOperations.Create(&propertyModel); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
