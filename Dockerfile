# 使用alpine作为基础镜像
FROM alpine:latest

# 安装必要的工具（wget用于下载文件）
RUN apk update && apk add --no-cache wget && rm -rf /var/cache/apk/*
        
#设置下载目录
WORKDIR /app

# 复用架构变量，自动匹配mixapi的x86_64/aarch64版本
RUN arch="$(apk --print-arch)"; \
    case "$arch" in \
        'x86_64') \
            wget -O mixapi https://github.com/aiprodcoder/MIXAPI/releases/download/v1.2/mixapi-v1.2-linux-amd64; \
            ;; \
        'aarch64') \
            wget -O mixapi https://github.com/aiprodcoder/MIXAPI/releases/download/v1.2/mixapi-v1.2-linux-arm64; \
            ;; \
    esac

# 设置文件可执行权限
RUN chmod +x mixapi

# 暴露3000端口
EXPOSE 3000

#设置工作目录
WORKDIR /data

# 启动命令
CMD ["/app/mixapi"]

