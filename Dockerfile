FROM golang:1.16.0-alpine3.13
RUN apk add --no-cache git
RUN apk add --no-cache sqlite-libs sqlite-dev
RUN apk add --no-cache build-base
ENV CGO_ENABLED 1
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go test ./...
CMD ["vroomberg"]
