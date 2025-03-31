FROM golang:1.24.1-bullseye
RUN apt update && apt install -y vim
WORKDIR /recco-server

