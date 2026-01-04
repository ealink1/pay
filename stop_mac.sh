#!/bin/bash

# 停止后台运行的应用程序
PID_FILE="./app.pid"

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    
    # 检查进程是否存在
    if ps -p $PID > /dev/null; then
        echo "正在停止应用程序 (PID: $PID)..."
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
        echo "应用程序已停止。"
    else
        echo "应用程序未在运行或已意外终止。"
        rm -f "$PID_FILE"
    fi
else
    echo "未找到 PID 文件，应用程序可能未在运行。"
fi
