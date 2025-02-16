FROM golang:1.23.2

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go build -o /build ./cmd/avito-test \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]