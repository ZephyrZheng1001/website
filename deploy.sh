#!/bin/bash
# ==============================================
# Zephyr 网站服务器部署脚本
# 用法: chmod +x deploy.sh && sudo ./deploy.sh
# ==============================================

set -e

# ========== 配置区域（请修改） ==========
DOMAIN="zephyrzheng.cn"          # 你的域名
EMAIL="admin@zephyrzheng.cn"    # Let's Encrypt 通知邮箱
APP_DIR="/opt/zephyr"             # 应用目录
MYSQL_ROOT_PASS="zephyr123"       # MySQL root 密码（与 main.go 一致）

echo "========================================"
echo " Zephyr 网站部署脚本"
echo " 域名: ${DOMAIN}"
echo "========================================"

# ========== 1. 更新系统 & 安装基础依赖 ==========
echo "[1/7] 更新系统并安装依赖..."
apt-get update -y && apt-get upgrade -y
apt-get install -y nginx certbot python3-certbot-nginx mysql-server curl

# ========== 2. 安装 Go（如果未安装） ==========
echo "[2/7] 检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "安装 Go 1.23..."
    curl -LO https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
    rm go1.23.4.linux-amd64.tar.gz
fi

# ========== 3. 启动 MySQL & 创建数据库 ==========
echo "[3/7] 配置 MySQL..."
systemctl enable mysql
systemctl start mysql

# 创建数据库（如果不存在）
mysql -u root -p"${MYSQL_ROOT_PASS}" -e "CREATE DATABASE IF NOT EXISTS zephyr_website CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || {
    echo "MySQL root 密码可能不匹配，尝试无密码连接..."
    mysql -u root -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '${MYSQL_ROOT_PASS}'; FLUSH PRIVILEGES;" 2>/dev/null || true
    mysql -u root -p"${MYSQL_ROOT_PASS}" -e "CREATE DATABASE IF NOT EXISTS zephyr_website CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
}

# ========== 4. 创建目录结构 ==========
echo "[4/7] 创建目录..."
mkdir -p ${APP_DIR}/{dist,server,logs}

# ========== 5. 上传并构建后端 ==========
echo "[5/7] 构建 Go 后端..."
cd ${APP_DIR}/server
# 注意：需要先把 server/main.go 上传到服务器
# 如果已有 main.go 则编译
if [ -f "main.go" ]; then
    go mod init zephyr-website 2>/dev/null || true
    go mod tidy
    CGO_ENABLED=0 go build -o server_linux main.go
    echo "后端编译完成"
else
    echo "【警告】${APP_DIR}/server/main.go 不存在，请先上传后端代码！"
fi

# ========== 6. 配置 systemd 服务 ==========
echo "[6/7] 配置 systemd 服务..."

# 后端服务
cat > /etc/systemd/system/zephyr-api.service << EOF
[Unit]
Description=Zephyr Website API
After=network.target mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=${APP_DIR}/server
ExecStart=${APP_DIR}/server/server_linux
Restart=always
RestartSec=5
StandardOutput=append:${APP_DIR}/logs/api.log
StandardError=append:${APP_DIR}/logs/api_error.log

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable zephyr-api

# ========== 7. HTTPS 证书 + Nginx ==========
echo "[7/7] 配置 Nginx 和 HTTPS..."

# 先配置临时 HTTP（用于证书验证）
cat > /etc/nginx/sites-available/zephyr << NGINX_HTTP
server {
    listen 80;
    listen [::]:80;
    server_name ${DOMAIN};

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}
NGINX_HTTP

ln -sf /etc/nginx/sites-available/zephyr /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
mkdir -p /var/www/certbot

nginx -t && systemctl reload nginx

# 申请 Let's Encrypt 证书
echo "申请 SSL 证书..."
certbot certonly --webroot -w /var/www/certbot \
    -d ${DOMAIN} \
    --email ${EMAIL} \
    --agree-tos \
    --non-interactive

# 写入完整 Nginx 配置（带 HTTPS）
cat > /etc/nginx/sites-available/zephyr << 'NGINX_FULL'
# ---- HTTP → HTTPS ----
server {
    listen 80;
    listen [::]:80;
    server_name DOMAIN_PLACEHOLDER;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

# ---- HTTPS 主站 ----
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name DOMAIN_PLACEHOLDER;

    ssl_certificate     /etc/letsencrypt/live/DOMAIN_PLACEHOLDER/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/DOMAIN_PLACEHOLDER/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;

    add_header Strict-Transport-Security "max-age=63072000" always;

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 静态前端
    location / {
        root /opt/zephyr/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    gzip on;
    gzip_types text/css application/javascript image/svg+xml;
}
NGINX_FULL

# 替换域名占位符
sed -i "s/DOMAIN_PLACEHOLDER/${DOMAIN}/g" /etc/nginx/sites-available/zephyr

nginx -t && systemctl reload nginx

# 配置证书自动续期
echo "0 3 * * * root certbot renew --quiet && systemctl reload nginx" > /etc/cron.d/certbot-renew

echo "========================================"
echo " 部署完成！"
echo " 前端: https://${DOMAIN}"
echo " API:  https://${DOMAIN}/api/"
echo "========================================"
echo ""
echo " 后续操作："
echo " 1. 上传前端: scp -r client/dist/* user@server:${APP_DIR}/dist/"
echo " 2. 上传后端: scp server/main.go user@server:${APP_DIR}/server/"
echo " 3. 启动后端: systemctl start zephyr-api"
echo " 4. 检查状态: systemctl status zephyr-api nginx"
echo ""
