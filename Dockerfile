FROM golang:1.19 AS builder

WORKDIR /app

ADD . .

RUN cd src && go run main.go && mkdir -p /out

CMD cp -r /app/docs/* /out/
