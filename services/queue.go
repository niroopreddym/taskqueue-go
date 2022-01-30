package services

import (
	"context"
	"fmt"
	"time"

	"github.com/niroopreddym/taskqueue-go/enums"
	"github.com/niroopreddym/taskqueue-go/models"
)

//Queue ...
type Queue struct {
	Items             chan []*models.Task
	ReadyToCleanItems chan []*models.Task
	IsEmpty           chan bool
}

//NewQueue is ctor
func NewQueue() *Queue {
	// create items and isEmpty channels
	items := make(chan []*models.Task, 1)
	readToCleanItems := make(chan []*models.Task, 1)
	isEmpty := make(chan bool, 1)
	// mark queue as empty
	isEmpty <- true
	// return queue
	return &Queue{items, readToCleanItems, isEmpty}
}

//Adder adds the items to Queue
func (q *Queue) Adder(ctx context.Context, i *models.Task) {
	// create items if queue is empty or get items
	var items []*models.Task
	select {
	case items = <-q.Items:
	case <-q.IsEmpty:
	}

	fmt.Println("added task data:", *i)
	// append item to items
	items = append(items, i)
	// update items
	q.Items <- items
}

//Cleaner removes the item to Queue
func (q *Queue) Cleaner(c context.Context) error {
	select {
	case <-c.Done():
		return c.Err()
	case items := <-q.ReadyToCleanItems:
		for index, val := range items {
			if val.Status == enums.Completed {
				removedItem := val
				fmt.Println("removed item:", removedItem)
				copy(items[index:], items[index+1:])
				items = items[:len(items)-1]
				q.Items <- items
			}
			time.Sleep(1 * 1000)
		}
	}

	if len(q.Items) == 0 {
		q.IsEmpty <- true
	}
	return nil
}

//Executor removes the item to Queue
func (q *Queue) Executor(c context.Context) error {
	var readyToCleanItems []*models.Task

	select {
	case items := <-q.Items:
		for _, item := range items {
			//push the data to chan and mke the status is true
			item.IsCompleted = true
			item.Status = enums.Completed
			time.Sleep(500 * time.Millisecond)
			q.ReadyToCleanItems <- append(readyToCleanItems, item)
			fmt.Println("altering data:", item)
		}

	// mark queue as empty if last item is dequeued or update items
	case <-c.Done():
		return c.Err()
	}

	// q.ReadyToCleanItems <- readyToCleanItems
	return nil
}

func dataTransform(readyToCleanItems []*models.Task) []models.Task {
	lastValues := []models.Task{}
	for _, val := range readyToCleanItems {
		data := *val
		lastValues = append(lastValues, data)
	}

	return lastValues
}

func dataTransform2(readyToCleanItems []models.Task) []*models.Task {
	lastValues := []*models.Task{}
	for _, val := range readyToCleanItems {
		lastValues = append(lastValues, &val)
	}

	return lastValues
}
