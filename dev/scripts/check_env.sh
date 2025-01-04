#!/bin/sh

BASE_DIR="$1"
ENV_FILE="${BASE_DIR}/.env"
ENV_EXAMPLE_FILE="${BASE_DIR}/env.example"

# Check if .env file exists
if [ -f "${ENV_FILE}" ]; then
  echo "✅ .env file exists, proceeding with the script..."
else
  # Check if env.example exists
  if [ -f "${ENV_EXAMPLE_FILE}" ]; then
    echo "⚠️  .env file is missing, copying env.example to .env..."
    cp "${ENV_EXAMPLE_FILE}" "${ENV_FILE}"
    echo "✅ .env has been created from env.example"
  else
    echo "❌ Error: Neither .env nor env.example file exists. Please reinstall the project."
    exit 1
  fi
fi
