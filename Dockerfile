FROM golang:1.16-alpine as build

ADD . /usr/local/go/src/github.com/wethedevelop/account

WORKDIR /usr/local/go/src/github.com/wethedevelop/account

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server

FROM alpine:3.7

ENV REDIS_ADDR=""
ENV REDIS_PW=""
ENV REDIS_DB=""
ENV MysqlDSN=""
ENV GIN_MODE="release"
ENV PORT=3000

RUN echo "http://mirrors.aliyun.com/alpine/v3.7/main/" > /etc/apk/repositories && \
    apk update && \
    apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    apk add ca-certificates && \
    echo "hosts: files dns" > /etc/nsswitch.conf

WORKDIR /www

COPY --from=build /usr/local/go/src/github.com/wethedevelop/account /usr/bin/api_server

RUN chmod +x /usr/bin/api_server

ENTRYPOINT ["api_server"]