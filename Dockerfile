FROM golang:latest

EXPOSE 9001 9002

LABEL maintainer="ssubedir <ssubedir@gmail.com>"

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build

CMD ["./hwatchdog-s"]