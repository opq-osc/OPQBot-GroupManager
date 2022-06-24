FROM alpine:latest as build
WORKDIR /apps
COPY opqbot-manager /apps/
RUN apk add upx && upx opqbot-manager

FROM alpine:latest
LABEL MAINTAINER enjoy<i@mcenjoy.cn>
ENV VERSION 1.0
# create a new dir
WORKDIR /apps
COPY --from=build /apps/opqbot-manager /apps/opqbot-manager

COPY config.yaml.example /apps/config.yaml.example

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

EXPOSE 8888

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# 开RUN
ENTRYPOINT ["/apps/opqbot-manager"]
