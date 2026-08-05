# Zephyr 网站部署总结 — 2026-08-03

## 项目结构
- 前端: `A:\ZEPHYR_S\website\client` (Vue 3 + Vite)
- 后端: `A:\ZEPHYR_S\website\server\main.go` (Go, 端口 8080)
- 远程服务器: `admin@47.116.136.145` (密钥已配)
- 域名: `zephyrzheng.cn`

## 服务器关键路径
| 用途 | 路径 |
|------|------|
| 前端静态文件 | `/var/www/zephyrzheng/` |
| 后端二进制 | `/opt/zephyr/server/server_linux`（运行时复制为 `/opt/zephyr/server/server`） |
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
### 仅改前端
1. `cd client && npm run build`
2. `ssh admin@47.116.136.145 "rm -rf ~/dist && mkdir ~/dist && sudo rm -rf /var/www/zephyrzheng/*"`
3. `scp -r client/dist/* admin@47.116.136.145:~/dist/`
4. `ssh admin@47.116.136.145 "sudo cp -r ~/dist/* /var/www/zephyrzheng/"`
5. 用户 Ctrl+Shift+R 强制刷新（Nginx 无缓存问题）

### 同时改后端
1. 本地改完 `server/main.go`
2. `cd server && $env:GOOS='linux'; $env:GOARCH='amd64'; go build -o server_linux main.go`
3. `scp server/server_linux admin@47.116.136.145:~/server_linux`
4. `ssh admin@47.116.136.145 "sudo systemctl stop zephyr-api && sudo cp ~/server_linux /opt/zephyr/server/server && sudo chmod +x /opt/zephyr/server/server && sudo systemctl start zephyr-api"`
   - **注意必须 `stop` 再 copy，不然报 `Text file busy`**

## 服务管理命令 (服务器上执行)
- 查看后端状态: `sudo systemctl status zephyr-api`
- 重启后端: `sudo systemctl restart zephyr-api`
- 查看日志: `sudo tail -50 /opt/zephyr/logs/api.log`
- 重载 Nginx: `sudo systemctl reload nginx`
- 查看端口占用: `sudo ss -tlnp | grep 8080`

## 踩坑记录 (2026-08-03)

### 坑1：`git checkout HEAD:file` 恢复代码覆盖了未提交的新功能
- 本地有很多未 commit 的新功能（子分类、年度进度、study 分类等）
- `git show HEAD:xxx | Set-Content` 恢复到旧版本后，所有新功能丢失
- **教训**：恢复单个文件前先确认 git 仓库里的版本是不是最新的；有未提交的重要改动先 `git stash` 或临时备份

### 坑2：PowerShell + Python + SSH 的三重转义地狱
- PowerShell 里写 Python 脚本，引号、`$`、`--` 都会被 PS 解析
- SSH 传 SQL 语句，单引号双引号又被 bash 解析一层
- **解决办法**：用 Python 脚本写到 `.py` 文件再执行，别直接在 CLI 里串联

### 坑3：scp + sudo cp 合并模式不覆盖旧文件
- `sudo cp -r ~/dist/* /var/www/zephyrzheng/` 不会删除目标目录的旧文件，只会合并
- 导致新旧 JS 文件共存，`index.html` 引用新 JS，但 chunk lazy-load 可能拉到旧的
- **解决办法**：先 `sudo rm -rf /var/www/zephyrzheng/*` 彻底清空再复制

### 坑4：Nginx 和浏览器双重缓存
- 前端文件更新后浏览器可能用缓存（304），Nginx 没有配置 cache-control
- **解决办法**：用户手动 Ctrl+Shift+R 强制刷新；F12 Network 勾选 Disable cache

### 坑5：前后端 API 返回格式不一致
- 公共接口 `getArticles` 返回 `{articles: [...], total: N}`（对象包数组）
- 管理接口 `adminGetArticles` 返回 `[...]`（裸数组）
- 前端取 `res.data.articles`，管理接口返回 undefined，文章列表显示空
- **解决办法**：统一所有接口返回包在对象里

### 坑6：Promise.all 请求顺序与 cats 数组不匹配
- `Promise.all([blog, leetcode, projects, notes, study])`
- `cats = [blog, leetcode, projects, study, notes]`
- `results[3]` 是 notes 但 `cats[3]` 是 study → 分类互换
- **解决办法**：两个数组顺序必须严格一致

### 坑7：Vue `<script setup>` 里 `computed` 有时不触发模板更新
- `studyColumns = computed(() => studyCatList.value.map(c => c.name))`
- 数据从 API 拿到后 `studyCatList` 更新了，但 computed 在模板的 `v-for` 里没重新渲染
- 原因可能是 Vite/Rolldown 编译优化导致 computed 依赖追踪断裂
- **解决办法**：改用 `ref([])` 并在 `fetchStudyCats` 里手动赋值 `studyColumns.value = ...`

### 坑8：`isLoggedIn` 初始化时序问题
- `useAuthStore` 的 `isLoggedIn` 从 `localStorage` 读 token 是同步的
- 但 `onMounted` 执行时 `isLoggedIn.value` 可能是 false（初始化时序）
- 导致 `if (isLoggedIn.value)` 判断失败，`fetchStudyCats` 不执行
- **解决办法**：`onMounted` 里无条件调用，不依赖 `isLoggedIn`

### 坑9：`<script setup>` 里定义的 ref 不在模板里使用就会被 tree-shake 掉
- Vite 编译时会分析哪些变量被模板引用，没引用的会被丢弃
- 只要模板里引用了（如 `v-for="sc in studyColumns"`），对应的 ref/computed 就会保留

### 坑10：Windows 下 Go 交叉编译 Linux 二进制
- `$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o server_linux main.go`
- 编译成功后 `server_linux` 可直接在 Linux 服务器运行

### 坑11：覆盖正在运行的二进制文件报 `Text file busy`
- `cp` 覆盖 `/opt/zephyr/server/server` 时如果服务在运行会失败
- **解决办法**：先 `systemctl stop` 再 copy 再 `systemctl start`

### 坑12：MySQL 密码 hash 通过 SSH 更新困难
- bcrypt hash 含 `$` 符号，经过 SSH + bash + SQL 多层转义极易出错
- **简单方案**：删掉 users 表数据后重启服务，`initDB` 自动创建 `admin/admin123`

## 日常部署一键命令（2026-08-03 更新）
### 仅前端
```powershell
cd A:\ZEPHYR_S\website\client
npm run build
ssh admin@47.116.136.145 "rm -rf ~/dist && mkdir ~/dist && sudo rm -rf /var/www/zephyrzheng/*"
scp -r dist/* admin@47.116.136.145:~/dist/
ssh admin@47.116.136.145 "sudo cp -r ~/dist/* /var/www/zephyrzheng/"
```

### 后端
```powershell
cd A:\ZEPHYR_S\website\server
$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o server_linux main.go
scp server_linux admin@47.116.136.145:~/server_linux
ssh admin@47.116.136.145 "sudo systemctl stop zephyr-api && sudo cp ~/server_linux /opt/zephyr/server/server && sudo chmod +x /opt/zephyr/server/server && sudo systemctl start zephyr-api"
```

---