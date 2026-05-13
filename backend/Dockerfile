# Build stage
FROM golang:1.22-alpine AS builder
RUN apk --no-cache add git
WORKDIR /app
COPY go.mod ./
ENV GONOSUMCHECK=*
ENV GONOSUMDB=*
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download || true
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["./server"]
