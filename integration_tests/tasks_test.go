package integration

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"
)

func TestTasksIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{
		ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
		Client: &http.Client{Transport: r},
	}

	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Tasks")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	cus1, err := api.CreateCustomer(&cm.NewCustomer{
		Name:           "Test Tasks",
		ExternalID:     "ext_customer_1",
		DataSourceUUID: ds.UUID,
	})
	if err != nil {
		t.Fatal(err)
	}

	newTaskParams := &cm.NewTask{
		CustomerUUID: cus1.UUID,
		Assignee:     "keith+test1@chartmogul.com",
		TaskDetails:  "This is some task details text.",
		DueDate:      "2025-04-30T00:00:00Z",
		CompletedAt:  "2025-04-20T00:00:00Z",
	}
	newTask, err := api.CreateTask(newTaskParams)
	if err != nil {
		t.Fatal(err)
	}
	allTasks, err := api.ListTasks(&cm.ListTasksParams{
		CustomerUUID: cus1.UUID,
		Cursor:       cm.Cursor{PerPage: 10},
	})
	if err != nil {
		t.Fatal(err)
	}

	var expectedAllTasks *cm.Tasks = &cm.Tasks{
		Entries: []*cm.Task{newTask},
	}
	expectedAllTasks.Cursor = allTasks.Cursor
	expectedAllTasks.HasMore = false

	if !reflect.DeepEqual(allTasks, expectedAllTasks) {
		spew.Dump(allTasks)
		spew.Dump(expectedAllTasks)
		t.Fatal("All tasks are not equal!")
	}

	retrievedTask, err := api.RetrieveTask(newTask.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(retrievedTask, newTask) {
		spew.Dump(retrievedTask)
		t.Fatal("Created task is not equal!")
	}

	updatedTaskParams := &cm.UpdateTask{
		TaskDetails: "This is some other task details text.",
	}
	updatedTask, err := api.UpdateTask(updatedTaskParams, retrievedTask.UUID)
	if err != nil {
		t.Fatal(err)
	}
	updatedRetrievedTask, err := api.RetrieveTask(updatedTask.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updatedTask, updatedRetrievedTask) {
		spew.Dump(updatedRetrievedTask)
		t.Fatal("Updated task is not equal!")
	}

	otherTaskParams := &cm.NewTask{
		CustomerUUID: cus1.UUID,
		Assignee:     "keith+test1@chartmogul.com",
		TaskDetails:  "This is details text for another task.",
		DueDate:      "2025-05-30T00:00:00Z",
		CompletedAt:  "2025-05-20T00:00:00Z",
	}
	otherTask, err := api.CreateTask(otherTaskParams)
	if err != nil {
		t.Fatal(err)
	}

	deleteErr := api.DeleteTask(otherTask.UUID)
	if deleteErr != nil {
		t.Fatal(err)
	}
}
