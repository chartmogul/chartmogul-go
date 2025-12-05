package chartmogul

import (
	"log"
	"testing"
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
	log.Println("Data source created.")

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
	log.Println("Data source found.")

	err = api.DeleteDataSource(ds.UUID)
	if err != nil {
		t.Error(err)
	}
	log.Println("Data source deleted.")
}

// TestListDataSourcesWithParams tests listing Data Sources with query parameters.
func TestListDataSourcesWithParams(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	ds, err := api.CreateDataSource(dsTestName)
	if err != nil {
		t.Error(err)
	}
	defer api.DeleteDataSource(ds.UUID)

	trueBool := true
	falseBool := false

	// Test with all parameters
	params := &ListDataSourcesParams{
		Name:   dsTestName,
		ExtraDataSourceParams: ExtraDataSourceParams{
			WithProcessingStatus:             &trueBool,
			WithAutoChurnSubscriptionSetting: &trueBool,
			WithInvoiceHandlingSetting:       &trueBool,
		},
	}

	res, err := api.ListDataSourcesWithFilters(params)
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, dataSource := range res.DataSources {
		if dataSource.Name == dsTestName {
			found = true
			// When with_processing_status=true, ProcessingStatus should be included
			if params.WithProcessingStatus != nil && *params.WithProcessingStatus {
				if dataSource.ProcessingStatus == nil {
					t.Errorf("Expected ProcessingStatus to be included when with_processing_status=true")
				}
			}
			// When with_auto_churn_subscription_setting=true, AutoChurnSubscriptionSetting should be included
			if params.WithAutoChurnSubscriptionSetting != nil && *params.WithAutoChurnSubscriptionSetting {
				if dataSource.AutoChurnSubscriptionSetting == nil {
					t.Errorf("Expected AutoChurnSubscriptionSetting to be included when with_auto_churn_subscription_setting=true")
				}
			}
			// When with_invoice_handling_setting=true, InvoiceHandlingSetting should be included
			if params.WithInvoiceHandlingSetting != nil && *params.WithInvoiceHandlingSetting {
				if dataSource.InvoiceHandlingSetting == nil {
					t.Errorf("Expected InvoiceHandlingSetting to be included when with_invoice_handling_setting=true")
				}
			}
			break
		}
	}

	if !found {
		t.Errorf("Data source not found in listing with filters! %+v", res)
	}

	// Test with parameters set to false
	params.WithProcessingStatus = &falseBool
	params.WithAutoChurnSubscriptionSetting = &falseBool
	params.WithInvoiceHandlingSetting = &falseBool

	res, err = api.ListDataSourcesWithFilters(params)
	if err != nil {
		t.Error(err)
	}

	found = false
	for _, dataSource := range res.DataSources {
		if dataSource.Name == dsTestName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Data source not found in listing with false filters! %+v", res)
	}

	log.Println("Data sources listing with parameters tested successfully.")
}

// TestRetrieveDataSourceWithParams tests retrieving a Data Source with query parameters.
func TestRetrieveDataSourceWithParams(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	ds, err := api.CreateDataSource(dsTestName)
	if err != nil {
		t.Error(err)
	}
	defer api.DeleteDataSource(ds.UUID)

	trueBool := true
	falseBool := false

	// Test with all parameters set to true
	params := &ExtraDataSourceParams{
		WithProcessingStatus:             &trueBool,
		WithAutoChurnSubscriptionSetting: &trueBool,
		WithInvoiceHandlingSetting:       &trueBool,
	}

	retrieved, err := api.RetrieveDataSource(ds.UUID, params)
	if err != nil {
		t.Error(err)
	}

	if retrieved.UUID != ds.UUID {
		t.Errorf("Data source UUIDs don't match - expected: %v, actual: %v", ds.UUID, retrieved.UUID)
	}

	// When with_processing_status=true, ProcessingStatus should be included
	if retrieved.ProcessingStatus == nil {
		t.Errorf("Expected ProcessingStatus to be included when with_processing_status=true")
	}

	// When with_auto_churn_subscription_setting=true, AutoChurnSubscriptionSetting should be included
	if retrieved.AutoChurnSubscriptionSetting == nil {
		t.Errorf("Expected AutoChurnSubscriptionSetting to be included when with_auto_churn_subscription_setting=true")
	}

	// When with_invoice_handling_setting=true, InvoiceHandlingSetting should be included
	if retrieved.InvoiceHandlingSetting == nil {
		t.Errorf("Expected InvoiceHandlingSetting to be included when with_invoice_handling_setting=true")
	}

	// Test with parameters set to false
	params = &ExtraDataSourceParams{
		WithProcessingStatus:             &falseBool,
		WithAutoChurnSubscriptionSetting: &falseBool,
		WithInvoiceHandlingSetting:       &falseBool,
	}

	retrieved, err = api.RetrieveDataSource(ds.UUID, params)
	if err != nil {
		t.Error(err)
	}

	if retrieved.UUID != ds.UUID {
		t.Errorf("Data source UUIDs don't match - expected: %v, actual: %v", ds.UUID, retrieved.UUID)
	}

	// Test with partial parameters - only processing status
	params = &ExtraDataSourceParams{
		WithProcessingStatus: &trueBool,
	}

	retrieved, err = api.RetrieveDataSource(ds.UUID, params)
	if err != nil {
		t.Error(err)
	}

	if retrieved.UUID != ds.UUID {
		t.Errorf("Data source UUIDs don't match - expected: %v, actual: %v", ds.UUID, retrieved.UUID)
	}

	// Test with mixed parameters
	params = &ExtraDataSourceParams{
		WithProcessingStatus:       &trueBool,
		WithInvoiceHandlingSetting: &falseBool,
	}

	retrieved, err = api.RetrieveDataSource(ds.UUID, params)
	if err != nil {
		t.Error(err)
	}

	if retrieved.UUID != ds.UUID {
		t.Errorf("Data source UUIDs don't match - expected: %v, actual: %v", ds.UUID, retrieved.UUID)
	}

	// Test backward compatibility - no parameters (can be omitted now)
	retrieved, err = api.RetrieveDataSource(ds.UUID)
	if err != nil {
		t.Error(err)
	}

	if retrieved.UUID != ds.UUID {
		t.Errorf("Data source UUIDs don't match - expected: %v, actual: %v", ds.UUID, retrieved.UUID)
	}

	// Test with nil parameters
	retrieved, err = api.RetrieveDataSource(ds.UUID, nil)
	if err != nil {
		t.Error(err)
	}

	if retrieved.UUID != ds.UUID {
		t.Errorf("Data source UUIDs don't match - expected: %v, actual: %v", ds.UUID, retrieved.UUID)
	}

	log.Println("Data source retrieval with parameters tested successfully.")
}

// TestListDataSourcesWithExtraParams tests listing Data Sources with ExtraDataSourceParams.
func TestListDataSourcesWithExtraParams(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	ds, err := api.CreateDataSource(dsTestName)
	if err != nil {
		t.Error(err)
	}
	defer api.DeleteDataSource(ds.UUID)

	trueBool := true
	falseBool := false

	// Test with all parameters
	params := &ExtraDataSourceParams{
		WithProcessingStatus:             &trueBool,
		WithAutoChurnSubscriptionSetting: &trueBool,
		WithInvoiceHandlingSetting:       &trueBool,
	}

	res, err := api.ListDataSources(params)
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, dataSource := range res.DataSources {
		if dataSource.Name == dsTestName {
			found = true
			// When with_processing_status=true, ProcessingStatus should be included
			if params.WithProcessingStatus != nil && *params.WithProcessingStatus {
				if dataSource.ProcessingStatus == nil {
					t.Errorf("Expected ProcessingStatus to be included when with_processing_status=true")
				}
			}
			// When with_auto_churn_subscription_setting=true, AutoChurnSubscriptionSetting should be included
			if params.WithAutoChurnSubscriptionSetting != nil && *params.WithAutoChurnSubscriptionSetting {
				if dataSource.AutoChurnSubscriptionSetting == nil {
					t.Errorf("Expected AutoChurnSubscriptionSetting to be included when with_auto_churn_subscription_setting=true")
				}
			}
			// When with_invoice_handling_setting=true, InvoiceHandlingSetting should be included
			if params.WithInvoiceHandlingSetting != nil && *params.WithInvoiceHandlingSetting {
				if dataSource.InvoiceHandlingSetting == nil {
					t.Errorf("Expected InvoiceHandlingSetting to be included when with_invoice_handling_setting=true")
				}
			}
			break
		}
	}

	if !found {
		t.Errorf("Data source not found in listing with extra params! %+v", res)
	}

	// Test with parameters set to false
	params.WithProcessingStatus = &falseBool
	params.WithAutoChurnSubscriptionSetting = &falseBool
	params.WithInvoiceHandlingSetting = &falseBool

	res, err = api.ListDataSources(params)
	if err != nil {
		t.Error(err)
	}

	found = false
	for _, dataSource := range res.DataSources {
		if dataSource.Name == dsTestName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Data source not found in listing with false extra params! %+v", res)
	}

	log.Println("Data sources listing with ExtraDataSourceParams tested successfully.")
}
