import argparse
import tempfile
from pathlib import Path

from scrape import scrape_website
from process_ollm import process_and_retrieve

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Scrape a website and use an Ollama LLM to process the content."
    )
    parser.add_argument(
        "-url",
        required=True,
        help="The URL of the website to scrape."
    )
    parser.add_argument(
        "-output-dir",
        help="Directory to save the intermediate scraped text file. Defaults to a temporary directory."
    )
    parser.add_argument(
        "-model",
        required=True,
        help="The name of the Ollama model to use (e.g., 'llama2', 'phi3')."
    )
    parser.add_argument(
        "-prompt",
        required=True,
        help="The specific question or retrieval instruction for the LLM."
    )
    parser.add_argument(
        "-output-file",
        help="File to save the final LLM response. Defaults to printing to console."
    )

    args = parser.parse_args()

    try:
        if args.output_dir:
            output_dir = Path(args.output_dir)
            output_dir.mkdir(parents=True, exist_ok=True)
            temp_file = tempfile.NamedTemporaryFile(
                delete=False, mode='w+', suffix=".txt", dir=output_dir, encoding='utf-8'
            )
            output_filepath = Path(temp_file.name)
            temp_file.close()
        else:
            with tempfile.NamedTemporaryFile(
                delete=False, mode='w+', suffix=".txt", encoding='utf-8'
            ) as temp_file:
                output_filepath = Path(temp_file.name)

        scrape_website(args.url, output_filepath)

        final_answer = process_and_retrieve(
            output_filepath,
            args.model,
            args.prompt
        )

        if args.output_file:
            with open(args.output_file, 'w', encoding='utf-8') as f:
                f.write(final_answer)
            print(f"Scraped response saved to: {args.output_file}")
        else:
            print(final_answer)

    except Exception as e:
        print(f"An error occurred in the main pipeline: {e}")
    finally:
        if "output_filepath" in locals() and output_filepath.exists():
            try:
                output_filepath.unlink()
            except OSError as e:
                print(f"Error cleaning up temporary file {output_filepath}: {e}")
