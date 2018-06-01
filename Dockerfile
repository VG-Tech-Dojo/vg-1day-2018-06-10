FROM golang:1.10.2-alpine

RUN apk --no-cache update && apk --no-cache add git build-base
RUN go get -v github.com/golang/dep/cmd/dep && \
    go get -v github.com/rubenv/sql-migrate/sql-migrate

ENTRYPOINT [ "make" ]
CMD [ "run" ]
