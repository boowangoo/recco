FROM golang:1.24.1-bullseye
RUN apt update && apt install -y vim unzip

EXPOSE 80

VOLUME /recco-server

WORKDIR /recco-server
# CMD [ "/bin/bash" ]