FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY gin/go.mod gin/go.sum ./gin/
RUN cd gin && go mod download

COPY gin/ ./gin/
COPY shared/ ./shared/

WORKDIR /build/gin
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /build/shared/templates/ ./shared/templates/

EXPOSE 8081
CMD ["./server"]
