FROM python:3.14-slim

WORKDIR /recco-dataloader

RUN apt-get update && apt-get install -y unzip
RUN pip install pandas

COPY ./dataloader /recco-dataloader/

# Download the dataset for the dataloader
RUN mkdir -p /dataset
ADD https://www.kaggle.com/api/v1/datasets/download/rounakbanik/the-movies-dataset /dataset.zip
RUN unzip /dataset.zip -d /dataset && \
    rm /dataset.zip
