package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func traverse(n *html.Node, w io.Writer) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "noscript", "head":
			return
		}
	}

	if n.Type == html.TextNode {
		trimmedData := strings.TrimSpace(n.Data)
		if len(trimmedData) > 0 {
			fmt.Fprint(w, trimmedData, " ")
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, w)
	}
}

func scrapeWebsite(url string) (string, error) {
	client := &http.Client{Timeout: 16 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching the URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing HTML: %w", err)
	}

	var sb strings.Builder
	traverse(doc, &sb)

	text := strings.Join(strings.Fields(sb.String()), " ")
	return text, nil
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func callOllamaAPI(model, prompt string) (string, error) {
	ollaHost := os.Getenv("OLLAMA_HOST")
	if ollaHost == "" {
		ollaHost = "http://localhost:11434"
	}
	apiURL := ollaHost + "/api/generate"

	payload := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	client := &http.Client{Timeout: 64 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error calling Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned non-200 status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading Ollama response body: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error unmarshalling Ollama response: %w", err)
	}

	return ollamaResp.Response, nil
}

func processAndRetrieve(textContent, llmModel, llmPrompt string) (string, error) {
	prompt := fmt.Sprintf(
		"Based ONLY on the following text, answer the user's question. "+
			"Do not use any outside knowledge. If the information is not present, state that. "+
			"\n\nTEXT:\n%s"+
			"\n\nUSER'S QUESTION:\n%s",
		textContent, llmPrompt,
	)

	finalAnswer, err := callOllamaAPI(llmModel, prompt)
	if err != nil {
		return "", fmt.Errorf("error during final answer retrieval: %w", err)
	}

	re := regexp.MustCompile(`(?s)<think>.*?</think>\s*`)
	finalAnswer = re.ReplaceAllString(finalAnswer, "")
	finalAnswer = strings.TrimSpace(finalAnswer)

	return finalAnswer, nil
}

func main() {
	url := flag.String("url", "", "URL of the website to scrape (required)")
	outputDir := flag.String("output-dir", "", "Directory to save intermediate scraped text. Defaults to OS temp dir")
	model := flag.String("model", "", "Name of the Ollama model to use (required)")
	prompt := flag.String("prompt", "", "The specific question for the LLM (required)")
	outputFile := flag.String("output-file", "", "File to save the final LLM response. Defaults to console")
	flag.Parse()

	if *url == "" || *model == "" || *prompt == "" {
		fmt.Println("Error: --url, --model, and --prompt are required flags.")
		flag.Usage()
		os.Exit(1)
	}

	scrapedText, err := scrapeWebsite(*url)
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	var tempFile *os.File
	var tempFilePath string

	writeDir := *outputDir
	if writeDir == "" {
		writeDir = os.TempDir()
	}

	tempFile, err = os.CreateTemp(writeDir, "scraped-*.txt")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFilePath = tempFile.Name()
	defer os.Remove(tempFilePath)

	if _, err := tempFile.WriteString(scrapedText); err != nil {
		log.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	fmt.Printf("Scraped and cleaned text saved to %s\n", tempFilePath)

	finalAnswer, err := processAndRetrieve(scrapedText, *model, *prompt)
	if err != nil {
		log.Fatalf("Processing with Ollama failed: %v", err)
	}

	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(finalAnswer), 0644)
		if err != nil {
			log.Fatalf("Failed to write to output file: %v", err)
		}
		fmt.Printf("Scraped response saved to: %s\n", *outputFile)
	} else {
		fmt.Println(finalAnswer)
	}
}
