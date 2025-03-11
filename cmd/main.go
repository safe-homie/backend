package main

import (
	"log"

	"github.com/safe-homie/backend/cmd/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
