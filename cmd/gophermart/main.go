package main

import (
	"log"

	"github.com/developerc/project_gophermart_2/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
