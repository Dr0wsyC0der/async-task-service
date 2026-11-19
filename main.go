package main

import (
	"fmt"

	"github.com/Dr0wsyC0der/async-task-service/internal/domain"
)

func main() {
	t := domain.Task{
		Code:     "print('Hello, world!')",
		Language: "python",
		Status:   domain.StatusInProgress,
		ID:       "some-id",
		Result:   "",
	}

	fmt.Printf("%+v\n", t)
}
