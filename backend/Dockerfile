FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /bin/parser ./cmd/parser
FROM alpine:3.21 AS api

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/api /bin/api

EXPOSE 8080

ENTRYPOINT ["/bin/api"]

FROM alpine:3.21 AS parser

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont

COPY --from=builder /bin/parser /bin/parser

ENTRYPOINT ["/bin/parser"]

