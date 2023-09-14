#!/bin/bash

# Directory where your source code is located
SRC_DIR="./cmd"

# Directory where you want to store the compiled executables
BIN_DIR="./bin"

# Ensure the bin directory exists, or create it if necessary
mkdir -p "$BIN_DIR"

# List of source files to compile (replace with your own file names)
SUB_DIRECTORIES=("coordinator" "store" "cli")

# Loop through each source file and compile it
for directory in "${SUB_DIRECTORIES[@]}"; do

    # Compile the source file and place the executable in the bin directory
    go build -o "$BIN_DIR/$directory" "$SRC_DIR/$directory/main.go"

    # Check if the compilation was successful
    if [ $? -eq 0 ]; then
        echo "Compiled $directory and moved to $BIN_DIR/$directory"
    else
        echo "Error compiling $directory"
    fi
done
