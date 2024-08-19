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

# 运行 bee run 并将输出重定向到 bee_run.log 文件, 后台运行
echo "Running bee run in $(pwd)..."
nohup bee run > bee_run.log 2>&1 &

# 输出后台运行的进程ID
echo "bee run started with PID: $!"
