FROM golang:1.10.2-alpine

RUN apk --no-cache update && apk --no-cache add git build-base
RUN go get -v github.com/golang/dep/cmd/dep && \
    go get -v github.com/rubenv/sql-migrate/sql-migrate

VOLUME /go/src/github.com/VG-Tech-Dojo/vg-1day-2018-06-10
WORKDIR /go/src/github.com/VG-Tech-Dojo/vg-1day-2018-06-10
ENTRYPOINT [ "make" ]
CMD [ "run" ]
