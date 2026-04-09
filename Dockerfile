FROM golang:1.25-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/shortener ./cmd/shortener

CMD ["/app/shortener"]