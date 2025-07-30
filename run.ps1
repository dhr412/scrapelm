Write-Host "Building Docker images..."
docker-compose build

Write-Host "Starting Docker containers in the background..."
docker-compose up -d

Write-Host "Waiting for the Ollama service to start..."
while ($true) {
    try {
        docker-compose exec ollama ollama list | Out-Null
        break
    } catch {
        Write-Host -NoNewline "."
        Start-Sleep -Seconds 1
    }
}

$model = Read-Host -Prompt "Enter the ollama model to pull (e.g., qwen2:0.5b)"
if (-not $model) {
    Write-Host "No model specified. Exiting."
    exit 0
}
Write-Host "Pulling the Ollama model ($model)..."
docker-compose exec ollama ollama pull $model

Write-Host "`nTo use the tool, run a command like this:"
Write-Host 'docker-compose exec app python src/cli.py -url "https://ollama.com" -model "qwen3:0.6b"  -prompt "Summarize the main features." -output-file "summary.txt"'
