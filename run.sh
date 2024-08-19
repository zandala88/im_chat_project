#!/bin/bash

# 查找占用8081端口的进程ID
pid=$(lsof -t -i:8081)

# 如果进程存在，终止该进程
if [ -n "$pid" ]; then
  echo "Killing process on port 8081 (PID: $pid)"
  kill -9 $pid
else
  echo "No process is running on port 8081"
fi

# 切换到 im_api 目录
cd "$(dirname "$0")/im_api" || { echo "Failed to switch to im_api directory"; exit 1; }

# 运行 bee run
echo "Running bee run in $(pwd)..."
bee run