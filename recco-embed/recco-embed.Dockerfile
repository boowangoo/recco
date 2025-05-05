FROM ghcr.io/huggingface/text-embeddings-inference:1.7


ENV MODEL_ID="BAAI/bge-large-en-v1.5"

ENTRYPOINT ["bash", "-c", "text-embeddings-router --model-id \"${MODEL_ID}\""]
