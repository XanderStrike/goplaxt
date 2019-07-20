FROM golang:1.12.7-alpine as builder
WORKDIR $GOPATH/src/github.com/xanderstrike/goplaxt/ 
RUN apk add --no-cache git
COPY . .
RUN mkdir /out
RUN mkdir /out/keystore
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o /out/goplaxt-docker

FROM scratch
LABEL maintainer="xanderstrike@gmail.com"
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out .
COPY static ./static
VOLUME /app/keystore/
EXPOSE 8000
ENTRYPOINT ["/app/goplaxt-docker"]
