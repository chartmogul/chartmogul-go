package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestListTasks(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks?per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"task_uuid": "00000000-0000-0000-0000-000000000000",
						"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
						"assignee": "keith+test1@chartmogul.com",
						"task_details": "This is some task details text.",
						"due_date": "2025-04-30T00:00:00Z",
						"completed_at": "2025-04-20T00:00:00Z",
						"created_at": "2025-04-01T12:00:00.000Z",
						"updated_at": "2025-04-01T12:00:00.000Z"
					}],
					"has_more": false,
					"cursor": "88abf99"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	params := &ListTasksParams{Cursor: Cursor{PerPage: 1}}
	tasks, err := tested.ListTasks(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(tasks.Entries) == 0 {
		spew.Dump(tasks)
		t.Fatal("Unexpected result")
	}
}

func TestListTasksWithCustomerUuid(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks?customer_uuid=cus_00000000-0000-0000-0000-000000000000&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"task_uuid": "00000000-0000-0000-0000-000000000000",
						"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
						"assignee": "keith+test1@chartmogul.com",
						"task_details": "This is some task details text.",
						"due_date": "2025-04-30T00:00:00Z",
						"completed_at": "2025-04-20T00:00:00Z",
						"created_at": "2025-04-01T12:00:00.000Z",
						"updated_at": "2025-04-01T12:00:00.000Z"
					}],
					"has_more": false,
					"cursor": "88abf99"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	params := &ListTasksParams{Cursor: Cursor{PerPage: 1}, CustomerUUID: "cus_00000000-0000-0000-0000-000000000000"}
	tasks, err := tested.ListTasks(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(tasks.Entries) == 0 {
		spew.Dump(tasks)
		t.Fatal("Unexpected result")
	}
}

func TestRetrieveTask(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks/00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"task_uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"assignee": "keith+test1@chartmogul.com",
					"task_details": "This is some task details text.",
					"due_date": "2025-04-30T00:00:00Z",
					"completed_at": "2025-04-20T00:00:00Z",
					"created_at": "2025-04-01T12:00:00.000Z",
					"updated_at": "2025-04-01T12:00:00.000Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	task, err := tested.RetrieveTask("00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if task == nil {
		spew.Dump(task)
		t.Fatal("Unexpected result")
	}
}

func TestCreateTask(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"task_uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"assignee": "keith+test1@chartmogul.com",
					"task_details": "This is some task details text.",
					"due_date": "2025-04-30T00:00:00Z",
					"completed_at": "2025-04-20T00:00:00Z",
					"created_at": "2025-04-01T12:00:00.000Z",
					"updated_at": "2025-04-01T12:00:00.000Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	task, err := tested.CreateTask(&NewTask{
		CustomerUUID: "cus_00000000-0000-0000-0000-000000000000",
		Assignee:     "keith+test1@chartmogul.com",
		TaskDetails:  "This is some task details text.",
		DueDate:      "2025-04-30T00:00:00Z",
		CompletedAt:  "2025-04-20T00:00:00Z",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if task.UUID != "00000000-0000-0000-0000-000000000000" {
		spew.Dump(task)
		t.Fatal("Unexpected result")
	}
}

func TestUpdateTask(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks/00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"task_uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"assignee": "keith+test1@chartmogul.com",
					"task_details": "This is some other task details text.",
					"due_date": "2025-04-30T00:00:00Z",
					"completed_at": "2025-04-20T00:00:00Z",
					"created_at": "2025-04-01T12:00:00.000Z",
					"updated_at": "2025-04-01T12:00:00.000Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	task, err := tested.UpdateTask(&UpdateTask{
		TaskDetails: "This is some other task details text.",
	}, "00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if task.TaskDetails != "This is some other task details text." {
		spew.Dump(task)
		t.Fatal("Unexpected result")
	}
}

func TestDeleteTask(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks/00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	err := tested.DeleteTask("00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
