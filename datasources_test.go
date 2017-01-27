package chartmogul

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

const dsTestName = "some name"

// TestImportDataSource tests creation, listing & deletion of Data Sources.
func TestImportDataSources(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	ds, err := api.CreateDataSource(dsTestName)
	if err != nil {
		t.Error(err)
	} else if ds.Name != dsTestName {
		t.Errorf("Data source names don't equal - expected: %v, actual: %v", dsTestName, ds.Name)
	} else if ds.UUID == "" {
		t.Errorf("Data source has no UUID!")
	} else if ds.CreatedAt == "" || ds.Status == "" {
		t.Errorf("Data source has empty attributes! %+v", ds)
	}
	logrus.Debug("Data source created.")

	res, err := api.ListDataSources()
	if err != nil {
		t.Error(err)
	}
	found := false
	for _, ds := range res.DataSources {
		if ds.Name == dsTestName {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Data source not found in listing! %+v", res)
	}
	logrus.Debug("Data source found.")

	err = api.DeleteDataSource(ds.UUID)
	if err != nil {
		t.Error(err)
	}
	logrus.Debug("Data source deleted.")
}
