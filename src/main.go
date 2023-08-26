package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"rotavator/migrate"
	"rotavator/web"
)

func main() {
	var app web.App
	f := flag.NewFlagSet("main", flag.ContinueOnError)
	action := f.String("action", "web", "action to perform (web, update, rollback")
	err := f.Parse(os.Args)
	if err != nil {
		log.Printf("FATAL ARGS: %v", err)
		os.Exit(1)
	}

	switch *action {
	case "update":
		migrate.Update()
		return
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Printf("FATAL: %v", err)
		os.Exit(1)
	}
}
