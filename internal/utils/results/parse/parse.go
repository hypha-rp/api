package parse

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/utils/results/structs"
)

func ParseJUnitResults(testSuites structs.JUnitTestSuites, dbOperations ops.DatabaseOperations, productId string) error {
	for _, suite := range testSuites.TestSuites {
		resultModel, err := createResultModel(productId)
		if err != nil {
			return err
		}
		if err := dbOperations.Create(&resultModel); err != nil {
			return err
		}

		testSuiteModel, err := createTestSuiteModel(suite, resultModel.ID)
		if err != nil {
			return err
		}
		if err := dbOperations.Create(&testSuiteModel); err != nil {
			return err
		}

		if err := createAndSaveProperties(suite.Properties, testSuiteModel.ID, dbOperations); err != nil {
			return err
		}

		if err := createAndSaveTestCases(suite.TestCases, testSuiteModel.ID, dbOperations); err != nil {
			return err
		}
	}
	return nil
}
