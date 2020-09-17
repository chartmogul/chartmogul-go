package main
import (
	"fmt"
	cm "github.com/chartmogul/chartmogul-go"
)

func main() {
api := &cm.API{
AccountToken: "xxx",
AccessKey:    "yyy",
}


	_, err := api.UploadCSVFile("./invoices.csv", &cm.CsvUploadRequest{
DataSourceUUID: "ds_uuid",
DataType:       "invoice",
BatchName:      "Invoices Upload",
})
fmt.Print(err)
}
