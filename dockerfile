FROM golang:1.16.2-buster

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build main.go
EXPOSE 8080
CMD [ "./main" ]