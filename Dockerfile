FROM alpine:latest

RUN addgroup freezer && adduser -H -D -s /bin/false -G freezer freezer
RUN mkdir -p /data && chown -R freezer:freezer /data
RUN mkdir /go && chown -R freezer:freezer /go

VOLUME /data

RUN apk update && apk add go \
    git \
    make \
    python py-setuptools py-pip \
    librsvg

COPY . /go/src/github.com/fhoehl/sunnyday4pm/

ENV GOPATH /go
ENV REDIS_ADDR db:6379
ENV REDIS_HOST db
ENV REDIS_PORT 6379

# Help with a bug on OSX
ENV GODEBUG netdns=cgo

RUN go get -v github.com/fhoehl/sunnyday4pm/makeicecream
RUN go get -v github.com/fhoehl/sunnyday4pm/icecreamd

WORKDIR /go/src/github.com/fhoehl/sunnyday4pm/

RUN pip install -r requirements.txt

RUN make

RUN cp ./bin/* /usr/local/bin/
RUN cp ./scripts/* /usr/local/bin/

USER freezer
