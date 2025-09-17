# 使用alpine作为基础镜像
FROM alpine:latest

# 安装必要的工具（wget用于下载文件）
RUN apk update && \
    apk add --no-cache wget && \
    # 清理缓存以减小镜像体积
    rm -rf /var/cache/apk/*

# 设置工作目录
WORKDIR /app

# 设置默认架构为amd64，允许通过--build-arg覆盖 =linux/arm64 为arm平台
ARG TARGETPLATFORM=linux/amd64

# 根据目标架构下载对应文件
RUN case "${TARGETPLATFORM}" in \
        "linux/amd64")  wget -O mixapi https://github.com/aiprodcoder/MIXAPI/releases/download/v1.1/mixapi-v1.1-linux-amd64 ;; \
        "linux/arm64")  wget -O mixapi https://github.com/aiprodcoder/MIXAPI/releases/download/v1.1/mixapi-v1.1-linux-arm64 ;; \
        *) echo "Unsupported architecture: ${TARGETPLATFORM}" && exit 1 ;; \
    esac

# 设置文件可执行权限
RUN chmod +x mixapi

# 暴露3000端口
EXPOSE 3000

# 启动命令
CMD ["./mixapi"]
