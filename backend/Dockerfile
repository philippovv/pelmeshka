FROM golang:1.22.1 as builder
WORKDIR /app
COPY . .
RUN go mod download && go mod verify 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend ./cmd/api

FROM alpine:3.19 
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/backend .
EXPOSE 8081
CMD ["./backend"]