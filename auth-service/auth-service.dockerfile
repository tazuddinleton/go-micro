FROM alpine:latest

RUN mkdir /app

COPY ./authService /app

CMD [ "/app/authService" ]