package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)


// addPackage adds a predefined package to the current project
func addPackage(pkgName string) {
	// Check if we're in a Go project directory
	if !isGoProject() {
		log.Fatal("Error: Not in a Go project directory (missing go.mod)")
	}

	// Create pkg directory if it doesn't exist
	err := os.MkdirAll("pkg", 0755)
	if err != nil {
		log.Fatalf("Error creating pkg directory: %v", err)
	}

	// Check if the requested package exists in our packages directory
	pkgPath := "packages/" + pkgName
	entries, err := packageFiles.ReadDir(pkgPath)
	if err != nil {
		log.Fatalf("Error: Package '%s' does not exist in the available packages", pkgName)
	}

	// Create the destination directory for the package
	destDir := filepath.Join("pkg", pkgName)
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		log.Fatalf("Error creating package directory: %v", err)
	}

	// Copy all files from the embedded package to the destination
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Skip directories for now, just handle files
		}

		// Read the file content from embedded FS
		content, err := packageFiles.ReadFile(pkgPath + "/" + entry.Name())
		if err != nil {
			log.Printf("Warning: Error reading file %s: %v", entry.Name(), err)
			continue
		}

		// Process the file as a template if it ends with .tmpl
		var processedContent []byte
		var outputFileName string

		if strings.HasSuffix(entry.Name(), ".tmpl") {
			// Create a template from the content
			tmplName := entry.Name()
			tmpl, err := template.New(tmplName).Parse(string(content))
			if err != nil {
				log.Printf("Warning: Error parsing template %s: %v", entry.Name(), err)
				continue
			}

			// Create project data for template processing
			// For now, we'll use a simple map; in a real scenario you might want to populate this with actual project data
			projectData := make(map[string]interface{})

			// Get project name from go.mod file
			goModContent, err := os.ReadFile("go.mod")
			if err == nil {
				// Extract module name from go.mod
				lines := strings.Split(string(goModContent), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if strings.HasPrefix(line, "module ") {
						moduleName := strings.TrimPrefix(line, "module ")
						moduleName = strings.TrimSpace(moduleName)
						projectData["ModuleName"] = moduleName

						// Extract project name from module path
						parts := strings.Split(moduleName, "/")
						projectName := parts[len(parts)-1]
						projectData["ProjectName"] = projectName
						break
					}
				}
			}

			// Execute the template with project data
			var buf strings.Builder
			err = tmpl.Execute(&buf, projectData)
			if err != nil {
				log.Printf("Warning: Error executing template %s: %v", entry.Name(), err)
				continue
			}

			processedContent = []byte(buf.String())

			// Remove .tmpl extension for output filename
			outputFileName = strings.TrimSuffix(entry.Name(), ".tmpl")
		} else {
			// Not a template file, use content as is
			processedContent = content
			outputFileName = entry.Name()
		}

		// Write the file to the destination
		destFile := filepath.Join(destDir, outputFileName)
		err = os.WriteFile(destFile, processedContent, 0644)
		if err != nil {
			log.Printf("Warning: Error writing file %s: %v", destFile, err)
			continue
		}

		fmt.Printf("Added %s/%s\n", pkgName, outputFileName)
	}

	fmt.Printf("Package '%s' has been successfully added to your project!\n", pkgName)

	// Suggest running go mod tidy to install dependencies
	fmt.Println("Consider running 'go mod tidy' to install any new dependencies.")
}

// listPackages lists all available packages
func listPackages() {
	fmt.Println("Available packages:")

	// Read all directories in the packages directory
	entries, err := packageFiles.ReadDir("packages")
	if err != nil {
		log.Printf("Error reading packages directory: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("- %s\n", entry.Name())
		}
	}

	if len(entries) == 0 {
		fmt.Println("No packages available.")
	}
}

// isGoProject checks if the current directory is a Go project by looking for go.mod
func isGoProject() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}