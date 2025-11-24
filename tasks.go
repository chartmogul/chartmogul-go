package chartmogul

// Task is the task as represented in the API.
type Task struct {
	UUID string `json:"task_uuid"`
	// Basic info
	CustomerUUID string `json:"customer_uuid"`
	Assignee     string `json:"assignee"`
	TaskDetails  string `json:"task_details"`
	DueDate      string `json:"due_date"`
	CompletedAt  string `json:"completed_at,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
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
	// Obligatory
	CustomerUUID string `json:"customer_uuid"`
	Assignee     string `json:"assignee"`
	TaskDetails  string `json:"task_details"`
	DueDate      string `json:"due_date"`

	// Optional
	CompletedAt string `json:"completed_at,omitempty"`
}

// ListTasksParams = parameters for listing tasks in API.
type ListTasksParams struct {
	CustomerUUID string `json:"customer_uuid,omitempty"`
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
//
// See https://dev.chartmogul.com/reference/create-a-task
func (api API) CreateTask(input *NewTask) (*Task, error) {
	result := &Task{}
	return result, api.create(tasksEndpoint, input, result)
}

// RetrieveTask returns one task from the API.
//
// See https://dev.chartmogul.com/reference/retrieve-a-task
func (api API) RetrieveTask(taskUUID string) (*Task, error) {
	result := &Task{}
	return result, api.retrieve(singleTaskEndpoint, taskUUID, result)
}

// UpdateTask updates one task through the API.
//
// See https://dev.chartmogul.com/reference/update-a-task
func (api API) UpdateTask(input *UpdateTask, taskUUID string) (*Task, error) {
	output := &Task{}
	return output, api.update(singleTaskEndpoint, taskUUID, input, output)
}

// ListTasks lists all tasks.
//
// See https://dev.chartmogul.com/reference/list-tasks
func (api API) ListTasks(listTasksParams *ListTasksParams) (*Tasks, error) {
	result := &Tasks{}
	query := make([]interface{}, 0, 1)
	if listTasksParams != nil {
		query = append(query, *listTasksParams)
	}
	return result, api.list(tasksEndpoint, result, query...)
}

// DeleteTask deletes one task by UUID.
//
// See https://dev.chartmogul.com/reference/delete-a-task
func (api API) DeleteTask(taskUUID string) error {
	return api.delete(singleTaskEndpoint, taskUUID)
}
