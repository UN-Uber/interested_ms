
FROM golang:1.17.4
RUN mkdir /app
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main .
WORKDIR /dist
RUN cp /app/main .
ENTRYPOINT ["/dist/main"]
