FROM golang:1.23-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -o redis_worker ./cmd/worker

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/redis_worker .
COPY --from=builder /app/data ./data
COPY --from=builder /app/config ./config

# Create startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
  echo '/app/redis_worker &' >> /app/start.sh && \
  echo 'exec /app/main' >> /app/start.sh && \
  chmod +x /app/start.sh

EXPOSE 8080

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/app/start.sh"]