import ollama
# from transformers import pipeline

# pipe = pipeline(
#     "text-generation",
#     model="dicta-il/DictaLM-3.0-1.7B-Instruct",
#     device_map="cpu",
# )


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
