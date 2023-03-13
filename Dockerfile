FROM golang:1.19 AS builder

WORKDIR /app

ADD . .

RUN cd src && go build -o /app/confparty .

ENTRYPOINT [ "/app/confparty" ]
CMD /app/confparty --out /out
