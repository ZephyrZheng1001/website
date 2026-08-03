# Zephyr 网站部署总结 — 2026-07-28

## 项目结构
- 前端: `A:\ZEPHYR_S\website\client` (Vue 3 + Vite)
- 后端: `A:\ZEPHYR_S\website\server\main.go` (Go, 端口 8080)
- 远程服务器: `admin@47.116.136.145` (密钥已配)
- 域名: `zephyrzheng.cn`

## 服务器关键路径
| 用途 | 路径 |
|------|------|
| 前端静态文件 | `/var/www/zephyrzheng/` |
| 后端二进制 | `/opt/zephyr/server/server_linux` |
| 后端日志 | `/opt/zephyr/logs/api.log` (错误: `api_error.log`) |
| Nginx 配置 | `/etc/nginx/sites-enabled/zephyrzheng` |

## 首次配置完成项
- [x] MySQL 已有 (root:zephyr123, 库: zephyr_website)
- [x] Go 1.22.2 已安装
- [x] SSH 密钥已配置 (id_ed25519)
- [x] zephyr-api systemd 服务已创建，开机自启
- [x] Nginx 已代理 `/api/` → `127.0.0.1:8080`
- [x] 默认管理员: admin / admin123

## 日常更新步骤
1. 本地修改代码后: `cd client && npm run build`
2. 上传前端: `scp -r client/dist/* admin@47.116.136.145:/var/www/zephyrzheng/`
3. 如果改了后端: 编译上传 `server_linux` 然后 `sudo systemctl restart zephyr-api`
4. 前端纯静态，上传即生效，无需重启

## 服务管理命令 (服务器上执行)
- 查看后端状态: `sudo systemctl status zephyr-api`
- 重启后端: `sudo systemctl restart zephyr-api`
- 查看日志: `sudo tail -50 /opt/zephyr/logs/api.log`
- 重载 Nginx: `sudo systemctl reload nginx`
- 查看端口占用: `sudo ss -tlnp | grep 8080`

## 已知问题 & 注意
- 老的 `server` 进程 (PID 6866) 之前占用了 8080 端口，已被 kill。如再次遇到冲突，`sudo ss -tlnp | grep 8080` 查占用然后 `sudo kill <PID>`
- 部署脚本在 `A:\ZEPHYR_S\website\deploy.sh` 供参考 (远程服务器不需要重新跑)
- Nginx 当前只有 HTTP (80端口)，未配 HTTPS，后续可加 Let's Encrypt
