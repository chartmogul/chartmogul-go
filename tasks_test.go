package chartmogul

import (
	"encoding/json"
	"io"
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

func TestListTasksWithFilters(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				if query.Get("contact_uuid") != "con_00000000-0000-0000-0000-000000000000" ||
					query.Get("assignee") != "keith+test1@chartmogul.com" ||
					query.Get("due_date_on_or_after") != "2025-04-01T00:00:00Z" ||
					query.Get("due_date_on_or_before") != "2025-04-30T00:00:00Z" ||
					query.Get("completed") != "false" ||
					query.Get("per_page") != "1" {
					t.Errorf("Unexpected query %v", r.URL.RawQuery)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"task_uuid": "00000000-0000-0000-0000-000000000000",
						"customer_uuid": null,
						"associated_object": "contact",
						"associated_object_uuid": "con_00000000-0000-0000-0000-000000000000",
						"assignee": "keith+test1@chartmogul.com",
						"task_details": "This is some task details text.",
						"due_date": "2025-04-30T00:00:00Z",
						"completed_at": null,
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
	completed := false
	params := &ListTasksParams{
		Cursor:            Cursor{PerPage: 1},
		ContactUUID:       "con_00000000-0000-0000-0000-000000000000",
		Assignee:          "keith+test1@chartmogul.com",
		DueDateOnOrAfter:  "2025-04-01T00:00:00Z",
		DueDateOnOrBefore: "2025-04-30T00:00:00Z",
		Completed:         &completed,
	}
	tasks, err := tested.ListTasks(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(tasks.Entries) != 1 {
		spew.Dump(tasks)
		t.Fatal("Unexpected result")
	}
	task := tasks.Entries[0]
	if task.CustomerUUID != "" || task.AssociatedObject != "contact" || task.AssociatedObjectUUID != "con_00000000-0000-0000-0000-000000000000" {
		spew.Dump(task)
		t.Fatal("Expected a contact task")
	}
}

func TestCreateTaskWithAssociatedObjectIdentifier(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				body := decodeNewTask(t, r)
				if body.CustomerUUID != "" {
					t.Errorf("Unexpected customer_uuid %v", body.CustomerUUID)
				}
				expected := AssociatedObjectIdentifier{
					AssociatedObject: "contact",
					Method:           "uuid",
					Value:            "con_00000000-0000-0000-0000-000000000000",
				}
				if body.AssociatedObjectIdentifier == nil || *body.AssociatedObjectIdentifier != expected {
					t.Errorf("Unexpected associated_object_identifier %v", body.AssociatedObjectIdentifier)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"task_uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": null,
					"associated_object": "contact",
					"associated_object_uuid": "con_00000000-0000-0000-0000-000000000000",
					"assignee": "keith+test1@chartmogul.com",
					"task_details": "This is some task details text.",
					"due_date": "2025-04-30T00:00:00Z",
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
		AssociatedObjectIdentifier: &AssociatedObjectIdentifier{
			AssociatedObject: AssociatedObjectContact,
			Method:           AssociatedObjectIdentifierMethodUUID,
			Value:            "con_00000000-0000-0000-0000-000000000000",
		},
		Assignee:    "keith+test1@chartmogul.com",
		TaskDetails: "This is some task details text.",
		DueDate:     "2025-04-30T00:00:00Z",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if task.AssociatedObjectUUID != "con_00000000-0000-0000-0000-000000000000" || task.CustomerUUID != "" {
		spew.Dump(task)
		t.Fatal("Unexpected result")
	}
}

func decodeNewTask(t *testing.T, r *http.Request) NewTask {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	var body NewTask
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	return body
}
