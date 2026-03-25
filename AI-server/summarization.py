import ollama


def summarize(text, model):
    model = model
    resp = ollama.chat(
        model=model,
        messages=[
            {
                "role": "system",
                "content": """
You are a meeting transcript summarizer.

Rules:
- Summarize ONLY what appears in the transcript.
- Do NOT invent information.
- Do NOT translate the text.
- Keep the original language of the transcript.
- Be concise and structured.

Output format (keep EXACT structure):

SUBJECT:
[1-2 sentence summary of the topic]

CONCLUSIONS_AND_ACTIONS:
- [clear action or conclusion]
- [clear action or conclusion]
- [clear action or conclusion]

If no actions exist, write:
NO_ACTION_ITEMS
""",
            },
            {
                "role": "user",
                "content": f"Summarize the following transcript:\n{text}",
            },
        ],
        options={"temperature": 0.3, "num_ctx": 4096},
    )
    return resp["message"]["content"]


def translate(text, model):
    response = ollama.generate(
        model=model,
        prompt=f"Translate the following English text to Hebrew:\n{text}",
        options={
            "temperature": 0.0,
            "num_predict": 1024,
        },
    )
    return response["response"].strip()


def generate_name(summery: str, model: str):
    response = ollama.generate(
        model=model,
        system="You are a professional assistant. Your task is to name meetings Based on their summary. "
        "Provide only the title. Use 5 words or less. No preamble.",
        prompt=f"Summary: {summery}",
        options={"temperature": 0.3},
    )
    return " ".join(response["response"].strip().split(" ")[:5])
