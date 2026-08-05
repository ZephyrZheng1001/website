# Zephyr 网站部署总结 — 2026-08-05

## 项目结构
- 前端: `A:\ZEPHYR_S\website\client` (Vue 3 + Vite + hash router)
- 后端: `A:\ZEPHYR_S\website\server\main.go` (Go, 端口 8080)
- 远程服务器: `admin@47.116.136.145` (密钥 id_ed25519)
- 域名: `zephyrzheng.cn` (HTTPS 已启用, 阿里云 SSL 证书)
- HTTPS 证书路径: `/etc/nginx/ssl/zephyrzheng.cn.pem` + `.key`
- ICP 备案: 订单 2035810496966, 管局审核中(预估16工作日), 目前可正常访问

## 服务器关键路径
| 用途 | 路径 |
|------|------|
| 前端静态文件 | `/var/www/zephyrzheng/` |
| 后端二进制 | `/opt/zephyr/server/server_linux`（运行时复制为 `/opt/zephyr/server/server`） |
| 后端日志 | `/opt/zephyr/logs/api.log` (错误: `api_error.log`) |
| Nginx 配置 | `/etc/nginx/sites-enabled/zephyrzheng` |
| SSL 证书 | `/etc/nginx/ssl/` |

## 日常部署

### 仅前端（一键）
```powershell
cd A:\ZEPHYR_S\website\client
npm run build
ssh admin@47.116.136.145 "rm -rf ~/dist && mkdir ~/dist"
scp -r dist/* admin@47.116.136.145:~/dist/
ssh admin@47.116.136.145 "sudo rm -rf /var/www/zephyrzheng/* && sudo cp -r ~/dist/* /var/www/zephyrzheng/"
```

### 同时改后端
```powershell
cd A:\ZEPHYR_S\website\server
$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o server_linux main.go
scp server_linux admin@47.116.136.145:~/server_linux
ssh admin@47.116.136.145 "sudo systemctl stop zephyr-api && sudo cp ~/server_linux /opt/zephyr/server/server && sudo chmod +x /opt/zephyr/server/server && sudo systemctl start zephyr-api"
```

## 服务管理命令 (服务器上执行)
- 查看后端状态: `sudo systemctl status zephyr-api`
- 重启后端: `sudo systemctl restart zephyr-api`
- 查看日志: `sudo tail -50 /opt/zephyr/logs/api.log`
- 重载 Nginx: `sudo systemctl reload nginx`
- 查看端口占用: `sudo ss -tlnp | grep 8080`

---

## 2026-08-05 新增功能总览

### 标签筛选修复
- Blog.vue: 侧边栏标签点击通过 `?tag=xxx` 筛选文章
- fetchArticles 从 route.query.tag 读取 tag 传给 API
- watch 监听 [category, route.query.tag] 变化触发刷新
- 侧边栏标签从全部分类文章收集（不受当前筛选影响）
- 点击文章卡片内标签也可跳转筛选

### About 头像
- `<div class="avatar">Z</div>` 替换为证件照 `avatar.jpg`
- 图片放在 `client/public/avatar.jpg`, 100x100 圆形裁剪

### KaTeX 内联公式渲染
- `processInlineKatex()` 在 marked 渲染后处理 `$...$` 和 `$$...$$`
- 使用 `katex.renderToString()` 渲染为 HTML
- **坑**: `[^$
]` 中的字面 `
` 会切断 Vue SFC 解析, 改用 `[^$]+`

### Mermaid 图表
- 通过 mermaid.ink 服务渲染为图片（在 marked code renderer 中处理）
- injectMermaid() 已清空（不再需要前端 JS 渲染）

### 文章置顶
- Admin.vue 编辑页加 "📌 置顶文章" checkbox
- editorForm.is_pinned 字段
- 置顶文章在各列表最前显示 📌 标记

### 阅读时长
- `readingTime(content)` 函数: 中文字数/400 + 英文单词/200
- 中文正则: `[一-鿿㐀-䶿豈-﫿]`
- 文章列表 + 文章详情页显示 "X min"

### 时光轴 / 首页最近更新
- 首页 Hero 下方的"最近更新"改为按年月分组的文章时间线
- 每条文章显示分类色点 + 分类标签 + 文章标签 + 日期
- 独立的 Timeline.vue 页面保留（路由存在但导航栏入口已移除）

### 图片点击放大 + 懒加载
- Article.vue: `injectImageZoom()` 给所有 `<img>` 加 `loading="lazy"`
- 点击图片弹出全屏 overlay（黑色半透明背景），再点击关闭
- 跳过 KaTeX 公式图和 Mermaid 图

### 搜索增强
- 前端 `articleAPI.search(q)` → `GET /api/search?q=xxx`
- 后端已有 searchArticles handler（搜索 title/content/summary/tags）
- Search.vue 页面结果列表显示分类图标和标签

### Admin 文章管理增强
- 搜索框（按标题/内容/摘要实时筛选）
- 分类下拉筛选 + 标签下拉筛选（标签从全部文章中自动收集）
- 列表显示分类标签和文章标签，点击标签直接筛选
- 显示筛选结果数量

### 学习笔记独立
- Blog.vue 侧边栏移除"学习笔记"，只保留 4 个分类（技术文章/算法笔记/项目/碎碎念）
- 学习笔记通过 `/study` 独立页面访问，有自己的子分类侧边栏

### 全局排版优化
- 正文行高 1.7 → 1.82，卡片内边距 22 → 24px
- 文章内容行高 1.85 → 1.92，段落间距 1 → 1.3em
- 标题间距加大，列表/引用块间距加大
- 移动端适配保持（卡片 18px，触控友好）

### HTTPS + 自定义 404
- 阿里云免费 SSL 证书（`.pem` + `.key`），文件在 `A:\ZEPHYR_S\website\ssl\`
- Nginx 配置: HTTP → HTTPS 301 跳转，TLSv1.2/1.3
- `client/public/404.html` — 访问不存在路径显示返回首页按钮

### 分类计数丢失修复
- `allColumns` 从裸数组改为 `ref([...])`，更新时整体替换触发响应式
- `fetchColumnCounts` 用 `Promise.allSettled` 容错
- watch `route.fullPath` 无条件调用 fetchColumnCounts

### SEO 优化
- index.html: title "Zephyr - 郑智毅的个人网站"
- meta description, keywords, Open Graph 标签
- lang="zh-CN"

---

## 核心注意事项（务必遵守）

1. **不要用 `git checkout <file>` 恢复单个文件** — 会覆盖未提交的新功能。用 `git show HEAD:file > file` 代替
2. **不要在 PowerShell CLI 里串联 Python/SSH** — `$` 引号等会被多重转义搞乱。写 `.py` 文件再 `python xxx.py` 执行
3. **部署前端先 `rm -rf` 目标目录再 `cp`** — `cp -r` 不会删除旧文件，新旧 JS 共存会导致 chunk 加载错乱
4. **覆盖运行中的 Go 二进制先 stop** — `cp` 覆盖 `/opt/zephyr/server/server` 时会报 `Text file busy`
5. **前后端 API 返回格式统一** — 所有接口返回 `{success: true, data: {...}}`，前端取 `res.data.xxx`
6. **Promise.all 和 cats 数组顺序必须严格一致**
7. **Vue ref 包裹的数组修改内部属性可能不触发更新** — 用整体替换 `arr.value = [...newArr]`
8. **PowerShell here-string 中的 `$` 会被变量展开** — 大段 JS 代码用 Python 写入文件
9. **Python raw string 中 `📝` 不会被解析为 emoji** — 直接在字符串里写 emoji 字符
10. **JS 正则中 `[^$
]` 的字面 `
` 会切断 Vue SFC/Babel 解析** — 改用 `[^$]+`
11. **Let's Encrypt 在阿里云国内服务器上可能被 WAF 拦截** — 用阿里云免费 SSL 证书替代
12. **SPA + try_files 时 404 页面需要 Nginx 特殊配置** — `error_page 404 /404.html`

## Git 提交历史（2026-08-05）
```
6404eec fix: 阅读时长修复中文检测; 首页最近更新改为时光轴; 导航栏移除时光轴入口
0c43f41 fix: 首页完整重写 - 时光轴替换最近更新
aaf6cb3 style: 全局排版优化 - 增大行高、间距、卡片内边距
dda6daf feat: 文章置顶、阅读时长、图片点击放大+懒加载、时光轴、SEO优化、搜索API
a291bfc fix: allColumns用ref包裹+整体替换数组确保响应式更新
b775bd3 fix: 计数始终刷新; 学习笔记移出文章侧边栏独立展示
f324b17 fix: 修复分类计数丢失; Admin加搜索/分类/标签筛选
52c697d fix: 修复标签筛选bug、About头像、KaTeX内联公式渲染、SESSION_SUMMARY乱码
4bdff92 feat: HTTPS证书配置 + 自定义404页面
```
