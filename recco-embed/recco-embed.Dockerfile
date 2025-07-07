# Using the Text Embeddings Inference (TEI)
FROM ghcr.io/huggingface/text-embeddings-inference:1.7

# 335M parameter, 1024-dimensional embedding model
# Benchmarks available: https://huggingface.co/docs/text-embeddings-inference/
# BGE models are the most downloaded English models on Hugging Face Hub.
ENV MODEL_ID="/data/bge-large-en-v1.5"

# Location of the embedding model 
VOLUME /data

ENV HOST=0.0.0.0
ENV PORT=80

EXPOSE 80
