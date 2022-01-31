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
	Items             []models.Task
	ExecutorItems     chan models.Task
	ReadyToCleanItems chan models.Task
}

//NewQueue is ctor
func NewQueue() *Queue {
	// create items and isEmpty channels
	items := []models.Task{}
	readToCleanItems := make(chan models.Task, 10)
	ExecutorItems := make(chan models.Task, 10)
	// return queue
	return &Queue{items, ExecutorItems, readToCleanItems}
}

//Adder adds the items to Queue
func (q *Queue) Adder(ctx context.Context, i models.Task) {
	fmt.Println("added task data:", i)
	// append item to items
	q.Items = append(q.Items, i)
	q.ExecutorItems <- i
}

//Cleaner removes the item to Queue
func (q *Queue) Cleaner(c context.Context, cleaningItem models.Task) error {
	items := q.Items
	for index, val := range items {
		if cleaningItem.Status == enums.Completed && val.ID == cleaningItem.ID {
			removedItem := cleaningItem
			fmt.Println("removed item:", removedItem)
			copy(items[index:], items[index+1:])
			items = items[:len(items)-1]
			q.Items = items
		}

	}

	select {
	case <-c.Done():
		return c.Err()
	default:
		time.Sleep(1 * 1000)
	}

	return nil
}

//Executor removes the item to Queue
func (q *Queue) Executor(c context.Context) error {
	select {
	case item := <-q.ExecutorItems:
		//push the data to chan and mke the status is true
		item.IsCompleted = true
		item.Status = enums.Completed
		time.Sleep(500 * time.Millisecond)
		fmt.Println("altering data:", item)
		data := item
		q.ReadyToCleanItems <- data

	// mark queue as empty if last item is dequeued or update items
	case <-c.Done():
		return c.Err()
	}

	return nil
}
