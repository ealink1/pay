#!/bin/bash

# 获取当前系统的架构
ARCH=$(uname -m)

# 根据系统架构设置 GOARCH 参数
if [ "$ARCH" = "arm64" ]; then
    # Apple Silicon Mac (M1, M2, etc.)
    GO_ARCH="arm64"
    echo "检测到 Apple Silicon Mac，使用 ARM64 架构构建..."
else
    # Intel Mac
    GO_ARCH="amd64"
    echo "检测到 Intel Mac，使用 AMD64 架构构建..."
fi

# 定义 PID 文件路径
PID_FILE="./app.pid"

# 函数：停止正在运行的应用程序
stop_existing_app() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")

        # 检查进程是否存在
        if ps -p $PID > /dev/null; then
            echo "发现正在运行的应用程序 (PID: $PID)，正在停止..."
            kill $PID

            # 等待进程结束
            sleep 2

            # 检查进程是否仍然存在
            if ps -p $PID > /dev/null; then
                echo "强制终止应用程序..."
                kill -9 $PID
            fi

            # 删除 PID 文件
            rm -f "$PID_FILE"
            echo "已停止正在运行的应用程序。"
        else
            echo "PID 文件存在但应用程序未在运行，清理 PID 文件..."
            rm -f "$PID_FILE"
        fi
    else
        echo "没有发现正在运行的应用程序。"
    fi
}

# 在构建前执行 go mod tidy 来确保依赖项是最新的
echo "正在更新和整理依赖项..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "go mod tidy 执行失败，构建中止。"
    exit 1
fi

# 构建应用程序，输出文件名为 paygo
echo "正在为 macOS ($GO_ARCH) 构建应用程序..."
CGO_ENABLED=0 GOOS=darwin GOARCH=$GO_ARCH go build -ldflags "-s -w" -o paygo main.go

# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "构建成功。二进制文件 'paygo' 已创建。"

    # 在构建成功后，停止任何正在运行的实例
    echo "检查并停止任何正在运行的应用程序..."
    stop_existing_app

    # 创建日志目录（如果不存在）
    LOG_DIR="./logs"
    if [ ! -d "$LOG_DIR" ]; then
        mkdir -p "$LOG_DIR"
    fi

    # 启动应用程序并将其放入后台运行
    echo "正在后台启动应用程序..."
    nohup ./paygo > ./logs/app.log 2>&1 &

    # 获取刚刚启动的进程 PID
    APP_PID=$!

    # 将 PID 保存到文件中以便后续管理
    echo $APP_PID > ./app.pid

    # 检查应用程序是否成功启动
    sleep 2
    if ps -p $APP_PID > /dev/null; then
        echo "应用程序已在后台启动成功，PID: $APP_PID"
        echo "日志文件位于: ./logs/app.log"
        echo "PID 文件位于: ./app.pid"
    else
        echo "应用程序启动失败，请检查日志文件: ./logs/app.log"
        exit 1
    fi
else
    echo "构建失败，未停止正在运行的应用程序。"
    exit 1
fi

