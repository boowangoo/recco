FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt ./
COPY db_load.py ./
RUN pip install --no-cache-dir -r requirements.txt

VOLUME [ "/app/dataset" ]


ENTRYPOINT ["python", "db_load.py"]
