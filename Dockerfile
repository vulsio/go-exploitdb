FROM golang:alpine as builder

RUN apk add --no-cache make git gcc musl-dev

ENV REPOSITORY github.com/prince-chrismc/go-exploitdb
COPY . $GOPATH/src/$REPOSITORY
RUN cd $GOPATH/src/$REPOSITORY && make install

FROM alpine:3.8

MAINTAINER princechrismc

ENV LOGDIR /var/log/vuls
ENV WORKDIR /vuls

RUN apk add --no-cache ca-certificates \
    && mkdir -p $WORKDIR $LOGDIR

COPY --from=builder /go/bin/go-exploitdb /usr/local/bin/

VOLUME [$WORKDIR, $LOGDIR]
WORKDIR $WORKDIR
ENV PWD $WORKDIR

ENTRYPOINT ["go-exploitdb"]
CMD ["--help"]
