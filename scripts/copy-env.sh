#!/bin/bash

if [ -f .env ]; then
  echo ".env file exists, proceeding with the script..."
else
  if [ -f env.example ]; then
    echo ".env file is missing, copying env.example to .env..."
    cp env.example .env
    echo ".env has been created from env.example"
  else
    echo "Error: Neither .env nor env.example file exists. Please create an env.example file."
    exit 1
  fi
fi
