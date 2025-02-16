FROM golang:1.24-alpine as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o xlog ./cmd/xlog

FROM alpine as final
COPY --from=builder /app/xlog /bin/xlog

CMD ["xlog", "-bind", "0.0.0.0:3000", "-source", "/files"]
