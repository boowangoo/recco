# Using the Text Embeddings Inference (TEI)
FROM ghcr.io/huggingface/text-embeddings-inference:1.7

# 109M parameter, 768-dimensional embedding model
# Benchmarks available: https://huggingface.co/docs/text-embeddings-inference/
# BGE models are the most downloaded English models on Hugging Face Hub.
ENV MODEL_ID="BAAI/bge-large-en-v1.5"

ENTRYPOINT ["bash", "-c", "text-embeddings-router --model-id \"${MODEL_ID}\""]
