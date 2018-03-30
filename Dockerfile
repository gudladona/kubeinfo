FROM golang:1.10

WORKDIR /go/src/app
COPY . .

RUN go get -d ./...
RUN go install ./...

CMD ["app"]