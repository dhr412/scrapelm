import json, os, requests, re
from pathlib import Path

OLLAMA_HOST = os.environ.get("OLLAMA_HOST", "http://localhost:11434")
OLLAMA_API_URL = f"{OLLAMA_HOST}/api/generate"

def _call_ollama_api(model: str, prompt: str) -> str:
    headers = {"Content-Type": "application/json"}
    payload = {
        "model": model,
        "prompt": prompt,
        "stream": False
    }

    try:
        response = requests.post(OLLAMA_API_URL, headers=headers, data=json.dumps(payload), timeout=64)
        response.raise_for_status()
        response_data = response.json()
        return response_data['response']
    except requests.exceptions.RequestException as e:
        print(f"Error calling Ollama API: {e}")
        raise
    except json.JSONDecodeError:
        print("Error: Failed to decode JSON response from Ollama API.")
        raise
    except KeyError:
        print("Error: 'response' key not found in Ollama API result.")
        raise

def process_and_retrieve(text_filepath: Path, llm_model: str, llm_prompt: str) -> str:
    try:
        with open(text_filepath, 'r', encoding="utf-8") as f:
            text_content = f.read()
    except FileNotFoundError:
        print(f"Error: The file {text_filepath} was not found.")
        raise
    except IOError as e:
        print(f"Error reading file {text_filepath}: {e}")
        raise

    prompt = (
        "Based ONLY on the following text, answer the user's question. "
        "Do not use any outside knowledge. If the information is not present, state that. "
        "\n\nTEXT:\n" + text_content +
        "\n\nUSER'S QUESTION:\n" + llm_prompt
    )

    try:
        final_answer = _call_ollama_api(llm_model, prompt)
        final_answer = re.sub(r"<think>.*?</think>\s*", "", final_answer, flags=re.DOTALL).strip()
        return final_answer
    except Exception as e:
        print(f"An unexpected error occurred during final answer retrieval: {e}")
        raise
