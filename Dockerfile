FROM pbxx/go-docker-base:master-latest as builder
COPY . /go/src/github.com/promoboxx/migrate
WORKDIR /go/src/github.com/promoboxx/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -ldflags "-s" -installsuffix cgo -o bin/migrate *.go

FROM alpine:latest

RUN apk update && \
    apk add ca-certificates bash && \
    rm -rf /var/cache/apk/*

WORKDIR /
COPY --from=builder /go/src/github.com/promoboxx/migrate/bin/migrate .
ENTRYPOINT ["/migrate"]
