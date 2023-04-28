FROM alpine:latest

RUN mkdir /app

COPY ./loggerService /app

CMD [ "/app/loggerService" ]