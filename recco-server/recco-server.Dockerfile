FROM golang:1.24.1-bullseye
RUN apt update && apt install -y vim unzip

EXPOSE 8080

VOLUME /recco-server

WORKDIR /recco-server
# CMD [ "/bin/bash" ]