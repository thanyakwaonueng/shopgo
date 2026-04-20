#!/bin/bash

# Move to the directory where this script is located, then go up one level to the root
cd "$(dirname "$0")/.."

echo "Current directory: $(pwd)"
echo "Starting Swagger documentation generation..."

# Now that we are at the root, main.go and swagger_doc are right here
swag init -g main.go -o swagger_doc

echo "Done! Documentation saved in ./swagger_doc"
