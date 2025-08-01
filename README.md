# ScrapeLM

`ScrapeLM` is a command-line tool that scrapes a website's text content and uses a local LLM via Ollama to answer questions based on the scraped information.

## Features

- **Web Scraping**: Extracts the main text content from a given URL, removing unnecessary HTML tags like `<script>` and `<style>`.
- **LLM Integration**: Connects to a local Ollama instance to process the scraped text.
- **Question & Answer**: Allows you to ask a specific question about the content of the URL.
- **Local First**: Operates entirely on your local machine, ensuring privacy and offline capability.
- **Flexible I/O**: Supports saving the scraped text and the final LLM response to files.

---

## Installation

### Docker (Recommended)

This is the easiest way to get started, as it handles all dependencies, including Ollama itself.

1. **Prerequisites**: Make sure you have [Docker](https://docs.docker.com/engine/install) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
2. **Clone the repository:**

    ```bash
    git clone https://github.com/dhr412/scrapelm.git
    cd scrapelm
    ```

3. **Run the setup script**:
    - For Windows (PowerShell):

      ```powershell
      ./run.ps1 gemma3:1b
      ```

    - For Linux/macOS:

      ```bash
      chmod +x run.sh
      ./run.sh gemma3:1b
      ```

    The script will build the Docker containers, start them and pull the Ollama model passed in the argument

### Local Installation

Follow these steps if you prefer to run the application without Docker.

1. **Clone the repository:**

    ```bash
    git clone https://github.com/dhr412/scrapelm.git
    cd scrapelm
    ```

2. **Install Python dependencies:**

    ```bash
    pip install -r requirements.txt
    ```

3. **Install and run Ollama:**
    - Download and install Ollama from the [official website](https://ollama.com/).
    - Pull a model to use. For example, to get the `gemma3:1b` model, run:

      ```bash
      ollama pull gemma3:1b
      ```

---

## Usage

### With Docker

Once the containers are running, you can execute commands like this:

**Basic Example**

```bash
docker-compose exec app python src/cli.py -url "https://example.com" -model "gemma3:1b" -prompt "What is this page about?"
```

**Saving the Output**

To save the LLM's response to a file in your local directory:

```bash
docker-compose exec app python src/cli.py -url "https://example.com" -model "gemma3:1b" -prompt "Summarize the main points." -output-file "summary.txt"
```

### Local Usage

Run the script from the root of the project directory. The tool requires a URL, a model name, and a prompt.

**Basic Example**

To scrape a website and ask a question, printing the output to the console:

```bash
python src/cli.py -url "https://example.com" -model "gemma3:1b" -prompt "What is this page about?"
```

**Saving the Output**

To save the LLM's response to a file:

```bash
python src/cli.py -url "https://example.com" -model "gemma3:1b" -prompt "Summarize the main points." -output-file "summary.txt"
```

---

## CLI Flags

| Flag            | Description                                                                 | Required |
|-----------------|-----------------------------------------------------------------------------|----------|
| `-url`          | The URL of the website to scrape.                                           | Yes      |
| `-model`        | The name of the Ollama model to use (e.g., `gemma3:1b`, `llama3`).               | Yes      |
| `-prompt`       | The specific question or instruction for the LLM.                           | Yes      |
| `-output-dir`   | Directory to save the intermediate scraped text file. Defaults to a temp dir. | No       |
| `-output-file`  | File to save the final LLM response. Defaults to printing to the console.   | No       |

## License

This project is licensed under the MIT license.
