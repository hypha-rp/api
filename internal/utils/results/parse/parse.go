package parse

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils/results/structs"
)

func ParseXUnitResults(assemblies structs.Assemblies, dbOperations ops.DatabaseOperations, productId string) error {
	for _, assembly := range assemblies.Assemblies {
		assemblyModel := tables.Assembly{
			ID:            assembly.ID,
			Name:          assembly.Name,
			TestFramework: assembly.TestFramework,
			RunDate:       assembly.RunDate,
			RunTime:       assembly.RunTime,
			Total:         assembly.Total,
			Passed:        assembly.Passed,
			Failed:        assembly.Failed,
			Skipped:       assembly.Skipped,
			Time:          assembly.Time,
			ProductID:     productId,
		}
		if err := dbOperations.Create(&assemblyModel); err != nil {
			return err
		}

		for _, collection := range assembly.Collections {
			collectionModel := tables.Collection{
				ID:         collection.ID,
				AssemblyID: assemblyModel.ID,
				Total:      collection.Total,
				Passed:     collection.Passed,
				Failed:     collection.Failed,
				Skipped:    collection.Skipped,
				Name:       collection.Name,
			}
			if err := dbOperations.Create(&collectionModel); err != nil {
				return err
			}

			for _, test := range collection.Tests {
				testModel := tables.Test{
					ID:           test.ID,
					CollectionID: collectionModel.ID,
					Name:         test.Name,
					Type:         test.Type,
					Method:       test.Method,
					Time:         test.Time,
					Result:       test.Result,
				}
				if err := dbOperations.Create(&testModel); err != nil {
					return err
				}

				for _, trait := range test.Traits {
					traitModel := tables.Trait{
						TestID: testModel.ID,
						Name:   trait.Name,
						Value:  trait.Value,
					}
					if err := dbOperations.Create(&traitModel); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
