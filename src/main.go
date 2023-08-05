package main

import (
	"log"
	"os"

	"rotavator/web"
)

func main() {
	var app web.App

	if err := app.Run(); err != nil {
		log.Printf("FATAL: %v", err)
		os.Exit(1)
	}
}
