FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o wordle-agent .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/wordle-agent .
COPY --from=builder /app/internal/dictionary/words.txt ./internal/dictionary/

EXPOSE 5001

CMD ["./wordle-agent"]