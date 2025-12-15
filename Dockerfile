FROM golang:1.25.5-alpine AS builder
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -trimpath \
    -o /build/bin \
    main.go

FROM alpine:3.21
ENV PORT=80
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    netcat-openbsd \
    && addgroup -g 1000 appuser \
    && adduser -D -u 1000 -G appuser appuser
WORKDIR /app
COPY --from=builder /build/bin /app/bin
COPY migrations /app/migrations
COPY entrypoint.sh /app/entrypoint.sh
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
RUN chmod +x /app/entrypoint.sh && \
    chown -R appuser:appuser /app
USER appuser
EXPOSE ${PORT}
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/bin", "serve"]
