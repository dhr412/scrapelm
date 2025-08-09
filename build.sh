#!/bin/bash

set -e

echo "Building Docker images..."
docker compose build

echo "Starting Docker containers in the background..."
docker compose up -d

echo "Waiting for the Ollama service to start..."
while ! docker compose exec ollama ollama list >/dev/null 2>&1; do
    echo -n "."
    sleep 1
done
echo -e "\nOllama service is ready."

model=$1
if [ -z "$model" ]; then
    echo "No model specified. Exiting."
    exit 0
fi
echo "Pulling the Ollama model ($model)..."
docker compose exec ollama ollama pull "$model"

echo -e "To use the tool, run a command like this:"
echo 'docker compose exec app ./scrapelm -url "https://ollama.com" -model "gemma3:1b"  -prompt "Summarize the main features." -output-file "summary.txt"'
