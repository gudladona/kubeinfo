FROM debian:stable

COPY ./bin/app .

ENTRYPOINT /app