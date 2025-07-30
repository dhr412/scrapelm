import re
import requests
from bs4 import BeautifulSoup
from pathlib import Path

def scrape_website(url: str, output_filepath: Path) -> None:
    try:
        response = requests.get(url, timeout=16)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"Error fetching the URL: {e}")
        raise

    soup = BeautifulSoup(response.content, "html.parser")

    for tag in soup(["script", "style", "noscript", "head"]):
        tag.decompose()

    if soup.body:
        html_clean = soup.body.get_text(separator=" ", strip=True)
    else:
        html_clean = soup.get_text(separator=" ", strip=True)

    text = re.sub(r"\s+", " ", html_clean).strip()

    text = re.sub(r"\n\s*\n", "\n", text)

    try:
        with open(output_filepath, 'w', encoding="utf-8") as f:
            f.write(text)
        print(f"Scraped and cleaned text saved to {output_filepath}")
    except IOError as e:
        print(f"Error writing to file: {e}")
        raise
