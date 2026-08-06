# Zephyr Website Deployment Guide — 2026-08-05

## Project Structure
- Frontend: `A:\ZEPHYR_S\website\client` (Vue 3 + Vite + Hash Router)
- Backend: `A:\ZEPHYR_S\website\server\main.go` (Go, port 8080)
- Server: `admin@47.116.136.145` (SSH key: id_ed25519)
- Domain: `zephyrzheng.cn` (HTTPS enabled, Alibaba Cloud SSL cert)
- SSL cert: `/etc/nginx/ssl/zephyrzheng.cn.pem` + `.key`
- ICP filing: Order 2035810496966, under review (est. 16 working days)

## Server Paths
| Purpose | Path |
|---------|------|
| Frontend static files | `/var/www/zephyrzheng/` |
| Backend binary | `/opt/zephyr/server/server_linux` (NOT `server`!) |
| Backend logs | `/opt/zephyr/logs/api.log` (errors: `api_error.log`) |
| Nginx config | `/etc/nginx/sites-enabled/zephyrzheng` |
| SSL certs | `/etc/nginx/ssl/` |

## Daily Deploy

### Frontend only
```powershell
cd A:\ZEPHYR_S\website\client
npm run build
ssh admin@47.116.136.145 "rm -rf ~/dist && mkdir ~/dist"
scp -r dist/* admin@47.116.136.145:~/dist/
ssh admin@47.116.136.145 "sudo rm -rf /var/www/zephyrzheng/* && sudo cp -r ~/dist/* /var/www/zephyrzheng/"
```

### With backend changes
```powershell
cd A:\ZEPHYR_S\website\server
$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o server_linux main.go
scp server_linux admin@47.116.136.145:~/server_linux
ssh admin@47.116.136.145 "sudo systemctl stop zephyr-api && sudo cp ~/server_linux /opt/zephyr/server/server_linux && sudo chmod +x /opt/zephyr/server/server_linux && sudo systemctl start zephyr-api"
```

## CRITICAL RULES (read before any work)

1. **NEVER use `git checkout <file>`** — overwrites uncommitted features. Use `git show HEAD:file > file` instead.
2. **NEVER chain Python/SSH in PowerShell CLI** — write a `.py` file and run it.
3. **Always `rm -rf` target dir before `cp`** when deploying frontend.
4. **Always `systemctl stop` before overwriting running Go binary** (`Text file busy`).
5. **systemd runs `server_linux`, NOT `server`** — copy to `/opt/zephyr/server/server_linux`!
6. **Keep API response format consistent** — `{success: true, data: {...}}`.
7. **Vue ref arrays: use full replacement** `arr.value = [...newArr]` for reactivity.
8. **PowerShell here-strings expand `$`** — use Python to write JS files.
9. **JS regex with literal `\n` breaks Vue SFC/Babel** — avoid `[^$\n]`, use `[^$]+`.
10. **Let's Encrypt blocked by Alibaba Cloud WAF** — use Alibaba Cloud free SSL cert.
11. **MySQL `ALTER TABLE ADD COLUMN IF NOT EXISTS` not supported** — check column exists manually or run ALTER directly.
12. **When restoring files, use `git show HEAD:path > file` in PowerShell** — ensure UTF-8 encoding.

## Recent Commits (2026-08-05)
```
f639870 fix: 网站标题乱码修复 + 复制按钮适配code-block-wrapper
242a7db fix: INSERT文章参数数量修复(SQL 9个?但传了10个参数)
784d5e7 feat: 暗色模式系统跟随+防闪白+平滑过渡动画
3110cda feat: Ctrl+K搜索快捷键 + 搜索结果高亮 + 无结果推荐 + 阅读进度记忆
2688b11 fix: 搜索结果列表修复 + 高亮 + 建议
636b5bf fix: 访客统计改为IP+UA去重; 访客数放到首页stats行
03868cb feat: 服务端独立访客统计(MySQL+IP按天去重)
01f7a83 feat: 页脚加独立访客统计(localStorage按天去重)
0f2691c fix: 字数统计修复 - systemd用server_linux而非server + MySQL缺view_count列
defcb58 fix: mermaid编码改用TextEncoder; 撤销首页卡片文章数
7b0d75f feat: 返回顶部按钮、代码语言标签、首页导航卡片文章数、文章PV统计
ae7bc48 fix: Mermaid改为客户端JS渲染(CDN加载mermaid.js)
aaf6cb3 style: 全局排版优化
dda6daf feat: 文章置顶、阅读时长、图片点击放大+懒加载、时光轴、SEO优化、搜索API
4bdff92 feat: HTTPS证书配置 + 自定义404页面
```


---

## 2026-08-05 Session: Remove reading position memory + cleanup

### Changes
- **Article.vue**: Removed `saveReadPos()` function, localStorage scroll-position restore in onMounted, `saveReadPos()` call in scroll listener, `restoreReadingProgress()` call in watch
- **inject_word.py**: Fixed dedup logic ? clears old `__DAILY_WORD__` blocks before injecting, prevents accumulation
- **index.html**: Cleaned from 6 duplicate DAILY_WORD blocks down to 1

### Cleanup
- Local: Deleted 15 temp scripts (fix*.py, patch.py, test*.py, show*.py)
- Server: Removed junk from `/opt/zephyr/server/` (go.mod, go.sum, main.go, app.log, old `server` binary)

### Effect
- Every article now opens from the top; no more auto-scrolling to last read position
- App.vue reading progress bar (the colored line at top) is kept

---

## 2026-08-06 Session: Admin panel upgrade - Time filter replaces tag filter

### Changes

**Backend (server/main.go)**
- `adminGetArticles`: Added `year`, `month`, `date` query parameter support
  - `date=2026-08-01` -> `DATE(created_at) = '2026-08-01'`
  - `month=2026-08` -> `DATE_FORMAT(created_at, '%Y-%m') = '2026-08'`
  - `year=2026` -> `YEAR(created_at) = 2026`
  - No params -> returns all articles (unchanged behavior)
- Rebuilt and deployed `server_linux` binary, restarted `zephyr-api` service

**Frontend (client/src/views/Admin.vue)**
- **Removed**: Tag filter dropdown and all related JS logic (`allTagsList`, `applyArticleFilter`)
- **Added**: Time precision filter with 3 levels:
  - Dropdown: All Time / By Year / By Month / By Date
  - By Year -> shows number input (e.g. 2026)
  - By Month -> shows `<input type="month">` native picker
  - By Date -> shows `<input type="date">` native picker
- Filter triggers server request only when a value is entered (not on precision change alone)
- Switching back to "All Time" clears the filter and reloads all articles
- Tags still display on article cards but are no longer clickable for filtering
- `fetchArticles()` now only resets time filter state on full reload (no-params call)

### Issues encountered
1. **Chinese character corruption**: Python scripts piped via stdin to `python -` caused UTF-8 Chinese characters to be mangled into `?` (0x3F bytes). Root cause: PowerShell pipeline encoding. Fixed by writing Python scripts as disk files with ASCII-safe hex-escaped Chinese strings, then executing them.
2. **Vite build permission errors**: `node_modules/.vite-temp/` write denied in sandbox. Resolved with escalated `npm run build`.
3. **Tag select replacement matching**: Indentation differences (10 vs 12 spaces) and previously-removed `@change` handler caused exact string matching to fail. Fixed by extracting the actual block from the file programmatically before replacement.
4. **Flash-back bug**: Choosing a precision (e.g. "By Year") immediately triggered `fetchArticles()` which reset `timePrecision` to empty, causing the dropdown to flash back to "All Time". Fixed by making `onTimeFilterChange` only fetch when a value is entered, and `fetchArticles` only reset time state on full reload.

### Effect
- Admin article management now filters by time instead of tags
- More scalable as tag count grows
- Cleaner UI with native browser date/month pickers


---

## 2026-08-07 Session: Homepage overhaul + LeetCode integration + UI polish

### Changes

**Homepage (Home.vue)**
- Removed daily word feature entirely (template, script, inject_word.py, index.html, package.json)
- Removed navigation cards (sidebar already has them)
- Added typing animation on hero title (120ms per char)
- Merged year progress + LeetCode stats into one compact card with two progress bars
- Stats line simplified: just articles + visitors

**LeetCode Backend**
- Server-side proxy to leetcode.cn/graphql
- Returns: totalSolved, totalQuestions, realName, avatar, skillTags, reputation, acSubmissions, totalSubmissions
- Difficulty breakdown (easy/medium/hard) not available on leetcode.cn public API

**Card Compact Design**
- Blog cards: title+tags+date in one row, summary hidden, expands on hover via grid-template-rows animation
- Tags shrunk globally, tag match highlighting when filtering
- Study/Search/Timeline cards also compacted

**Card Hover Effect**
- Replaced translateY with left-border slide (::before 3px accent bar) + shadow diffusion

**Page Transition + Copy Animation**
- router-view fade/slide transition (mode out-in)
- Copy button pop animation, turns green with checkmark

**Mobile Responsiveness**
- Body font 0.92rem, line-height 1.6
- Cards 14px padding, title 0.84rem, tags 0.62rem

### Issues
1. Python Unicode escapes written literally to files, fixed with regex conversion
2. leetcode.cn field exploration: difficultyAcCount, starRating, reputationDiff not available
3. totalQues -> totalQuestions replacement cascading bug
4. Orphaned LC card styles from multiple edit cycles, fixed with full rewrite
