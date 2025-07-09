# Using the Text Embeddings Inference (TEI)
FROM ghcr.io/huggingface/text-embeddings-inference:1.7

ENV HOST=0.0.0.0
ENV PORT=80

EXPOSE 80

# Add healthcheck for container
COPY healthcheck.sh /app/healthcheck.sh
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD /bin/bash /app/healthcheck.sh
