FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o orderdelay .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
USER appuser

COPY --from=builder /app/orderdelay .

EXPOSE 8888

CMD ["./orderdelay"]