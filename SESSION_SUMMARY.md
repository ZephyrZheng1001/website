# Zephyr Website Deployment Guide — 2026-08-05

## Project Structure
- Frontend: `A:\ZEPHYR_S\website\client` (Vue 3 + Vite + Hash Router)
- Backend: `A:\ZEPHYR_S\website\server\main.go` (Go, port 8080)
- Server: `admin@47.116.136.145` (SSH key: id_ed25519)
- Domain: `zephyrzheng.cn` (HTTPS enabled, Alibaba Cloud SSL cert)
- SSL cert path on server: `/etc/nginx/ssl/zephyrzheng.cn.pem` + `.key`
- ICP filing: Order 2035810496966, under review (est. 16 working days), site accessible now

## Server Paths
| Purpose | Path |
|---------|------|
| Frontend static files | `/var/www/zephyrzheng/` |
| Backend binary | `/opt/zephyr/server/server_linux` (copied to `server` at runtime) |
| Backend logs | `/opt/zephyr/logs/api.log` (errors: `api_error.log`) |
| Nginx config | `/etc/nginx/sites-enabled/zephyrzheng` |
| SSL certs | `/etc/nginx/ssl/` |

## Daily Deploy

### Frontend only (one-liner)
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
ssh admin@47.116.136.145 "sudo systemctl stop zephyr-api && sudo cp ~/server_linux /opt/zephyr/server/server && sudo chmod +x /opt/zephyr/server/server && sudo systemctl start zephyr-api"
```

## Server Management
- Backend status: `sudo systemctl status zephyr-api`
- Restart backend: `sudo systemctl restart zephyr-api`
- View logs: `sudo tail -50 /opt/zephyr/logs/api.log`
- Reload Nginx: `sudo systemctl reload nginx`
- Check port: `sudo ss -tlnp | grep 8080`

---

## Features Added on 2026-08-05

### Tag Filtering Fix (Blog.vue)
- Sidebar tags filter articles via `?tag=xxx` in URL
- `fetchArticles()` reads `route.query.tag` and passes to API
- `watch([category, route.query.tag])` triggers refresh
- Tag list collected from ALL articles (not filtered subset)
- Clicking tags inside article cards also navigates to filter

### About Page Avatar
- Replaced `<div class="avatar">Z</div>` with photo `avatar.jpg`
- Image at `client/public/avatar.jpg`, 100x100px rounded

### KaTeX Inline Formula Rendering
- `processInlineKatex()` processes `$...$` (inline) and `$$...$$` (block) after marked
- Uses `katex.renderToString()` to render to HTML
- **Pitfall**: literal `\n` inside `[^$\n]` regex breaks Vue SFC/Babel parser — use `[^$]+` instead

### Mermaid Diagrams
- Rendered as images via mermaid.ink service (in marked code renderer)
- `injectMermaid()` is now empty (no client-side JS rendering needed)

### Pin Articles
- Admin editor has "📌 Pin Article" checkbox
- `editorForm.is_pinned` field
- Pinned articles show 📌 at top of all lists

### Reading Time
- `readingTime(content)`: Chinese chars/400 + English words/200
- Chinese regex: `[\u4e00-\u9fff\u3400-\u4dbf\uf900-\ufaff]`
- Shows "X min" on article lists and article detail page

### Timeline / Homepage Recent Updates
- Replaced old "Recent Updates" section with year-month grouped article timeline
- Each article shows category dot + category label + tags + date
- Standalone Timeline.vue page exists but removed from nav (already on homepage)

### Image Click-to-Zoom + Lazy Loading
- Article.vue: `injectImageZoom()` adds `loading="lazy"` to all `<img>`
- Click to open fullscreen overlay, click again to close
- Skips KaTeX formulas and Mermaid diagrams

### Search Enhancement
- Frontend: `articleAPI.search(q)` → `GET /api/search?q=xxx`
- Backend: `searchArticles` handler (searches title/content/summary/tags)

### Admin Article Management Enhancements
- Search box (filters by title/content/summary in real-time)
- Category dropdown filter + Tag dropdown filter (tags auto-collected from all articles)
- Article list shows category label + tags, clicking a tag filters by it
- Shows filtered result count

### Study Notes Independent Page
- Removed "Study Notes" from Blog.vue sidebar (now 4 categories)
- Study notes accessed via `/study` with own subcategory sidebar

### Global Typography Polish
- Body line-height: 1.7 → 1.82, Card padding: 22 → 24px
- Article content line-height: 1.85 → 1.92, paragraph spacing: 1 → 1.3em
- Heading margins increased, list/blockquote spacing increased
- Mobile breakpoints preserved

### HTTPS + Custom 404
- Alibaba Cloud free SSL certificate (`.pem` + `.key`), stored in `ssl/` folder
- Nginx config: HTTP → HTTPS 301 redirect, TLSv1.2/1.3
- `client/public/404.html` — custom 404 page with navigation buttons

### Category Count Bug Fix
- `allColumns` changed from raw array to `ref([...])`, full replacement triggers reactivity
- `fetchColumnCounts` uses `Promise.allSettled` for resilience
- `watch(route.fullPath)` always calls `fetchColumnCounts()` unconditionally

### SEO Optimization
- index.html title: "Zephyr - Zheng Zhiyi's Personal Site"
- Meta description, keywords, Open Graph tags
- `lang="zh-CN"`

---

## CRITICAL RULES (read before any work)

1. **NEVER use `git checkout <file>` to restore a single file** — it overwrites uncommitted features. Use `git show HEAD:file > file` instead.
2. **NEVER chain Python/SSH in PowerShell CLI** — `$` and quotes get mangled by multiple escaping layers. Write a `.py` file and run it.
3. **Always `rm -rf` target dir before `cp`** when deploying frontend — `cp -r` merges, old JS files coexist and break chunk loading.
4. **Always `systemctl stop` before overwriting running Go binary** — `cp` on a running binary fails with `Text file busy`.
5. **Keep API response format consistent** — all endpoints return `{success: true, data: {...}}`, frontend reads `res.data.xxx`.
6. **Promise.all and cats array must be in exact same order** — mismatch causes category swaps.
7. **Vue ref-wrapped plain arrays: modifying object properties may not trigger update** — use full replacement `arr.value = [...newArr]`.
8. **PowerShell here-strings expand `$` as variables** — use Python to write files containing JavaScript with `$`.
9. **Python raw strings don't interpret `\U0001f4dd` as emoji** — write the actual emoji character directly.
10. **JS regex `[^$\n]` with literal newline breaks Vue SFC/Babel parsing** — use `[^$]+` instead.
11. **Let's Encrypt on Alibaba Cloud mainland servers gets blocked by WAF** — use Alibaba Cloud free SSL cert instead.
12. **SPA + `try_files` makes 404 page tricky** — use `error_page 404 /404.html` in Nginx config.

## Git History (2026-08-05)
```
4bdff92 feat: HTTPS certificate + custom 404 page
6404eec fix: reading time Chinese detection; homepage timeline replaces recent updates
0c43f41 fix: homepage full rewrite - timeline replaces recent updates
aaf6cb3 style: global typography improvements
dda6daf feat: pin articles, reading time, image zoom+lazy, timeline, SEO, search API
a291bfc fix: allColumns reactive ref + full array replacement
b775bd3 fix: always refresh counts; study notes removed from article sidebar
f324b17 fix: category count fix; Admin search/category/tag filters
52c697d fix: tag filtering, avatar, KaTeX inline, SESSION_SUMMARY garbled
47cc19a docs: update SESSION_SUMMARY - all 2026-08-05 features + pitfalls
```
