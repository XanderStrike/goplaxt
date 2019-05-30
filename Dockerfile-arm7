FROM golang:alpine as builder
WORKDIR $GOPATH/src/github.com/xanderstrike/goplaxt/ 
COPY . .
RUN mkdir /out
RUN mkdir /out/keystore
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o /out/goplaxt-docker

FROM scratch
LABEL maintainer="xanderstrike@gmail.com"
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out .
COPY static ./static
VOLUME /app/keystore/
EXPOSE 8000
ENTRYPOINT ["/app/goplaxt-docker"]
