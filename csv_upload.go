package chartmogul

import "strings"

// Customer is the customer as represented in the API.

type CsvUploadRequest struct {
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	DataType       string `json:"data_type,omitempty"`
	FilePath       string `json:"file_path,omitempty"`
	BatchName      string `json:"batch_name,omitempty"`
}

type CsvUploadResponse struct {
	ID              string  `json:"id,omitempty"`
	OriginalName    string  `json:"original_name,omitempty"`
	DataType        string  `json:"data_type,omitempty"`
	StoragePath     string  `json:"storage_path,omitempty"`
	PercentComplete float32 `json:"percent_complete,omitempty"`
	CreatedAt       string  `json:"created_at,omitempty"`
	UpdatedAt       string  `json:"updated_at,omitempty"`
	BatchName       string  `json:"batch_name,omitempty"`
}

const (
	uploadEndoint = "data_sources/:data_source_uuid/uploads"
)

func (api API) UploadCSVFile(filename string, uploadRequest *CsvUploadRequest) (*CsvUploadResponse, error) {
	result := &CsvUploadResponse{}

	path := strings.Replace(uploadEndoint, ":data_source_uuid", uploadRequest.DataSourceUUID, 1)

	return result, api.upload(path, filename, uploadRequest, result)
}
