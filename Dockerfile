FROM golang:1.14

WORKDIR /go/src/app
COPY . .
RUN mkdir -p /go/src/spankes/sample
COPY . /go/src/spankes/sample

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3000

CMD ["app"]
