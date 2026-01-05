#!/bin/sh

set -e # Exit immediately if a command exits with a non-zero status.

# --- Configuration ---
API_DIR="proto"
PKG_DIR="pkg/grpc"

# --- Function to compile protos in a given subdirectory ---
compile_protos() {
  local SUB_DIR="$1"
  local SRC_DIR="$API_DIR/$SUB_DIR"

  echo "--- Compiling protos for '$SUB_DIR' ---"

  # Check if the source directory exists
  if [ ! -d "$SRC_DIR" ]; then
    echo "Directory not found, skipping: $SRC_DIR"
    return
  fi

  # Check if there are any .proto files to compile
  if [ -z "$(find "$SRC_DIR" -name '*.proto' -print -quit)" ]; then
      echo "No .proto files found in $SRC_DIR, skipping."
      return
  fi

  # Find and compile all .proto files in the subdirectory
  find "$SRC_DIR" -type f -name "*.proto" | while read -r proto_file; do
    echo "  -> Compiling $proto_file"
    protoc -I="$API_DIR" \
           --go_out="$PKG_DIR" --go_opt=paths=source_relative \
           --go-grpc_out="$PKG_DIR" --go-grpc_opt=paths=source_relative \
           --experimental_allow_proto3_optional \
           "$proto_file"
  done
}

# --- Main Execution ---
echo "Starting Protobuf/gRPC code generation..."

# Ensure the root destination directory exists
mkdir -p "$PKG_DIR"

compile_protos "provide"
compile_protos "consume"
compile_protos "event"

echo "Code generation finished successfully."