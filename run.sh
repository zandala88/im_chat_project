#!/bin/bash

echo "killing process..."
nc -z localhost 8080

if [ $? -eq 0 ]; then
    echo "im 存在"
    kill $(lsof -t -i:8080)
fi

current_date=$(date +'%Y-%m-%d-%H:%M:%S')
sleep 1
mkdir -p runlog
echo "running code..."
nohup go run main.go &
echo "app is running, fetching port info..."

sleep 2
echo "port: "
lsof -i:8080