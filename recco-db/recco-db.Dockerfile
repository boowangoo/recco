# Currently using Qdrant as backend db
FROM qdrant/qdrant:v1.13.6

# Qdrant persistent data locations are specified at https://qdrant.tech/documentation/guides/configuration/
# Qdrant treats service configuration and content data as persistent data
VOLUME /qdrant/storage
VOLUME /qdrant/config

EXPOSE 6333