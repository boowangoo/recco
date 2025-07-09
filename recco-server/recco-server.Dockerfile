FROM golang:1.24.1-bullseye
RUN apt update && apt install -y vim unzip

EXPOSE 80

VOLUME /recco-server

WORKDIR /recco-server

# Comment out the entrypoint for development
ENTRYPOINT [ "/bin/bash", "start.sh" ]
