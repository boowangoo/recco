FROM golang:1.24.1-bullseye
RUN apt update && apt install -y vim unzip
WORKDIR /recco-dataset
ADD https://www.kaggle.com/api/v1/datasets/download/rounakbanik/the-movies-dataset /dataset.zip
RUN unzip /dataset.zip -d /recco-dataset
RUN rm -f /dataset.zip
WORKDIR /recco-server