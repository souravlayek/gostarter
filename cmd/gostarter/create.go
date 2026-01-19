package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// createProject handles the project creation functionality
func createProject(args []string, isWebServer bool, includeTests bool, includeAir bool, middlewares string, prefix string) {
	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	} else {
		// Get project name from user if not provided
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter project name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading project name:", err)
		}
		projectName = strings.TrimSpace(input)
	}

	// Validate project name
	if projectName == "" {
		fmt.Println("Project name cannot be empty. Exiting.")
		return
	}

	// Construct the module path
	modulePath := projectName
	if prefix != "" {
		modulePath = strings.TrimSuffix(prefix, "/") + "/" + projectName
	}

	fmt.Printf("Creating project: %s (module: %s)\n", projectName, modulePath)

	// Create project directory
	err := os.MkdirAll(projectName, 0755)
	if err != nil {
		log.Fatal("Error creating project directory:", err)
	}

	// Change to project directory
	err = os.Chdir(projectName)
	if err != nil {
		log.Fatal("Error changing to project directory:", err)
	}

	// Initialize go module
	initCmd := exec.Command("go", "mod", "init", modulePath)
	output, err := initCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error initializing go module: %s\n", output)
		return
	}

	// Determine template directory based on project type
	templateDir := "base"
	if isWebServer {
		templateDir = "web-server"
	}

	// Define project data for templates
	projectData := map[string]interface{}{
		"ProjectName": projectName,
		"ModuleName":  modulePath,
		"Description": fmt.Sprintf("A new Go project created with gostarter."),
		"Author":      "Your Name",
		"Port":        "8080",
		"Host":        "localhost",
	}

	// Copy the entire template directory structure
	copyTemplateDir(templateDir, projectData, includeTests, includeAir)

	// If it's a web server project, process middleware based on user selection
	if isWebServer {
		// Always create the middleware directory
		err := os.MkdirAll("internal/middleware", 0755)
		if err != nil {
			log.Printf("Warning: Error creating middleware directory: %v", err)
		}

		// Always copy the main.go middleware file which is always needed
		mainContent, err := templateFiles.ReadFile("templates/web-server/internal/middleware/main.go.tmpl")
		if err != nil {
			log.Printf("Warning: Error reading middleware main template: %v", err)
		} else {
			// Create the file with processed template
			tmpl, err := template.New("main.go.tmpl").Parse(string(mainContent))
			if err != nil {
				log.Printf("Warning: Error parsing middleware main template: %v", err)
			} else {
				file, err := os.Create("internal/middleware/main.go")
				if err != nil {
					log.Printf("Warning: Error creating middleware main.go: %v", err)
				} else {
					err = tmpl.Execute(file, projectData)
					if err != nil {
						log.Printf("Warning: Error executing middleware main template: %v", err)
						file.Close()
					} else {
						file.Close()
					}
				}
			}
		}

		// Process specific middleware if any were selected
		if middlewares != "" {
			processMiddleware(middlewares, projectData)
		}
	}


	// If Air is enabled and it's a web server project, copy the air.toml file
	if includeAir && isWebServer {
		// Process and copy the air.toml file
		content, err := templateFiles.ReadFile("templates/web-server/air.toml.tmpl")
		if err != nil {
			log.Printf("Warning: Error reading air.toml template: %v", err)
		} else {
			// Create the file with processed template
			tmpl, err := template.New("air.toml.tmpl").Parse(string(content))
			if err != nil {
				log.Printf("Warning: Error parsing air.toml template: %v", err)
			} else {
				file, err := os.Create("air.toml")
				if err != nil {
					log.Printf("Warning: Error creating air.toml: %v", err)
				} else {
					err = tmpl.Execute(file, projectData)
					if err != nil {
						log.Printf("Warning: Error executing air.toml template: %v", err)
						file.Close()
					} else {
						file.Close()
					}
				}
			}
		}

		// Also create a Makefile that supports Air
		airMakefileContent := `# Makefile for {{.ProjectName}}

.PHONY: build start install-air

# Build the project
build:
	go build -o bin/{{.ProjectName}} ./cmd/{{.ProjectName}}

# Install Air for live reloading (optional)
install-air:
	go install github.com/cosmtrek/air@latest

# Start the project with Air
start:
	air
`
		// Create the Makefile with processed template
		tmpl, err := template.New("Makefile").Parse(airMakefileContent)
		if err != nil {
			log.Printf("Warning: Error parsing Makefile template: %v", err)
		} else {
			file, err := os.Create("Makefile")
			if err != nil {
				log.Printf("Warning: Error creating Makefile: %v", err)
			} else {
				err = tmpl.Execute(file, projectData)
				if err != nil {
					log.Printf("Warning: Error executing Makefile template: %v", err)
					file.Close()
				} else {
					file.Close()
				}
			}
		}
	}

	fmt.Printf("Project '%s' has been created successfully!\n", projectName)

	// Display next steps
	fmt.Println("\nNext steps:")
	fmt.Printf("1. cd %s\n", projectName)
	fmt.Println("2. Run 'make build' to build the project")
	fmt.Println("3. Run 'make start' to start the project")

	if includeAir {
		fmt.Println("4. Install air using 'make install-air' or 'go install github.com/cosmtrek/air@latest'")
		fmt.Println("5. Run 'air' to start with hot reload")
	}

	fmt.Printf("\nYour project is ready! Enjoy coding!\n")

	// Print instructions for installing third-party packages
	fmt.Println("\nTo install required third-party packages, run:")
	fmt.Println("go mod tidy")
}


// processMiddleware handles the middleware files based on user selection
func processMiddleware(middlewares string, projectData map[string]interface{}) {
	// Split the comma-separated middleware list
	middlewareList := strings.Split(middlewares, ",")

	// Process each middleware
	for _, middleware := range middlewareList {
		middleware = strings.TrimSpace(middleware)

		// Process the middleware file based on the name
		var fileName string
		switch middleware {
		case "realIP":
			fileName = "realIP.go.tmpl"
		case "logger":
			fileName = "logger.go.tmpl"
		case "ratelimit":
			fileName = "ratelimit.go.tmpl"
		default:
			log.Printf("Warning: Unknown middleware '%s'", middleware)
			continue
		}

		// Process and copy the middleware file
		content, err := templateFiles.ReadFile("templates/web-server/internal/middleware/" + fileName)
		if err != nil {
			log.Printf("Warning: Error reading middleware template %s: %v", fileName, err)
			continue
		}

		// Create the file with processed template
		tmpl, err := template.New(fileName).Parse(string(content))
		if err != nil {
			log.Printf("Warning: Error parsing middleware template %s: %v", fileName, err)
			continue
		}

		// Create the output filename by removing .tmpl extension
		outputFileName := strings.TrimSuffix(fileName, ".tmpl")
		file, err := os.Create("internal/middleware/" + outputFileName)
		if err != nil {
			log.Printf("Warning: Error creating middleware file %s: %v", outputFileName, err)
			continue
		}

		err = tmpl.Execute(file, projectData)
		if err != nil {
			log.Printf("Warning: Error executing middleware template %s: %v", fileName, err)
			file.Close()
			continue
		}

		file.Close()
	}
}

// copyTemplateDir copies the entire template directory structure to the project
func copyTemplateDir(templateDir string, projectData map[string]interface{}, includeTests bool, includeAir bool) error {
	// Get all files in the template directory
	entries, err := templateFiles.ReadDir("templates/" + templateDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip the middleware directory during general copy since it's handled separately
			if entry.Name() == "middleware" {
				continue
			}

			// Recursively copy subdirectories
			// Replace "project-name" with actual project name in directory names
			dirName := entry.Name()
			if dirName == "project-name" {
				dirName = projectData["ProjectName"].(string)
			}

			err = copyTemplateDirRecursive("templates/"+templateDir+"/"+entry.Name(), dirName, projectData, includeTests, includeAir)
			if err != nil {
				log.Printf("Warning: Error copying directory %s: %v", dirName, err)
			}
		} else {
			// Skip test files if not including tests
			if !includeTests && strings.Contains(entry.Name(), "_test") {
				continue
			}

			// Process and copy the file
			content, err := templateFiles.ReadFile("templates/" + templateDir + "/" + entry.Name())
			if err != nil {
				log.Printf("Warning: Error reading file %s: %v", entry.Name(), err)
				continue
			}

			// Create the file with processed template
			tmpl, err := template.New(entry.Name()).Parse(string(content))
			if err != nil {
				log.Printf("Warning: Error parsing template %s: %v", entry.Name(), err)
				continue
			}

			// Remove .tmpl extension for the output filename
			outputFilename := strings.TrimSuffix(entry.Name(), ".tmpl")

			file, err := os.Create(outputFilename)
			if err != nil {
				log.Printf("Warning: Error creating file %s: %v", outputFilename, err)
				continue
			}

			err = tmpl.Execute(file, projectData)
			if err != nil {
				log.Printf("Warning: Error executing template %s: %v", entry.Name(), err)
				file.Close()
				continue
			}

			file.Close()
		}
	}

	return nil
}

// copyTemplateDirRecursive copies a subdirectory recursively
func copyTemplateDirRecursive(srcPath, destPath string, projectData map[string]interface{}, includeTests bool, includeAir bool) error {
	// Create destination directory
	err := os.MkdirAll(destPath, 0755)
	if err != nil {
		return err
	}

	// Change to destination directory for file operations
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(destPath)
	if err != nil {
		return err
	}
	defer os.Chdir(currentDir) // Restore original directory

	// Read source directory
	entries, err := templateFiles.ReadDir(srcPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcFilePath := srcPath + "/" + entry.Name()

		if entry.IsDir() {
			// Skip the middleware directory during general copy since it's handled separately
			if entry.Name() == "middleware" {
				continue
			}

			// Replace "project-name" with actual project name in directory names
			dirName := entry.Name()
			if dirName == "project-name" {
				dirName = projectData["ProjectName"].(string)
			}

			// Recursively copy subdirectory
			err = copyTemplateDirRecursive(srcFilePath, dirName, projectData, includeTests, includeAir)
			if err != nil {
				log.Printf("Warning: Error copying subdirectory %s: %v", dirName, err)
			}
		} else {
			// Skip test files if not including tests
			if !includeTests && strings.Contains(entry.Name(), "_test") {
				continue
			}

			// Process and copy the file
			content, err := templateFiles.ReadFile(srcFilePath)
			if err != nil {
				log.Printf("Warning: Error reading file %s: %v", srcFilePath, err)
				continue
			}

			// Create the file with processed template
			tmpl, err := template.New(entry.Name()).Parse(string(content))
			if err != nil {
				log.Printf("Warning: Error parsing template %s: %v", entry.Name(), err)
				continue
			}

			// Remove .tmpl extension for the output filename
			outputFilename := strings.TrimSuffix(entry.Name(), ".tmpl")

			file, err := os.Create(outputFilename)
			if err != nil {
				log.Printf("Warning: Error creating file %s: %v", outputFilename, err)
				continue
			}

			err = tmpl.Execute(file, projectData)
			if err != nil {
				log.Printf("Warning: Error executing template %s: %v", entry.Name(), err)
				file.Close()
				continue
			}

			file.Close()
		}
	}

	return nil
}
