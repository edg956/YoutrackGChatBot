FROM golang:1.17.2 AS BUILDER

WORKDIR /work/build

COPY . .

RUN go build .

FROM ubuntu:focal-20210921

USER gchatbot
WORKDIR /home/gchatbot

COPY --from=BUILDER --chown=gchatbot:gchatbot server .

ENTRYPOINT ./server