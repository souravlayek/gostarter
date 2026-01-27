package main

import (
	"log"
)

// setupLogging configures the logging system
func setupLogging() {
	// Set log format with timestamp and source file info
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// Optionally, you can set a custom log file instead of stderr
	// Uncomment the next lines if you want to log to a file
	/*
	file, err := os.OpenFile("gostarter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Println("Failed to open log file, using default stderr")
	}
	*/
}

// logInfo logs an informational message
func logInfo(message string) {
	log.Printf("[INFO] %s", message)
}

// logError logs an error message
func logError(message string) {
	log.Printf("[ERROR] %s", message)
}

// logUserInput logs user input specifically
func logUserInput(input string) {
	log.Printf("[USER_INPUT] %s", input)
}