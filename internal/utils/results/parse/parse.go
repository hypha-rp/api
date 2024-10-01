package parse

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils/results/structs"
	"time"
)

func ParseJUnitResults(testSuites structs.JUnitTestSuites, dbOperations ops.DatabaseOperations, productId string) error {
	for _, suite := range testSuites.TestSuites {
		resultModel := tables.Result{
			ID:           ops.GenerateUniqueID(),
			ProductID:    productId,
			DateReported: time.Now().UTC(),
		}
		if err := dbOperations.Create(&resultModel); err != nil {
			return err
		}

		testSuiteModel := tables.TestSuite{
			ID:         ops.GenerateUniqueID(),
			ResultID:   resultModel.ID,
			Name:       suite.Name,
			Tests:      suite.Tests,
			Failures:   suite.Failures,
			Errors:     suite.Errors,
			Skipped:    suite.Skipped,
			Assertions: suite.Assertions,
			Time:       suite.Time,
			File:       suite.File,
			SystemOut:  suite.SystemOut,
			SystemErr:  suite.SystemErr,
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
			status := "pass"
			var message *string
			var testCaseType *string

			if testCase.Failure != nil {
				status = "fail"
				message = &testCase.Failure.Message
				testCaseType = &testCase.Failure.Type
			} else if testCase.Error != nil {
				status = "error"
				message = &testCase.Error.Message
				testCaseType = &testCase.Error.Type
			} else if testCase.Skipped != nil {
				status = "skipped"
				message = &testCase.Skipped.Message
				testCaseType = nil
			}

			testCaseModel := tables.TestCase{
				ID:          ops.GenerateUniqueID(),
				TestSuiteID: testSuiteModel.ID,
				ClassName:   testCase.ClassName,
				Name:        testCase.Name,
				Time:        testCase.Time,
				Status:      status,
				Message:     message,
				Type:        testCaseType,
				Assertions:  testCase.Assertions,
				File:        testCase.File,
				Line:        testCase.Line,
				SystemOut:   testCase.SystemOut,
				SystemErr:   testCase.SystemErr,
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
