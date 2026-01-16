FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/ragbench ./cmd/ragbench/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/ragbench /usr/local/bin/ragbench

ENTRYPOINT ["/usr/local/bin/ragbench"]
