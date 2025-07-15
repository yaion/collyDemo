#!/bin/bash

echo "启动考古家数据采集系统..."
echo

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go环境，请先安装Go"
    exit 1
fi

# 检查MongoDB是否运行
echo "检查MongoDB连接..."
if ! mongosh --eval "db.runCommand('ping')" &> /dev/null; then
    echo "警告: MongoDB可能未运行，请确保MongoDB服务已启动"
    echo
fi

# 安装依赖
echo "安装依赖包..."
go mod tidy

# 编译程序
echo "编译程序..."
go build -o collyDemo

# 运行程序
echo "启动程序..."
echo "按 Ctrl+C 停止程序"
echo
./collyDemo 