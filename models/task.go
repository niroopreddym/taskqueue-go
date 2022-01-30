package models

import (
	"time"

	"github.com/niroopreddym/taskqueue-go/enums"
)

//Task ...
type Task struct {
	ID           string
	IsCompleted  bool
	Status       enums.Status // untouched, completed, failed, timeout
	CreationTime time.Time    // when was the task created
	TaskData     string       // field containing data about the task
}
