FROM golang:1.21.5

WORKDIR /coredog

COPY . .

RUN go build -o coredog cmd/main.go

CMD ["/coredog/coredog"]
