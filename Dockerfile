FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o weather-app

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/weather-app /weather-app
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/weather-app"]