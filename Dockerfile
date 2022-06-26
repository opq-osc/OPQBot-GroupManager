#FROM golang:alpine AS build
#WORKDIR $GOPATH/src
FROM alpine:latest AS build
WORKDIR /apps
ARG TARGETPLATFORM
ARG BUILDPLATFORM
COPY . .
RUN ls -lh && echo $TARGETPLATFORM \
    && [[ "$TARGETPLATFORM" == "linux/amd64" ]] \
    && mv /apps/opqbot-manager-amd64 /apps/opqbot-manager || echo "not amd64" \
    && [[ "$TARGETPLATFORM" == "linux/arm64" ]] \
    && mv /apps/opqbot-manager-arm64 /apps/opqbot-manager || echo "not arm64" \
    && [[ "$TARGETPLATFORM" == "linux/arm/v7" ]] \
    && mv /apps/opqbot-manager-arm /apps/opqbot-manager || echo "not arm" \
    && [[ "$TARGETPLATFORM" == "linux/386" ]] \
    && mv /apps/opqbot-manager-386 /apps/opqbot-manager || echo "not 386"



# if [[ "$TARGETPLATFORM" = "linux/amd64" ]]; \
#     then \
#         mv ./opqbot-manager-amd64 ./opqbot-manager; \
#     fi \
#     && if ["$TARGETPLATFORM" = "linux/arm64"]; \
#     then \
#         mv ./opqbot-manager-arm64 ./opqbot-manager; \
#     fi \
#     && if ["$TARGETPLATFORM" = "linux/arm/v7"]; \
#     then \
#         mv ./opqbot-manager-arm ./opqbot-manager; \
#     fi \
#     && if ["$TARGETPLATFORM" = "linux/386"]; \
#     then \
#         mv ./opqbot-manager-386 ./opqbot-manager; \
#     fi

RUN apk add upx \
    && upx opqbot-manager \
    || echo "UPX Install Failed!"
# RUN go mod tidy\
#     && go build -o opqbot-manager -ldflags="-s -w" . \
#     && apk add upx \
#     && upx opqbot-manager \
#     || echo "UPX Install Failed!"

FROM alpine:latest
LABEL MAINTAINER enjoy<i@mcenjoy.cn>
ENV VERSION 1.0
# create a new dir
WORKDIR /apps
COPY --from=build /apps/opqbot-manager /apps/opqbot-manager
COPY config.yaml.example /apps/
COPY font.ttf /apps/
COPY dictionary.txt /apps/
# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

EXPOSE 8888

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# 开RUN
ENTRYPOINT ["/apps/opqbot-manager"]
