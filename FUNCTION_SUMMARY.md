# Zephyr Website — Complete Feature Summary (2026-08-05)

## Tech Stack
| Layer | Technology |
|-------|-----------|
| Frontend | Vue 3, Vite, Vue Router (hash mode), Pinia |
| Styling | CSS custom properties (light/dark themes), scoped styles, smooth transitions |
| Backend | Go (stdlib net/http), MySQL |
| Server | Ubuntu 24.04, Nginx (HTTPS with Alibaba Cloud SSL), systemd |
| Markdown | marked.js + KaTeX (inline & block math) + highlight.js |
| Domain | zephyrzheng.cn (ICP filing pending, accessible now) |

---

## 1. Site Architecture

### Routing (Hash Mode)
| Path | Component | Purpose |
|------|-----------|---------|
| `/` | Home.vue | Landing page with timeline |
| `/blog` | Blog.vue | Tech articles |
| `/leetcode` | Blog.vue | Algorithm notes |
| `/projects` | Blog.vue | Project writeups |
| `/notes` | Blog.vue | Random thoughts |
| `/study` | Study.vue | Structured learning notes (independent page) |
| `/blog/:id` `/leetcode/:id` etc. | Article.vue | Article detail |
| `/timeline` | Timeline.vue | Full timeline (route exists, shown on homepage) |
| `/search` | Search.vue | Full-text search with highlight |
| `/about` | About.vue | About page + resume download |
| `/admin` | Admin.vue | Admin panel (auth required) |
| `*` | 404.html | Custom 404 page |

### Layout (App.vue)
- **Top bar**: Logo, GitHub icon, QQ icon, search icon, dark mode toggle
- **Desktop sidebar**: Home, Articles (4 categories), Study Notes, Timeline, About
- **Mobile**: Hamburger menu + bottom navigation bar
- **Reading progress bar**: 3px gradient line at top, fills as user scrolls
- **Back-to-top button**: Appears after scrolling 400px, smooth scroll to top
- **Site footer**: Article count, word count, total reads, admin link, copyright

---

## 2. Homepage Features

### Hero Section
- "Hello, I'm Zephyr" intro with link to About page

### Stats Line
- `X articles · Y visitors` (visitor count from MySQL, per IP+UserAgent dedup)

### Annual Progress Bar
- Shows current year progress by days with gradient progress bar
- Real-time clock display

### Navigation Cards
- 4 cards: Tech Articles, Algorithm Notes, Projects, Random Thoughts
- Each with emoji icon, title, description, arrow hover animation

### Timeline (Recent Updates)
- All articles grouped by year-month, displayed chronologically
- Each article: category color dot, pin indicator, title, category label, tags (up to 3), date

---

## 3. Article Features

### Article Detail Page
- **Reading time**: Chinese chars/400 + English words/200 (list pages use `content` field)
- **View count**: PV counter incremented on each visit (MySQL `view_count` column)
- ~~**Reading position memory**~~: Removed (2026-08-05) ? articles now always start from top
- **Table of Contents**: Auto-extracted from h2/h3 headings, sticky sidebar, active heading tracking
- **Prev/Next navigation**: Links to previous/next article in same category

### Markdown Rendering
- **KaTeX**: Inline `$...$` and block `$$...$$` math formulas rendered client-side
- **Code blocks**: Syntax highlighting via highlight.js with language tag label
- **Copy button**: Each code block has a "copy" button
- **Tables, blockquotes**: Styled with proper spacing

### Image Features
- **Lazy loading**: All images get `loading="lazy"` attribute
- **Click-to-zoom**: Click any image for fullscreen overlay, click to close
- **Exclusions**: KaTeX formulas excluded from zoom

### Mermaid Diagrams
- Rendered as plain code blocks with syntax highlighting (user uploads images instead)

---

## 4. Tag & Category System

### Sidebar Tag Filtering (Blog.vue)
- Tags collected from ALL articles in current category
- Click a tag -> URL becomes `?tag=xxx` -> API filters articles
- Active tag highlighted, "clear filter" button when active
- Clicking tags inside article cards also navigates to filter

### Category Sidebar
- 4 categories with article counts
- `allColumns` is Vue `ref` — full array replacement for reactivity
- Watches `route.fullPath` to refresh counts on any navigation

---

## 5. Study Notes (Independent Page)

- Accessed via `/study`
- Dynamic subcategory sidebar from MySQL `study_categories`
- Status filters: All / Learning / Todo / Done (color-coded dots)
- Overall progress bar showing completed/total percentage

---

## 6. Admin Panel

### Authentication
- Bearer token JWT auth, default admin: `admin` / `admin123`
- Password change support (bcrypt hashed)

### Article Management
- **Search**: Real-time filter by title/content/summary
- **Category filter**: Dropdown selector
- **Tag filter**: Dropdown auto-populated from all articles, clickable tags
- **List view**: Title, pin indicator, category badge, tags, date, edit/delete buttons
- **CRUD**: Create, edit (pre-filled form), delete (with confirmation)

### Article Editor
- Title, category, subcategory, summary, tags, content
- **Pin toggle**: "📌 Pin Article" checkbox
- **Live preview**: Side-by-side Markdown editor & rendered preview
- **Study status**: Learning / Todo / Done (for study notes)

---

## 7. Search

- Dedicated `/search` page
- **Ctrl+K / Cmd+K** global shortcut to open search
- Backend `GET /api/search?q=xxx` searches title, content, summary, tags
- **Keyword highlighting**: Matches highlighted with `<mark>` in title and summary
- **Suggestions**: When no results, shows 5 latest articles as recommendations
- Dark mode compatible highlighting

---

## 8. Visitor & Reading Statistics

### Server-Side Visitor Tracking
- MySQL `site_visitors` table: IP + UserAgent + date
- `POST /api/visitor` silently records visits
- `GET /api/visitors/count` returns `COUNT(DISTINCT CONCAT(ip, user_agent))`
- Displayed on homepage stats line: "X visitors"

### Article View Count
- MySQL `view_count` column, incremented on each article view
- Displayed on article detail page: "X reads"

### Site Stats Footer
- "X articles · ~Y words · Z reads · Admin"
- Words calculated from article content + summary

---

## 9. Dark Mode

- **System-following**: `matchMedia('prefers-color-scheme: dark')` on first visit
- **Anti-FOUC**: Inline `<script>` in `<head>` sets `data-theme` before page renders
- **Smooth transitions**: `transition: background-color 0.3s, color 0.3s` on body and cards
- **Color scheme**: `html { color-scheme: light dark }` for native browser controls
- All components support `[data-theme="dark"]` overrides

---

## 10. SEO & Meta

- `<title>Zephyr - Zheng Zhiyi's Personal Site</title>`
- Meta description, keywords, author, Open Graph tags
- `sitemap.xml` and `robots.txt`
- `lang="zh-CN"` on `<html>`

---

## 11. About Page

- Bio section, contact cards (email, GitHub, QQ, WeChat)
- Resume download: `/resume.pdf`
- Site info: "Powered by Go · Deployed on my own server"

---

## 12. Typography & Design

- Max content width: 720px centered
- Body line-height: 1.82, article line-height: 1.92
- Card padding: 24px, border-radius: 10px
- Card hover: subtle shadow elevation
- Green accent (#2d8a7b), warm secondary (#d4a574)
- Mobile responsive: hamburger menu, bottom nav bar, collapsed sidebars

---

## 13. HTTPS & Infrastructure

- SSL: Alibaba Cloud free DV certificate
- Nginx: HTTP -> HTTPS 301 redirect, custom 404 page
- systemd: `zephyr-api` auto-starts on boot, runs `/opt/zephyr/server/server_linux`
- Logs: `/opt/zephyr/logs/api.log` and `api_error.log`

---

## 14. API Endpoints

### Public
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/articles` | List articles (category, tag, page, limit params) |
| GET | `/api/articles/:id` | Get article detail + prev/next |
| GET | `/api/search?q=` | Full-text search |
| GET | `/api/study-categories` | List study subcategories |
| POST | `/api/visitor` | Record visitor |
| GET | `/api/visitors/count` | Get unique visitor count |

### Admin (JWT required)
| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api/admin/login` | Login (returns JWT) |
| GET | `/api/admin/articles` | List all articles |
| POST | `/api/admin/articles` | Create article |
| PUT | `/api/admin/articles/:id` | Update article |
| DELETE | `/api/admin/articles/:id` | Delete article |
| PUT | `/api/admin/password` | Change password |
| POST | `/api/admin/study-categories` | Create study category |
| PUT | `/api/admin/study-categories/:id` | Update study category |
| DELETE | `/api/admin/study-categories/:id` | Delete study category |

---

## 15. Database Schema

### articles
```
id, title, content, summary, tags, category, subcategory,
study_status, is_pinned, view_count, created_at, updated_at
```

### users
```
id, username, password_hash, created_at
```

### study_categories
```
id, name, icon, sort_order, created_at
```

### site_visitors
```
id, ip, user_agent, visit_date, created_at
```

---

## 16. Deploy Quick Reference

```powershell
# Frontend only
cd A:\ZEPHYR_S\website\client && npm run build
ssh admin@47.116.136.145 "rm -rf ~/dist && mkdir ~/dist"
scp -r dist/* admin@47.116.136.145:~/dist/
ssh admin@47.116.136.145 "sudo rm -rf /var/www/zephyrzheng/* && sudo cp -r ~/dist/* /var/www/zephyrzheng/"

# Backend (CRITICAL: overwrite server_linux, NOT server)
cd A:\ZEPHYR_S\website\server
$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o server_linux main.go
scp server_linux admin@47.116.136.145:~/server_linux
ssh admin@47.116.136.145 "sudo systemctl stop zephyr-api && sudo cp ~/server_linux /opt/zephyr/server/server_linux && sudo chmod +x /opt/zephyr/server/server_linux && sudo systemctl start zephyr-api"
```


---

## Recent Changes

### 2026-08-05: Remove reading position memory
- Removed `saveReadPos()` and all localStorage scroll-position save/restore logic in Article.vue
- Articles now always start from top; no more auto-scrolling to last position
- Cleaned 15 temp/dev scripts from repo root
- Fixed `inject_word.py` to deduplicate `__DAILY_WORD__` blocks in index.html (was accumulating 6 blocks)
- Cleaned server: removed old `server` binary, `go.mod`, `go.sum`, `main.go`, `app.log` from `/opt/zephyr/server/`
