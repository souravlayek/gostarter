# gostarter

A CLI tool to quickly setup fresh Go projects with a standard directory structure and files.

## Installation

### Using Go
```bash
go install github.com/souravlayek/gostarter@latest
```

### Using Package Managers

#### Homebrew (macOS/Linux)
```bash
brew install souravlayek/tap/gostarter
```

#### Chocolatey (Windows)
```powershell
choco install gostarter
```

#### APT (Debian/Ubuntu)
```bash
# Add the repository
curl -fsSL https://example.com/gpg-key | sudo gpg --dearmor -o /usr/share/keyrings/gostarter-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/gostarter-archive-keyring.gpg] https://example.com/apt-repo stable main" | sudo tee /etc/apt/sources.list.d/gostarter.list

# Update and install
sudo apt-get update
sudo apt-get install gostarter
```

### From Source
Clone and build from source:

```bash
git clone https://github.com/souravlayek/gostarter.git
cd gostarter
go build -o gostarter ./cmd/gostarter
```

## Usage

### Create a new project

```bash
gostarter create [project-name]
```

#### Options:
- `--web-server`: Create a web server project with HTTP handlers
- `--include-tests`: Include test files in the project (default: true)
- `--air`: Setup Air for auto-reloading (web server projects only)
- `--middlewares`: Comma-separated middleware list (realIP,logger,ratelimit) (web server projects only)
- `--prefix`: Module prefix (e.g., github.com/username/)

### Examples:

Create a basic Go project:
```bash
gostarter create myproject
```

Create a web server project:
```bash
gostarter create mywebserver --web-server
```

Create a web server project with Air auto-reloader:
```bash
gostarter create mywebserver --web-server --air
```

Create a web server project with specific middleware:
```bash
gostarter create mywebserver --web-server --middlewares realIP,logger
```

Create a project with a custom module prefix:
```bash
gostarter create myproject --prefix github.com/username/
```

Create a project without test files:
```bash
gostarter create myproject --include-tests=false
```

## Adding Packages

Add predefined packages to your project:

List all available packages:
```bash
gostarter add-pkg --list
```

Add a specific package to your project:
```bash
gostarter add-pkg logger
```

This will add the logger package to your project's `pkg/logger` directory.

## Adding Middleware

Add predefined middleware to your project:

List all available middleware:
```bash
gostarter add-middleware --list
```

Add a specific middleware to your project:
```bash
gostarter add-middleware logger
```

This will add the logger middleware to your project's `internal/middleware/logger.go` file.

## Features

- Quick project scaffolding with standard Go project structure
- Two project types: base Go projects and web server projects
- Pre-configured templates for common files (main.go, README.md, .gitignore, etc.)
- Optional test file generation
- Custom module prefix support
- Air auto-reload setup for development
- Configurable middleware selection (realIP, logger, ratelimit)
- Package management system to add reusable packages to projects
- Middleware management system to add reusable middleware to projects
- Command-line interface with intuitive flags
- Standard Go project conventions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.