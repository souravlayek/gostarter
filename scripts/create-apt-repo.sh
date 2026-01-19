#!/bin/bash

# Script to create APT repository for gostarter
# This script should be run on release to update the APT repository

set -e

# Variables
PACKAGE_NAME="gostarter"
VERSION=$(git describe --tags --abbrev=0 | sed 's/v//')
ARCHITECTURE="amd64"
DISTRO="focal"  # Ubuntu 20.04, but this can be generalized
REPO_DIR="apt-repo"
BIN_DIR="${REPO_DIR}/pool/main/g/${PACKAGE_NAME}"

# Create directory structure
mkdir -p "${BIN_DIR}"

# Download the binary for the target architecture
BINARY_URL="https://github.com/souravlayek/gostarter/releases/download/v${VERSION}/gostarter-linux-amd64"
BINARY_NAME="${PACKAGE_NAME}_${VERSION}_linux_amd64.tar.gz"
wget -O "${BIN_DIR}/${BINARY_NAME}" "${BINARY_URL}"

# Create control file for the package
CONTROL_FILE="${REPO_DIR}/control"
cat > "${CONTROL_FILE}" << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: devel
Priority: optional
Architecture: ${ARCHITECTURE}
Depends: 
Maintainer: Sourav Layek <sourav@example.com>
Description: A CLI tool to quickly setup fresh Go projects with a standard directory structure and files.
 A Go project starter CLI tool that helps you setup fresh Go projects from scratch.
 .
 Features:
  * Quick project scaffolding with standard Go project structure
  * Two project types: base Go projects and web server projects
  * Pre-configured templates for common files
  * Optional test file generation
  * Custom module prefix support
  * Air auto-reload setup for development
  * Configurable middleware selection
  * Package management system
  * Middleware management system
EOF

# Create the package file
DEBIAN_PACKAGE="${BIN_DIR}/${PACKAGE_NAME}_${VERSION}_${ARCHITECTURE}.deb"

# Create temporary directory for package building
TEMP_DIR=$(mktemp -d)
mkdir -p "${TEMP_DIR}/DEBIAN"
mkdir -p "${TEMP_DIR}/usr/bin"

# Copy binary to the package structure
# For now, we'll just create a placeholder - in reality, we'd need to create a proper deb package
cp "${BIN_DIR}/${BINARY_NAME}" "${TEMP_DIR}/usr/bin/"

# Create control file for the package
cat > "${TEMP_DIR}/DEBIAN/control" << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: devel
Priority: optional
Architecture: ${ARCHITECTURE}
Depends: 
Maintainer: Sourav Layek <sourav@example.com>
Description: A CLI tool to quickly setup fresh Go projects with a standard directory structure and files.
 A Go project starter CLI tool that helps you setup fresh Go projects from scratch.
EOF

# Build the package
dpkg-deb --build "${TEMP_DIR}" "${DEBIAN_PACKAGE}"

# Clean up
rm -rf "${TEMP_DIR}"

echo "APT repository structure created:"
echo "  - ${DEBIAN_PACKAGE}"
echo "Version: ${VERSION}"

# Instructions for APT repo hosting
echo ""
echo "To host the APT repository:"
echo "  1. Upload the repository to a web server"
echo "  2. Create a Release file in the dists directory"
echo "  3. Sign the Release file with GPG"
echo "  4. Users can add the repo with: sudo add-apt-repository 'deb [arch=${ARCHITECTURE}] https://your-domain.com/apt-repo focal main'"
echo "  5. Then run: sudo apt-get update && sudo apt-get install ${PACKAGE_NAME}"