# Build Geth in a stock Go builder container
FROM golang:1.14-alpine as builder


RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add --no-cache make gcc musl-dev linux-headers git

ENV GOPROXY https://goproxy.cn

ADD . /go-ethereum
RUN cd /go-ethereum/cmd/geth && go build -tags sm2
RUN cd /go-ethereum/cmd/bootnode && go build -tags sm2

# 打包, 使用 alpine 作为基础镜像
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add --no-cache ca-certificates

COPY --from=builder /go-ethereum/cmd/geth/geth /usr/local/bin/
COPY --from=builder /go-ethereum/cmd/bootnode/bootnode /usr/local/bin/
COPY --from=builder /go-ethereum/genesis.json /
RUN geth  --nousb init genesis.json
RUN bootnode -genkey=boot.key

