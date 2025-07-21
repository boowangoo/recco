FROM python:3.11-slim

VOLUME [ "/app/dataset" ]

COPY code/ /app/
WORKDIR /app
RUN pip install --no-cache-dir -r requirements.txt

ENTRYPOINT [ "/bin/bash", "start.sh" ]
