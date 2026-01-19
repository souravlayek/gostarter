# Release Scripts

This directory contains scripts to automate the release process for gostarter across different package managers.

## Scripts

### `release.sh`
Main release script that:
- Builds binaries for different platforms (Linux, macOS, Windows)
- Creates archives with checksums
- Prepares package manager distributions (Homebrew, Chocolatey)

### `update-homebrew.sh`
Creates a Homebrew formula file for the current version.

### `create-chocolatey.sh`
Creates a Chocolatey package for Windows users.

### `create-apt-repo.sh`
Sets up an APT repository structure for Debian/Ubuntu users.

## Usage

### For Maintainers
1. Tag a new release: `git tag v1.2.3 && git push origin v1.2.3`
2. Run the release script: `./scripts/release.sh`
3. Upload the generated binaries to the GitHub release
4. Update package manager repositories with the new version

### For Users
The release process automatically publishes to:
- GitHub Releases (manual upload required)
- Homebrew (tap repository)
- Chocolatey (community repository)
- APT repository (for Linux users)

## Package Manager Installation

### Homebrew (macOS/Linux)
```bash
brew install souravlayek/tap/gostarter
```

### Chocolatey (Windows)
```powershell
choco install gostarter
```

### Manual Installation
Download the appropriate binary from the [GitHub releases page](https://github.com/souravlayek/gostarter/releases).