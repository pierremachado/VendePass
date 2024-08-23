package main

import (
	"fmt"
	"storkie/internal/models"

	"github.com/google/uuid"
)

func main() {

	fly := models.Fly{
		Id: uuid.New(),
	}
	fmt.Println(fly.Id)
}
