FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/rag-cli ./cmd/rag-cli/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/rag-cli /usr/local/bin/rag-cli

ENTRYPOINT ["/usr/local/bin/rag-cli"]
