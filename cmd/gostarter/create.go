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
func createProject(args []string, isWebServer bool, includeTests bool, includeAir bool, middlewares string, database string, prefix string) {
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
		"Description": "A new Go project created with gostarter.",
		"Author":      "Your Name",
		"Port":        "8080",
		"Host":        "localhost",
		"Air":         true,
		"Database":    database,
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

	// Handle database option if specified
	if database != "" {
		setupDatabase(database, projectData)
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
	}

	fmt.Printf("Project '%s' has been created successfully!\n", projectName)

	// Display next steps
	fmt.Println("\nNext steps:")
	fmt.Printf("1. cd %s\n", projectName)
	fmt.Println("2. run `go mod tidy`")
	fmt.Println("\nRun 'make build' to build the project")
	fmt.Println("Run 'make start' to start the project")

	if includeAir {
		fmt.Println("\nInstall air using go install github.com/cosmtrek/air@latest")
	}
	if database == "postgresql" {
		fmt.Println("\nInstall sqlc using go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest")
		fmt.Println("We use sqlc for ORM you can checkout more details at https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html")
	}

	fmt.Printf("\nYour project is ready! Enjoy coding!\n")
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
			// Skip the middleware and database directories during general copy since they're handled separately
			if entry.Name() == "middleware" || entry.Name() == "database" {
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

// setupDatabase creates the database structure based on the selected database type
func setupDatabase(dbType string, projectData map[string]interface{}) {
	// Create the database directory inside internal
	err := os.MkdirAll("internal/database", 0755)
	if err != nil {
		log.Printf("Warning: Error creating internal/database directory: %v", err)
		return
	}

	// Create the main database file
	var dbContent []byte

	switch dbType {
	case "mongodb":
		dbContent, err = templateFiles.ReadFile("templates/web-server/internal/database/mongodb.go.tmpl")
		if err != nil {
			log.Printf("Warning: Error reading MongoDB template: %v", err)
			// Use a default MongoDB implementation
			dbContent = []byte(`package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB connects to MongoDB
func ConnectDB() {
	// Replace with your MongoDB connection string
	connectionString := "mongodb://localhost:27017"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	// Set global DB variable
	DB = client.Database("{{.ProjectName}}")

	// Send ping to confirm successful connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully!")
}
`)
		}
	case "postgresql":
		dbContent, err = templateFiles.ReadFile("templates/web-server/internal/database/postgresql.go.tmpl")
		if err != nil {
			log.Printf("Warning: Error reading PostgreSQL template: %v", err)
			// Use a default PostgreSQL implementation
			dbContent = []byte(`package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDB connects to PostgreSQL
func ConnectDB() {
	// Replace with your PostgreSQL connection string
	connectionString := "host=localhost port=5432 user=username password=password dbname={{.ProjectName}} sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping PostgreSQL:", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
}
`)
		}
	default:
		log.Printf("Warning: Unsupported database type '%s'. Supported types: mongodb, postgresql", dbType)
		return
	}

	// Process the template
	tmpl, err := template.New("database.go.tmpl").Parse(string(dbContent))
	if err != nil {
		log.Printf("Warning: Error parsing database template: %v", err)
		return
	}

	// Create the database file
	file, err := os.Create("internal/database/main.go")
	if err != nil {
		log.Printf("Warning: Error creating database file: %v", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, projectData)
	if err != nil {
		log.Printf("Warning: Error executing database template: %v", err)
		return
	}

	// For PostgreSQL, create additional directories and files
	if dbType == "postgresql" {
		// Create migrations directory inside database
		err = os.MkdirAll("internal/database/migrations", 0755)
		if err != nil {
			log.Printf("Warning: Error creating migrations directory: %v", err)
		}

		// Create query directory inside database
		err = os.MkdirAll("internal/database/query", 0755)
		if err != nil {
			log.Printf("Warning: Error creating internal/database/query directory: %v", err)
		}

		// Create sqlc.yaml file for PostgreSQL projects
		sqlcContent, err := templateFiles.ReadFile("templates/web-server/sqlc.yaml.tmpl")
		if err != nil {
			log.Printf("Warning: Error reading sqlc.yaml template: %v", err)
		} else {
			// Create the file with processed template
			tmpl, err := template.New("sqlc.yaml.tmpl").Parse(string(sqlcContent))
			if err != nil {
				log.Printf("Warning: Error parsing sqlc.yaml template: %v", err)
			} else {
				file, err := os.Create("sqlc.yaml")
				if err != nil {
					log.Printf("Warning: Error creating sqlc.yaml: %v", err)
				} else {
					err = tmpl.Execute(file, projectData)
					if err != nil {
						log.Printf("Warning: Error executing sqlc.yaml template: %v", err)
						file.Close()
					} else {
						file.Close()
					}
				}
			}
		}

	}
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
			// Skip the middleware and database directories during general copy since they're handled separately
			if entry.Name() == "middleware" || entry.Name() == "database" {
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
