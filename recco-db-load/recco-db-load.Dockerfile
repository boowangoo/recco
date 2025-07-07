FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt ./
COPY db_load.py ./
RUN pip install --no-cache-dir -r requirements.txt
ADD https://files.grouplens.org/datasets/movielens/ml-32m.zip /app/
RUN apt-get update && apt-get install -y unzip && \
    unzip /app/ml-32m.zip -d /app/ && \
    mv /app/ml-32m/ /app/dataset/ && \
    rm /app/ml-32m.zip

CMD ["python", "db_load.py"]
