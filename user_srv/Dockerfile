# 启动编译环境
FROM golang:1.19-alpine AS builder

# 配置编译环境
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 拷贝源代码到镜像中
COPY . /go/src/shop/user_srv
COPY config.yaml /go/bin
# 编译
WORKDIR /go/src/shop
RUN go install ./user_srv/...

#FROM  alpine:3.18
#COPY --from=builder /go/bin/user_srv /bin/user_srv


# 申明暴露的端口
EXPOSE 8080

# 设置服务入口
ENTRYPOINT [ "/bin/user_srv" ]
