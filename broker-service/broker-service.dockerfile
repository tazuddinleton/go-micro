# Base Go Image

FROM golang:1.20-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerService ./cmd/api

RUN chmod +x /app/brokerService


# Build a tiny docker image

FROM alpine:latest

RUN mkdir /app 

COPY --from=builder /app/brokerService /app 

CMD [ "/app/brokerService" ]

