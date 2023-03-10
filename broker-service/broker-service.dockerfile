# Build a tiny docker image

FROM alpine:latest

RUN mkdir /app 

COPY ./brokerService /app 

CMD [ "/app/brokerService" ]

