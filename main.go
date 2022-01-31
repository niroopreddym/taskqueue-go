package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/niroopreddym/taskqueue-go/enums"
	"github.com/niroopreddym/taskqueue-go/models"
	"github.com/niroopreddym/taskqueue-go/services"
)

func main() {
	ctx := context.Background()
	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	queue := services.NewQueue()

	go dataProcessor(ctx, queue)
	go dataCleaner(ctx, queue)
	go addDataToQueue(ctx, queue)

	select {
	case <-c:
		fmt.Println("cancel operation")
		cancel()
	case <-ctx.Done():
		time.Sleep(600 * time.Millisecond)
	}

	fmt.Println("done")
}

func addDataToQueue(ctx context.Context, queue *services.Queue) {
	for i := 0; ; i++ {
		task := models.Task{
			ID:           uuid.NewString(),
			IsCompleted:  false,
			Status:       enums.Untouched,
			CreationTime: time.Now(),
			TaskData:     "test data",
		}

		queue.Adder(ctx, task)

		select {
		case <-ctx.Done():
			fmt.Println("from data add")
			return
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func dataProcessor(ctx context.Context, queue *services.Queue) {
	// fmt.Println("inside data processor")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("from data processor")
			return
		default:
			err := queue.Executor(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func dataCleaner(ctx context.Context, queue *services.Queue) {
	// fmt.Println("inside data cleaner")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("from data cleaner")
			return
		case data := <-queue.ReadyToCleanItems:
			queue.Cleaner(ctx, data)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
