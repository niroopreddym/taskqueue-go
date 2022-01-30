package main

import (
	"fmt"

	"github.com/niroopreddym/taskqueue-go/services"
)

func main() {
	queue := services.NewQueue()
	fmt.Println(len(queue.Items))
}
