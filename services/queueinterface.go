package services

import (
	"context"

	"github.com/niroopreddym/taskqueue-go/models"
)

//IQueue ...
type IQueue interface {
	Adder(ctx context.Context, task *models.Task) error
	Cleaner(c context.Context, cleaningItem models.Task) error
	Executor(c context.Context) error
}
