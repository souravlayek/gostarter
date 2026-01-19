#!/bin/bash

# Script to create Chocolatey package
# This script should be run on release to create the Chocolatey package

set -e

# Variables
PACKAGE_NAME="gostarter"
VERSION=$(git describe --tags --abbrev=0 | sed 's/v//')
TOOLS_DIR="chocolatey/tools"
NUSPEC_FILE="${PACKAGE_NAME}.nuspec"

# Create tools directory
mkdir -p "${TOOLS_DIR}"

# Create chocolateyinstall.ps1
cat > "${TOOLS_DIR}/chocolateyinstall.ps1" << EOF
\$packageName = '${PACKAGE_NAME}'
\$toolsDir = "$(Split-Path -parent \$MyInvocation.MyCommand.Definition)"
\$url64 = "https://github.com/souravlayek/gostarter/releases/download/v${VERSION}/gostarter-windows-amd64.exe"

\$packageArgs = @{
  packageName   = \$packageName
  unzipLocation = \$toolsDir
  fileType      = 'exe'
  url64bit      = \$url64

  softwareName  = 'gostarter*'

  checksum64    = ''
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs
EOF

# Create the nuspec file from template
sed "s/{{VERSION}}/${VERSION}/g" chocolatey-template.nuspec > "${NUSPEC_FILE}"

echo "Chocolatey package files created:"
echo "  - ${NUSPEC_FILE}"
echo "  - ${TOOLS_DIR}/chocolateyinstall.ps1"
echo "Version: ${VERSION}"

# Instructions for packaging
echo ""
echo "To create the Chocolatey package:"
echo "  1. Install Chocolatey (choco install chocolatey)"
echo "  2. Run: choco pack ${NUSPEC_FILE}"
echo "  3. Test locally: choco install ${PACKAGE_NAME} --source . --force"
echo "  4. Push to Chocolatey: choco push ${PACKAGE_NAME}.${VERSION}.nupkg --source https://push.chocolatey.org/"