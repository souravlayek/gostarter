#!/bin/bash

# Script to update Homebrew formula
# This script should be run on release to update the Homebrew tap

set -e

# Variables
FORMULA_NAME="gostarter"
FORMULA_FILE="${FORMULA_NAME}.rb"
REPO_NAME="souravlayek/homebrew-tap"  # Replace with actual tap name
VERSION=$(git describe --tags --abbrev=0 | sed 's/v//')

# Calculate SHA256 for the main archive
ARCHIVE_URL="https://github.com/souravlayek/gostarter/archive/v${VERSION}.tar.gz"
SHA256=$(curl -sL "${ARCHIVE_URL}" | shasum -a 256 | cut -d' ' -f1)

# Create the formula file
cat > "${FORMULA_FILE}" << EOF
class Gostarter < Formula
  desc "A CLI tool to quickly setup fresh Go projects with a standard directory structure and files."
  homepage "https://github.com/souravlayek/gostarter"
  url "https://github.com/souravlayek/gostarter/archive/v#{version}.tar.gz"
  sha256 "${SHA256}"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = HOMEBREW_CACHE/"go_cache"
    (buildpath/"src/github.com/souravlayek/gostarter").install buildpath.children

    cd "src/github.com/souravlayek/gostarter/cmd/gostarter" do
      system "go", "build", "-o", bin/"gostarter", "."
    end
  end

  test do
    system "#{bin}/gostarter", "--help"
  end
end
EOF

echo "Homebrew formula created: ${FORMULA_FILE}"
echo "Version: ${VERSION}"
echo "SHA256: ${SHA256}"