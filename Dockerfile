FROM golang:1.19-alpine

WORKDIR /usr/src/app

RUN apk update && apk upgrade && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o YT-API .

CMD ["./YT-API"]