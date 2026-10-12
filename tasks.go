package chartmogul

// Task is the task as represented in the API.
type Task struct {
	UUID string `json:"task_uuid"`
	// Basic info
	CustomerUUID         string `json:"customer_uuid,omitempty"`          // empty for tasks attached to a contact
	AssociatedObject     string `json:"associated_object,omitempty"`      // "customer" | "contact"
	AssociatedObjectUUID string `json:"associated_object_uuid,omitempty"` // cus_… or con_…
	Assignee             string `json:"assignee"`
	TaskDetails          string `json:"task_details"`
	DueDate              string `json:"due_date"`
	CompletedAt          string `json:"completed_at,omitempty"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

// UpdateTask allows for updating a task through the update endpoint.
type UpdateTask struct {
	Assignee    string `json:"assignee,omitempty"`
	TaskDetails string `json:"task_details,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// NewTask allows for creating a task through the new endpoint.
type NewTask struct {
	// Exactly one of CustomerUUID or AssociatedObjectIdentifier is required.
	CustomerUUID               string                      `json:"customer_uuid,omitempty"`
	AssociatedObjectIdentifier *AssociatedObjectIdentifier `json:"associated_object_identifier,omitempty"`

	// Obligatory
	Assignee    string `json:"assignee"`
	TaskDetails string `json:"task_details"`
	DueDate     string `json:"due_date"`

	// Optional
	CompletedAt string `json:"completed_at,omitempty"`
}

// ListTasksParams = parameters for listing tasks in API.
type ListTasksParams struct {
	CustomerUUID      string `json:"customer_uuid,omitempty"`
	ContactUUID       string `json:"contact_uuid,omitempty"`
	Assignee          string `json:"assignee,omitempty"`
	DueDateOnOrAfter  string `json:"due_date_on_or_after,omitempty"`  // ISO 8601
	DueDateOnOrBefore string `json:"due_date_on_or_before,omitempty"` // ISO 8601
	Completed         *bool  `json:"completed,omitempty"`
	Cursor
}

// Tasks is result of listing tasks in API.
type Tasks struct {
	Entries []*Task `json:"entries"`
	Pagination
}

const (
	singleTaskEndpoint = "tasks/:uuid"
	tasksEndpoint      = "tasks"
)

// CreateTask creates the task through the API.
func (api API) CreateTask(input *NewTask) (*Task, error) {
	result := &Task{}
	return result, api.create(tasksEndpoint, input, result)
}

// RetrieveTask returns one task from the API.
func (api API) RetrieveTask(taskUUID string) (*Task, error) {
	result := &Task{}
	return result, api.retrieve(singleTaskEndpoint, taskUUID, result)
}

// UpdateTask updates one task through the API.
// The API responds with 304 and an empty body when the input carries no updatable field,
// which is reported as an error.
func (api API) UpdateTask(input *UpdateTask, taskUUID string) (*Task, error) {
	output := &Task{}
	return output, api.update(singleTaskEndpoint, taskUUID, input, output)
}

// ListTasks lists all tasks.
func (api API) ListTasks(listTasksParams *ListTasksParams) (*Tasks, error) {
	result := &Tasks{}
	query := make([]interface{}, 0, 1)
	if listTasksParams != nil {
		query = append(query, *listTasksParams)
	}
	return result, api.list(tasksEndpoint, result, query...)
}

// DeleteTask deletes one task by UUID.
func (api API) DeleteTask(taskUUID string) error {
	return api.delete(singleTaskEndpoint, taskUUID)
}
