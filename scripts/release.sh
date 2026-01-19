#!/bin/bash

# Comprehensive release script for gostarter
# This script automates the release process for different package managers

set -e

echo "Gostarter Release Script"
echo "======================="

# Get version from git tag
VERSION=$(git describe --tags --abbrev=0 | sed 's/v//')
if [ $? -ne 0 ]; then
    echo "Error: No tags found. Please create a version tag (e.g., v1.0.0)."
    exit 1
fi

echo "Detected version: ${VERSION}"
echo ""

# Confirm release
read -p "Proceed with releasing version ${VERSION}? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Release cancelled."
    exit 1
fi

echo ""
echo "Step 1: Building binaries for different platforms..."
echo "----------------------------------------------------"

# Create dist directory
mkdir -p dist

# Build for different platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
    GOOS=$(echo "$platform" | cut -d'/' -f1)
    GOARCH=$(echo "$platform" | cut -d'/' -f2)
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    export GOOS GOARCH
    
    if [ "$GOOS" = "windows" ]; then
        env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -ldflags="-s -w" -o "dist/gostarter-${GOOS}-${GOARCH}.exe" ./cmd/gostarter
    else
        env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -ldflags="-s -w" -o "dist/gostarter-${GOOS}-${GOARCH}" ./cmd/gostarter
    fi
    
    unset GOOS GOARCH
done

echo ""
echo "Binaries built successfully:"
ls -la dist/

echo ""
echo "Step 2: Creating archives..."
echo "-----------------------------"

# Create archives for each platform
for binary in dist/*; do
    if [ -f "$binary" ]; then
        filename=$(basename "$binary")
        platform=${filename#*-}  # Remove 'gostarter-' prefix
        arch_suffix="${platform#*-}"  # Get architecture part
        os_arch="${platform%-*}"  # Get os-arch part
        
        # Create archive
        archive_name="gostarter_${VERSION}_${os_arch}.tar.gz"
        tar -czf "dist/${archive_name}" -C dist "$filename"
        
        # Calculate checksum
        checksum=$(shasum -a 256 "dist/${archive_name}" | cut -d' ' -f1)
        echo "${checksum}  ${archive_name}" >> "dist/checksums.txt"
        
        echo "Created: ${archive_name} (checksum: ${checksum})"
    fi
done

echo ""
echo "Step 3: Preparing package manager distributions..."
echo "---------------------------------------------------"

# Prepare Homebrew formula
echo "Preparing Homebrew formula..."
ARCHIVE_URL="https://github.com/souravlayek/gostarter/releases/download/v${VERSION}/gostarter_${VERSION}_linux_amd64.tar.gz"
SHA256=$(curl -sL "${ARCHIVE_URL}" | shasum -a 256 | cut -d' ' -f1)

cat > "dist/gostarter.rb" << EOF
class Gostarter < Formula
  desc "A CLI tool to quickly setup fresh Go projects with a standard directory structure and files."
  homepage "https://github.com/souravlayek/gostarter"
  url "https://github.com/souravlayek/gostarter/releases/download/v#{version}/gostarter_#{version}_linux_amd64.tar.gz"
  sha256 "${SHA256}"

  def install
    bin.install "gostarter-#{OS.mac? ? 'darwin' : 'linux'}-#{Hardware::CPU.intel? ? 'amd64' : 'arm64'}" => "gostarter"
  end

  test do
    system "#{bin}/gostarter", "--help"
  end
end
EOF

echo "Homebrew formula created: dist/gostarter.rb"

# Prepare Chocolatey package
echo "Preparing Chocolatey package..."
mkdir -p dist/chocolatey/tools

# Create chocolateyinstall.ps1
cat > "dist/chocolatey/tools/chocolateyinstall.ps1" << EOF
\$ErrorActionPreference = 'Stop'

\$packageName = 'gostarter'
\$softwareName = 'gostarter*'
\$toolsDir   = "$(Split-Path -parent \$MyInvocation.MyCommand.Definition)"
\$url64      = 'https://github.com/souravlayek/gostarter/releases/download/v${VERSION}/gostarter-windows-amd64.exe'

\$packageArgs = @{
  packageName   = \$packageName
  unzipLocation = \$toolsDir
  fileType      = 'exe'
  url64bit      = \$url64

  softwareName  = \$softwareName
  checksum64    = ''
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs
EOF

# Create nuspec file
cat > "dist/gostarter.nuspec" << EOF
<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
  <metadata>
    <id>gostarter</id>
    <version>${VERSION}</version>
    <title>gostarter</title>
    <authors>Sourav Layek</authors>
    <description>A CLI tool to quickly setup fresh Go projects with a standard directory structure and files.</description>
    <summary>A Go project starter CLI tool</summary>
    <projectUrl>https://github.com/souravlayek/gostarter</projectUrl>
    <licenseUrl>https://github.com/souravlayek/gostarter/blob/main/LICENSE</licenseUrl>
    <tags>go golang cli development tools</tags>
    <owners>souravlayek</owners>
    <requireLicenseAcceptance>false</requireLicenseAcceptance>
  </metadata>
  <files>
    <file src="tools\**" target="tools" />
  </files>
</package>
EOF

echo "Chocolatey package prepared in dist/chocolatey/"

echo ""
echo "Step 4: Release preparation complete!"
echo "--------------------------------------"
echo "Files created in dist/:"
ls -la dist/

echo ""
echo "Next steps:"
echo "1. Create a GitHub release with the binaries in dist/"
echo "2. Update the Homebrew tap with the formula in dist/gostarter.rb"
echo "3. Build and publish the Chocolatey package: choco pack dist/gostarter.nuspec"
echo "4. Update documentation with the new version"