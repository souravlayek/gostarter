package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// addMiddleware adds a predefined middleware to the current project
func addMiddleware(middlewareName string) {
	// Check if we're in a Go project directory
	if !isGoProject() {
		log.Fatal("Error: Not in a Go project directory (missing go.mod)")
	}

	// Create internal/middleware directory if it doesn't exist
	err := os.MkdirAll("internal/middleware", 0755)
	if err != nil {
		log.Fatalf("Error creating internal/middleware directory: %v", err)
	}

	// Check if the requested middleware exists in our templates directory
	// We need to look for the specific middleware file
	middlewareFileName := middlewareName + ".go.tmpl"
	middlewarePath := "templates/web-server/internal/middleware/" + middlewareFileName
	
	// Read the file content from embedded FS
	content, err := templateFiles.ReadFile(middlewarePath)
	if err != nil {
		log.Fatalf("Error: Middleware '%s' does not exist in the available middleware", middlewareName)
	}

	// Create the file with processed template
	tmpl, err := template.New(middlewareFileName).Parse(string(content))
	if err != nil {
		log.Printf("Warning: Error parsing middleware template %s: %v", middlewareFileName, err)
		return
	}

	// Create the output filename by removing .tmpl extension
	outputFileName := strings.TrimSuffix(middlewareFileName, ".tmpl")
	file, err := os.Create("internal/middleware/" + outputFileName)
	if err != nil {
		log.Printf("Warning: Error creating middleware file %s: %v", outputFileName, err)
		return
	}

	// For middleware, we need to pass project data for template processing
	// For now, we'll pass an empty map, but in a real scenario you might want to pass actual project data
	projectData := make(map[string]interface{})
	err = tmpl.Execute(file, projectData)
	if err != nil {
		log.Printf("Warning: Error executing middleware template %s: %v", middlewareFileName, err)
		file.Close()
		return
	}

	file.Close()

	// Also copy the main.go middleware file which is always needed
	mainContent, err := templateFiles.ReadFile("templates/web-server/internal/middleware/main.go.tmpl")
	if err != nil {
		log.Printf("Warning: Error reading middleware main template: %v", err)
	} else {
		// Check if main.go already exists to avoid overwriting
		if _, err := os.Stat("internal/middleware/main.go"); os.IsNotExist(err) {
			// Create the file with processed template
			tmpl, err := template.New("main.go.tmpl").Parse(string(mainContent))
			if err != nil {
				log.Printf("Warning: Error parsing middleware main template: %v", err)
			} else {
				mainFile, err := os.Create("internal/middleware/main.go")
				if err != nil {
					log.Printf("Warning: Error creating middleware main.go: %v", err)
				} else {
					err = tmpl.Execute(mainFile, projectData)
					if err != nil {
						log.Printf("Warning: Error executing middleware main template: %v", err)
						mainFile.Close()
					} else {
						mainFile.Close()
					}
				}
			}
		}
	}

	fmt.Printf("Added middleware/%s\n", outputFileName)
	fmt.Printf("Middleware '%s' has been successfully added to your project!\n", middlewareName)
}

// listMiddlewares lists all available middleware
func listMiddlewares() {
	fmt.Println("Available middleware:")
	
	// Read all files in the middleware templates directory
	// We need to get the directory listing from the embedded FS
	entries, err := templateFiles.ReadDir("templates/web-server/internal/middleware")
	if err != nil {
		log.Printf("Error reading middleware directory: %v", err)
		return
	}
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.tmpl") {
			// Remove the .tmpl extension to get the middleware name
			middlewareName := strings.TrimSuffix(entry.Name(), ".go.tmpl")
			// Skip main.go.tmpl as it's not a specific middleware
			if middlewareName != "main" {
				fmt.Printf("- %s\n", middlewareName)
			}
		}
	}
	
	// Check if we found any middleware
	found := false
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.tmpl") && strings.TrimSuffix(entry.Name(), ".tmpl") != "main" {
			found = true
			break
		}
	}
	
	if !found {
		fmt.Println("No middleware available.")
	}
}