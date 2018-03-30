FROM scratch

COPY ./bin/app .

CMD ["/app"]