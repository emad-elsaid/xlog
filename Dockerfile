FROM golang:1.17
ENV CGO_ENABLED=0
WORKDIR /xlog
COPY . .
RUN  go build -o xlog -installsuffix cgo *.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /xlog/xlog /bin/xlog
WORKDIR /srv

CMD ["xlog"]
