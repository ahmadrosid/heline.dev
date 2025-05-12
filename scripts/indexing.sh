#!/bin/bash
set -e


cd $(pwd)/heline-indexer

# Check if the binary exists
if [ ! -f "./target/release/heline-indexer" ]; then
    echo "Binary not found. Building heline-indexer..."
    # Check if Rust is installed
    if ! command -v cargo &> /dev/null; then
        echo "Rust not found. Please install Rust first: https://rustup.rs/"
        exit 1
    fi
    
    # Build the binary in release mode
    cargo build --release
    
    echo "Build completed."
fi

# Run the indexer
echo "Running heline-indexer..."
./target/release/heline-indexer repo.json --delete-dir
