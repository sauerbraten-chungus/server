FROM golang:1.24.5 as builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server ./server
RUN chmod +x ./server
CMD ["./server"]
# CMD ["tail", "-f", "/dev/null"]
