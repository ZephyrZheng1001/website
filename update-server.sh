#!/bin/bash
# ==============================================
# Zephyr 网站更新脚本（已有服务器，只更新代码）
# 用法: 在本地项目根目录运行
#   chmod +x update-server.sh
#   ./update-server.sh user@your-server-ip
# ==============================================

set -e

SERVER="$1"
APP_DIR="/opt/zephyr"

if [ -z "$SERVER" ]; then
    echo "用法: ./update-server.sh user@your-server"
    echo "示例: ./update-server.sh root@192.168.1.100"
    exit 1
fi

echo "========================================"
echo " 更新 Zephyr 网站到 ${SERVER}"
echo "========================================"

# ---- 1. 构建前端 ----
echo "[1/4] 构建前端..."
cd client
npm install --silent
npm run build
cd ..

# ---- 2. 上传前端 ----
echo "[2/4] 上传前端静态文件..."
ssh ${SERVER} "mkdir -p ${APP_DIR}/dist"
scp -r client/dist/* ${SERVER}:${APP_DIR}/dist/

# ---- 3. 上传并编译后端 ----
echo "[3/4] 上传并编译后端..."
ssh ${SERVER} "mkdir -p ${APP_DIR}/server"
scp server/main.go server/go.mod server/go.sum ${SERVER}:${APP_DIR}/server/

ssh ${SERVER} "cd ${APP_DIR}/server && \
    export PATH=\$PATH:/usr/local/go/bin && \
    go mod tidy && \
    CGO_ENABLED=0 go build -o server_linux main.go"

# ---- 4. 重启服务 ----
echo "[4/4] 重启服务..."
ssh ${SERVER} "systemctl restart zephyr-api && systemctl reload nginx"

echo "========================================"
echo " 更新完成！"
echo " 访问: https://${SERVER#*@}"
echo "========================================"
