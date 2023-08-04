package main

import (
	"log"
	"os"

	"rotorvator/rotavator"
)

func main() {
	var app rotavator.App

	if err := app.Run(); err != nil {
		log.Printf("FATAL: %v", err)
		os.Exit(1)
	}
}
