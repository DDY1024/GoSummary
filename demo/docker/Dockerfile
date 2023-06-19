# 使用 golang 作为基础镜像
FROM golang:latest

# 设置工作目录
WORKDIR /go/src/app

# 复制应用程序代码到镜像中
COPY . .

# 安装依赖
RUN go mod tidy
# RUN go get -d -v ./...

# 构建应用程序
RUN go build -o app

# 指定容器启动命令
CMD ["./app"]

# 暴露容器的端口（如果应用程序监听了某个端口）
EXPOSE 18080
