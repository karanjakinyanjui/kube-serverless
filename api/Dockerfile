FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kube-serverless-api .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/kube-serverless-api .

EXPOSE 8080 9090

CMD ["./kube-serverless-api"]
