package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// godotenv
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(os.Getenv("name"))
	fmt.Println(os.Getenv("age"))
}
