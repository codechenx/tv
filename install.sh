#!/usr/bin/env bash
set -e

# ftv installer script
# Automatically detects platform and installs the latest release

get_latest_release() {
    curl --silent "https://api.github.com/repos/codechenx/FastTableViewer/releases/latest" |
        grep '"tag_name":' |
        sed -E 's/.*"v([^"]+)".*/\1/'
}

detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        linux*)   os="Linux" ;;
        darwin*)  os="Darwin" ;;
        *)        echo "Unsupported OS: $os" >&2; exit 1 ;;
    esac
    
    case "$arch" in
        x86_64)  arch="x86_64" ;;
        amd64)   arch="x86_64" ;;
        arm64)   arch="arm64" ;;
        aarch64) arch="arm64" ;;
        armv7l)  arch="armv7" ;;
        i386|i686) arch="i386" ;;
        *)       echo "Unsupported architecture: $arch" >&2; exit 1 ;;
    esac
    
    echo "${os}_${arch}"
}

main() {
    echo "ftv installer"
    echo "============="
    echo
    
    # Check if ftv already exists
    if [ -f "$PWD/ftv" ]; then
        echo "Error: ftv already exists in current directory" >&2
        echo "Please remove it first or install to a different location" >&2
        exit 1
    fi
    
    # Detect platform
    platform=$(detect_platform)
    echo "Detected platform: $platform"
    
    # Get latest version
    version=$(get_latest_release)
    if [ -z "$version" ]; then
        echo "Error: Could not determine latest version" >&2
        exit 1
    fi
    echo "Latest version: v$version"
    
    # Construct download URL
    filename="FastTableViewer_${version}_${platform}.tar.gz"
    url="https://github.com/codechenx/FastTableViewer/releases/download/v${version}/${filename}"
    
    echo "Downloading: $url"
    
    # Create temp directory
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT
    
    # Download and extract
    if ! curl -sSL "$url" | tar -xz -C "$tmp_dir"; then
        echo "Error: Download or extraction failed" >&2
        exit 1
    fi
    
    # Move binary to current directory
    mv "$tmp_dir/FastTableViewer" "$PWD/ftv"
    chmod +x "$PWD/ftv"
    
    echo
    echo "✓ Successfully installed ftv to: $PWD/ftv"
    echo
    echo "To make it globally accessible, run:"
    echo "  sudo mv ftv /usr/local/bin/"
    echo
    echo "Or add the current directory to your PATH:"
    echo "  export PATH=\"\$PATH:$PWD\""
}

main "$@"
