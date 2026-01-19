package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// setupProject handles the project setup functionality
func setupProject() {
	fmt.Println("Setting up a new Go project...")

	// Get project name from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter project name: ")
	projectName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading project name:", err)
	}
	projectName = strings.TrimSpace(projectName)

	// Validate project name
	if projectName == "" {
		fmt.Println("Project name cannot be empty. Exiting.")
		return
	}

	// Create project directory
	err = os.MkdirAll(projectName, 0755)
	if err != nil {
		log.Fatal("Error creating project directory:", err)
	}

	// Change to project directory
	err = os.Chdir(projectName)
	if err != nil {
		log.Fatal("Error changing to project directory:", err)
	}

	// Initialize go module
	initCmd := exec.Command("go", "mod", "init", projectName)
	output, err := initCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error initializing go module: %s\n", output)
		return
	}

	// Create basic directory structure
	dirs := []string{"cmd", "internal", "pkg", "docs", "test", "config"}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("Warning: Error creating directory %s: %v", dir, err)
		}
	}

	// Define project data for templates
	projectData := map[string]interface{}{
		"ProjectName": projectName,
		"Description": "A new Go project created with gostarter.",
		"Author":      "Your Name",
		"Port":        "8080",
		"Host":        "localhost",
	}

	// Determine the executable's directory to find templates
	execPath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Could not determine executable path: %v", err)
		// Fallback to current directory
		execPath, _ = os.Getwd()
	}
	execDir := filepath.Dir(execPath)
	// Look for templates in different possible locations
	possibleTemplatePaths := []string{
		filepath.Join(execDir, "..", "..", "templates"), // From cmd/gostarter to root/templates
		filepath.Join(execDir, "templates"),             // Templates in same directory as executable
		filepath.Join(".", "templates"),                 // Templates in current directory
		filepath.Join("..", "templates"),                // Templates in parent directory
	}

	templatesDir := ""
	for _, path := range possibleTemplatePaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				templatesDir = absPath
				break
			}
		}
	}

	// If no templates directory found, use current directory
	if templatesDir == "" {
		templatesDir = "./templates"
		log.Printf("Warning: Could not find templates directory in standard locations, using: %s", templatesDir)
	}

	// Create main.go in cmd directory using template
	cmdDir := "cmd/" + projectName
	err = os.MkdirAll(cmdDir, 0755)
	if err != nil {
		log.Printf("Warning: Error creating cmd directory: %v", err)
	} else {
		// Parse and execute the main.go template
		tmpl, err := template.ParseFiles(filepath.Join(templatesDir, "main.go.tmpl"))
		if err != nil {
			log.Printf("Warning: Error parsing main.go template: %v", err)
			// Fallback to basic content
			mainContent := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello, %s!")
}`, projectName)
			err = os.WriteFile(cmdDir+"/main.go", []byte(mainContent), 0644)
			if err != nil {
				log.Printf("Warning: Error creating main.go: %v", err)
			}
		} else {
			file, err := os.Create(cmdDir + "/main.go")
			if err != nil {
				log.Printf("Warning: Error creating main.go: %v", err)
			} else {
				err = tmpl.ExecuteTemplate(file, "main.go.tmpl", projectData)
				if err != nil {
					log.Printf("Warning: Error executing main.go template: %v", err)
				}
				file.Close()
			}
		}

		// Create test file in the same directory
		testTmpl, err := template.ParseFiles(filepath.Join(templatesDir, "main_test.go.tmpl"))
		if err != nil {
			log.Printf("Warning: Error parsing main_test.go template: %v", err)
			// Fallback to basic test content
			testContent := fmt.Sprintf(`package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// TODO: Add your tests here
	expected := "Hello, %s!"
	actual := getGreeting() // assuming you have this function

	if expected != actual {
		t.Errorf("Expected %%s, got %%s", expected, actual)
	}
}

func getGreeting() string {
	return "Hello, %s!"
}`, projectName, projectName)
			err = os.WriteFile(cmdDir+"/main_test.go", []byte(testContent), 0644)
			if err != nil {
				log.Printf("Warning: Error creating main_test.go: %v", err)
			}
		} else {
			testFile, err := os.Create(cmdDir + "/main_test.go")
			if err != nil {
				log.Printf("Warning: Error creating main_test.go: %v", err)
			} else {
				err = testTmpl.ExecuteTemplate(testFile, "main_test.go.tmpl", projectData)
				if err != nil {
					log.Printf("Warning: Error executing main_test.go template: %v", err)
				}
				testFile.Close()
			}
		}
	}

	// Create README.md using template
	readmeTmpl, err := template.ParseFiles(filepath.Join(templatesDir, "README.md.tmpl"))
	if err != nil {
		log.Printf("Warning: Error parsing README.md template: %v", err)
		// Fallback to basic content
		readmeContent := fmt.Sprintf("# %s\n\nA new Go project created with gostarter.\n", projectName)
		err = os.WriteFile("README.md", []byte(readmeContent), 0644)
		if err != nil {
			log.Printf("Warning: Error creating README.md: %v", err)
		}
	} else {
		file, err := os.Create("README.md")
		if err != nil {
			log.Printf("Warning: Error creating README.md: %v", err)
		} else {
			err = readmeTmpl.ExecuteTemplate(file, "README.md.tmpl", projectData)
			if err != nil {
				log.Printf("Warning: Error executing README.md template: %v", err)
			}
			file.Close()
		}
	}

	// Create .gitignore using template
	gitignoreTmpl, err := template.ParseFiles(filepath.Join(templatesDir, ".gitignore.tmpl"))
	if err != nil {
		log.Printf("Warning: Error parsing .gitignore template: %v", err)
		// Fallback to basic content
		gitignoreContent := `# Binaries for programs and plugins
*.exe
*.dll
*.so
*.dylib

# Test binary, built with ` + "`go test -c`" + `
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove if you don't want them)
vendor/

# Go workspace file
go.work
`
		err = os.WriteFile(".gitignore", []byte(gitignoreContent), 0644)
		if err != nil {
			log.Printf("Warning: Error creating .gitignore: %v", err)
		}
	} else {
		file, err := os.Create(".gitignore")
		if err != nil {
			log.Printf("Warning: Error creating .gitignore: %v", err)
		} else {
			err = gitignoreTmpl.ExecuteTemplate(file, ".gitignore.tmpl", projectData)
			if err != nil {
				log.Printf("Warning: Error executing .gitignore template: %v", err)
			}
			file.Close()
		}
	}

	// Create config file using template
	configTmpl, err := template.ParseFiles(filepath.Join(templatesDir, "config.env.tmpl"))
	if err != nil {
		log.Printf("Warning: Error parsing config.env template: %v", err)
		// Fallback to basic content
		configContent := `# Configuration for ` + projectName + `

# Server settings
PORT=8080
HOST=localhost

# Logging level (debug, info, warn, error)
LOG_LEVEL=info

# Environment (development, staging, production)
ENVIRONMENT=development
`
		err = os.WriteFile("config/config.env", []byte(configContent), 0644)
		if err != nil {
			log.Printf("Warning: Error creating config.env: %v", err)
		}
	} else {
		file, err := os.Create("config/config.env")
		if err != nil {
			log.Printf("Warning: Error creating config.env: %v", err)
		} else {
			err = configTmpl.ExecuteTemplate(file, "config.env.tmpl", projectData)
			if err != nil {
				log.Printf("Warning: Error executing config.env template: %v", err)
			}
			file.Close()
		}
	}

	fmt.Printf("Project '%s' has been created successfully!\n", projectName)
}
