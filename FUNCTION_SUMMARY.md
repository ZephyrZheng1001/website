# Zephyr Website — Complete Feature Summary (2026-08-05)

## Tech Stack
| Layer | Technology |
|-------|-----------|
| Frontend | Vue 3, Vite, Vue Router (hash mode), Pinia |
| Styling | CSS custom properties (light/dark themes), scoped styles |
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
| `/timeline` | Timeline.vue | Full timeline (route exists, nav removed — shown on homepage) |
| `/search` | Search.vue | Full-text search |
| `/about` | About.vue | About page + resume download |
| `/admin` | Admin.vue | Admin panel (auth required) |
| `*` | 404.html | Custom 404 page |

### Layout (App.vue)
- **Top bar**: Logo, GitHub icon, QQ icon, search icon, dark mode toggle
- **Desktop sidebar**: Home, Articles (4 categories), Study Notes, Timeline, About
- **Mobile**: Hamburger menu + bottom navigation bar
- **Reading progress bar**: 3px gradient line at top, fills as user scrolls
- **Back-to-top button**: Appears after scrolling 400px, smooth scroll to top

---

## 2. Homepage Features

### Hero Section
- "Hello, I'm Zephyr" intro with link to About page

### Stats Line
- `7 articles · X visitors` (visitor count from MySQL, per IP+UserAgent)

### Annual Progress Bar
- Shows current year progress by days (e.g. "2026 is 60.3% complete")
- Real-time clock display
- Gradient progress bar

### Navigation Cards
- 4 cards: Tech Articles, Algorithm Notes, Projects, Random Thoughts
- Each with emoji icon, title, description, arrow hover animation

### Timeline (Recent Updates)
- All articles grouped by year-month, displayed chronologically
- Each article shows: category color dot, pin indicator (if pinned), title, category label, tags (up to 3), date
- Built from a single API call (`GET /api/articles?limit=200`)

---

## 3. Article Features

### Article Detail Page
- **Reading time**: Chinese chars/400 + English words/200
- **View count**: PV counter incremented on each visit (MySQL `view_count` column)
- **Table of Contents**: Auto-extracted from h2/h3 headings, sticky sidebar, active heading tracking on scroll
- **Prev/Next navigation**: Links to previous/next article in same category

### Markdown Rendering
- **KaTeX**: Inline `$...$` and block `$$...$$` math formulas rendered client-side
- **Code blocks**: Syntax highlighting via highlight.js with language tag label (supports: JavaScript, Python, Go, Bash, CSS, JSON, XML/HTML, SQL, YAML, Java, C/C++, MySQL — unknown languages fall back to highlightAuto)
- **Copy button**: Each code block has a "copy" button (copies to clipboard)
- **Tables**: Styled with header rows, alternating row hover
- **Blockquotes**: Styled with left border and background

### Image Features
- **Lazy loading**: All images get `loading="lazy"` attribute
- **Click-to-zoom**: Click any image for fullscreen overlay (black background), click anywhere to close
- **Exclusions**: KaTeX formulas and Mermaid diagrams are excluded from zoom behavior

### Mermaid Diagrams
- Currently rendered as plain code blocks with syntax highlighting
- (Client-side mermaid.js rendering tested but removed due to encoding issues with mermaid.ink; user uploads diagrams as images instead)

---

## 4. Tag & Category System

### Sidebar Tag Filtering (Blog.vue)
- Tags collected from ALL articles in current category (not just filtered results)
- Click a tag → URL becomes `?tag=xxx` → API filters articles
- Active tag highlighted with accent color
- "Clear filter" button appears when a tag is active
- Clicking tags inside article cards also navigates to filter

### Category Sidebar
- Shows 4 categories with article counts
- Counts fetched via `fetchColumnCounts()` using `Promise.allSettled`
- `allColumns` is a Vue `ref` — counts updated by full array replacement (triggers reactivity)
- Watches `route.fullPath` to refresh counts on any navigation

---

## 5. Study Notes (Independent Page)

- Separate from article categories, accessed via `/study`
- Subcategory sidebar: dynamic categories from `study_categories` MySQL table
- Status filters: All / Learning / Todo / Done (with color-coded dots)
- Overall progress bar showing completed/total percentage
- Each article card shows: status dot, pin indicator, title, subcategory badge, tags, date

---

## 6. Admin Panel

### Authentication
- Bearer token JWT auth
- Default admin: `admin` / `admin123`
- Token stored in localStorage
- Password change support (bcrypt hashed)

### Article Management
- **Search**: Real-time filter by title/content/summary
- **Category filter**: Dropdown selector
- **Tag filter**: Dropdown auto-populated from all articles, clickable tags in list
- **List view**: Title, pin indicator, category badge, tags, date, edit/delete buttons
- **CRUD**: Create, edit (pre-filled form), delete (with confirmation)

### Article Editor
- Title, category, subcategory (for study notes), summary, tags, content
- **Pin toggle**: Checkbox "📌 Pin Article" (pinned articles appear first in lists)
- **Live preview**: Side-by-side Markdown editor & rendered preview
- **Study status**: Learning / Todo / Done (for study notes category)

### Study Category Management
- Add/delete/reorder subcategories
- Each with name, icon, sort order

---

## 7. Search

- Dedicated `/search` page
- Backend `GET /api/search?q=xxx` searches title, content, summary, tags
- Results show: category icon, title, summary excerpt, category badge, tags, date
- Results link to correct article path based on category

---

## 8. Visitor & Reading Statistics

### Server-Side Visitor Tracking
- MySQL table `site_visitors` records IP + UserAgent + date on each visit
- `POST /api/visitor` silently records visits
- `GET /api/visitors/count` returns `COUNT(DISTINCT CONCAT(ip, user_agent))`
- Displayed on homepage stats line: "X visitors"

### Article View Count
- MySQL `view_count` column, incremented on each `getArticle` call
- Displayed on article detail page: "X reads"

### Site Stats Footer
- Bottom of every page: "X articles · ~Y words · Z reads · Admin"
- Aggregated from all articles via `fetchSiteStats()`

---

## 9. SEO & Meta

- `<title>Zephyr - Zheng Zhiyi's Personal Site</title>`
- Meta description, keywords, author
- Open Graph tags (og:title, og:description, og:type, og:url)
- `lang="zh-CN"` on `<html>`
- `sitemap.xml` and `robots.txt` in public directory

---

## 10. About Page

- Bio section: intro text
- Contact cards: school email, GitHub, QQ, WeChat
- **Resume download**: Link to `/resume.pdf` (uploaded PDF)
- Site info: "Powered by Go · Deployed on my own server"

---

## 11. Dark Mode

- System-following dark/light theme via CSS custom properties
- Toggle button in top bar (sun/moon emoji)
- Persisted via composable (useTheme)
- All components support `[data-theme="dark"]` overrides
- Article content: code blocks, tables, blockquotes, difficulty tags all dark-adapted

---

## 12. Typography & Design

- Max content width: 720px centered
- Body line-height: 1.82, article line-height: 1.92
- Card padding: 24px, border-radius: 10px
- Card hover: subtle shadow elevation + translateY(-2px)
- Color palette: green accent (#2d8a7b), warm secondary (#d4a574)
- Mobile responsive: hamburger menu, bottom nav bar, collapsed sidebars

---

## 13. HTTPS & Infrastructure

- SSL certificate: Alibaba Cloud free DV certificate
- Nginx config: HTTP → HTTPS 301 redirect, custom 404 page
- systemd service: `zephyr-api` auto-starts on boot
- Backend binary: `/opt/zephyr/server/server_linux`
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
cd A:\ZEPHYR_S\website\client
npm run build
ssh admin@47.116.136.145 "rm -rf ~/dist && mkdir ~/dist"
scp -r dist/* admin@47.116.136.145:~/dist/
ssh admin@47.116.136.145 "sudo rm -rf /var/www/zephyrzheng/* && sudo cp -r ~/dist/* /var/www/zephyrzheng/"

# Backend
cd A:\ZEPHYR_S\website\server
$env:GOOS='linux'; $env:GOARCH='amd64'; go build -o server_linux main.go
scp server_linux admin@47.116.136.145:~/server_linux
ssh admin@47.116.136.145 "sudo systemctl stop zephyr-api && sudo cp ~/server_linux /opt/zephyr/server/server_linux && sudo chmod +x /opt/zephyr/server/server_linux && sudo systemctl start zephyr-api"
```

**Critical**: systemd runs `/opt/zephyr/server/server_linux`, NOT `server`. Overwrite `server_linux`, not `server`.
