<template>
  <div style="padding-top:48px;padding-bottom:32px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">

    <!-- Hero -->
    <section style="margin-bottom:56px;">
      <h1 style="font-size:2.4rem;font-weight:800;letter-spacing:-0.5px;line-height:1.3;margin-bottom:18px;">
        你好，我是Zephyr！
      </h1>
      <p style="font-size:1.08rem;color:var(--text-muted);max-width:560px;line-height:1.8;">
        欢迎来到我的主页，我是郑智毅，就读于同济大学计算机专业。在这里你可以了解我所学的知识，也可以<a href="/#/about" style="color:var(--accent);text-decoration:underline;">了解我</a>
      </p>
    </section>

    <!-- Stats line -->
    <section style="display:flex;gap:48px;padding:18px 0;border-top:1px solid var(--border);border-bottom:1px solid var(--border);margin-bottom:44px;font-size:0.88rem;color:var(--text-muted);">
      <div><strong style="color:var(--accent);">{{ stats.articles }}</strong> 篇文章</div>
      <div><strong style="color:var(--accent);">Always</strong> 在线</div>
    </section>

    <!-- 年度进度 -->
    <section style="margin-bottom:40px;">
      <div style="background:var(--bg-card);border:1px solid var(--border);border-radius:12px;padding:22px 28px;box-shadow:var(--shadow);">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:10px;">
          <span style="font-weight:600;font-size:0.95rem;color:var(--text);">{{ currentYear }} 年度进度 · 当前时间 {{ currentTime }}</span>
          <span style="font-size:0.9rem;font-weight:700;color:var(--accent);">{{ yearPercent }}%</span>
        </div>
        <div style="position:relative;height:10px;background:var(--bg-hover);border-radius:10px;overflow:hidden;">
          <div :style="{width:yearPercent+'%'}" style="height:100%;border-radius:10px;background:linear-gradient(90deg,#2d8a7b,#5cc4ae);transition:width 1s ease;"></div>
        </div>
        <div style="display:flex;justify-content:space-between;margin-top:8px;font-size:0.78rem;color:var(--text-muted);">
          <span>1月1日</span>
          <span>已过 {{ daysElapsed }} 天 / 全年 {{ totalDays }} 天</span>
          <span>12月31日</span>
        </div>
      </div>
    </section>

    <!-- Navigation cards -->
    <section>
      <div class="nav-grid">
        <router-link v-for="col in columns" :key="col.key" :to="col.link" class="card nav-card">
          <div class="nav-card-icon">{{ col.icon }}</div>
          <div>
            <h3>{{ col.title }}</h3>
            <p>{{ col.desc }}</p>
          </div>
          <span class="nav-card-arrow">→</span>
        </router-link>
      </div>
    </section>

    <!-- Latest from each column -->
    <section v-if="latestList.length" style="margin-top:52px;">
      <h3 style="font-size:1rem;color:var(--text-muted);font-weight:500;margin-bottom:20px;">最近更新</h3>
      <div style="display:flex;flex-direction:column;gap:12px;">
        <router-link
          v-for="item in latestList"
          :key="item.id"
          :to="item.link"
          class="card"
          style="display:flex;align-items:center;gap:14px;padding:16px 20px;"
        >
          <span style="font-size:1.2rem;">{{ item.icon }}</span>
          <div style="flex:1;min-width:0;">
            <div style="font-size:0.92rem;font-weight:600;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ item.title }}</div>
            <div style="font-size:0.78rem;color:var(--text-muted);">{{ item.column }} · {{ formatDate(item.created_at) }}</div>
          </div>
          <span style="color:var(--text-muted);font-size:0.85rem;">→</span>
        </router-link>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { articleAPI } from '../api'

const columns = [
  { key: 'blog', title: '技术文章', desc: '技术分享与教程', icon: '📝', link: '/blog' },
  { key: 'leetcode', title: '算法笔记', desc: '算法题解笔记', icon: '💡', link: '/leetcode' },
  { key: 'projects', title: '项目', desc: '项目复盘方案', icon: '🚀', link: '/projects' },
  { key: 'notes', title: '碎碎念', desc: '日常随想记录', icon: '💬', link: '/notes' },
]

const timelineItems = ref([])
const stats = ref({ articles: 0 })

const yp = calcYearProgress()
const currentYear = ref(yp.year)
const daysElapsed = ref(yp.daysElapsed)
const totalDays = ref(yp.totalDays)
const yearPercent = ref(yp.percent)
const currentTime = ref('')

function catLabel(c) {
  const m = { blog: '技术文章', leetcode: '算法笔记', projects: '项目', study: '学习笔记', notes: '碎碎念' }
  return m[c] || c
}
function dotClass(c) {
  const m = { blog: 'dot-blog', leetcode: 'dot-leetcode', projects: 'dot-projects', study: 'dot-study', notes: 'dot-notes' }
  return m[c] || 'dot-blog'
}
function articleLink(a) {
  const cat = a.category || 'blog'
  return cat === 'blog' ? '/blog/' + a.id : '/' + cat + '/' + a.id
}
function dayLabel(d) {
  return d ? new Date(d).toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' }) : ''
}

function formatNow() { const n = new Date(); return n.getFullYear() + '-' + String(n.getMonth()+1).padStart(2,'0') + '-' + String(n.getDate()).padStart(2,'0') + ' ' + String(n.getHours()).padStart(2,'0') + ':' + String(n.getMinutes()).padStart(2,'0') + ':' + String(n.getSeconds()).padStart(2,'0') }
function calcYearProgress() {
  const now = new Date()
  const year = now.getFullYear()
  const start = new Date(year, 0, 1)
  const end = new Date(year + 1, 0, 1)
  const elapsed = Math.floor((now - start) / (1000 * 60 * 60 * 24))
  const total = Math.floor((end - start) / (1000 * 60 * 60 * 24))
  return {
    year,
    daysElapsed: elapsed,
    totalDays: total,
    percent: Math.round((elapsed / total) * 100 * 10) / 10
  }
}

function formatDate(d) {
  return d ? new Date(d).toLocaleDateString('zh-CN') : ''
}

onMounted(async () => {
  currentTime.value = formatNow()
  setInterval(() => { currentTime.value = formatNow() }, 1000)
  try {
    const res = await articleAPI.list({ limit: 200 })
    const all = res.data.articles || []
    stats.value.articles = res.data.total || all.length
    // Group by year-month
    const groups = {}
    all.forEach(a => {
      const d = new Date(a.created_at)
      const key = d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0')
      if (!groups[key]) groups[key] = { year: d.getFullYear(), month: d.getMonth() + 1, key, articles: [] }
      groups[key].articles.push(a)
    })
    timelineItems.value = Object.values(groups).sort((a, b) => {
      if (a.year !== b.year) return b.year - a.year
      return b.month - a.month
    })
  } catch (e) { console.error(e) }
})
</script>

<style scoped>
.nav-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}
.nav-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 18px;
  cursor: pointer;
}
.nav-card:hover { text-decoration: none; }
.nav-card-icon {
  font-size: 1.3rem;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent-light);
  border-radius: 10px;
  flex-shrink: 0;
}
.nav-card h3 {
  font-size: 0.9rem;
  color: var(--text);
  margin-bottom: 2px;
}
.nav-card p {
  font-size: 0.78rem;
  color: var(--text-muted);
}
.nav-card-arrow {
  margin-left: auto;
  color: var(--text-muted);
  font-size: 0.9rem;
  transition: transform 0.2s;
}
.nav-card:hover .nav-card-arrow { transform: translateX(3px); color: var(--accent); }

@media (max-width: 768px) {
  .nav-grid { grid-template-columns: 1fr; }
}

.home-timeline {
  position: relative;
  padding-left: 16px;
  border-left: 2px solid var(--border);
}
.home-tl-group { margin-bottom: 20px; }
.home-tl-marker {
  display: flex; align-items: baseline; gap: 8px;
  margin-bottom: 8px; margin-left: -26px;
}
.home-tl-year { font-size: 1rem; font-weight: 700; color: var(--accent); }
.home-tl-month { font-size: 0.8rem; color: var(--text-muted); font-weight: 500; }
.home-tl-cards { display: flex; flex-direction: column; gap: 6px; }
.home-tl-card { padding: 12px 16px !important; display: block; }
.home-tl-card:hover { text-decoration: none; transform: translateY(-1px); }
.home-tl-dot {
  width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; margin-top: 5px;
}
.dot-blog { background: #2d8a7b; }
.dot-leetcode { background: #e65100; }
.dot-projects { background: #6c3fb5; }
.dot-study { background: #2d7dd2; }
.dot-notes { background: #d4a574; }
@media (max-width: 768px) {
  .home-timeline { padding-left: 12px; }
  .home-tl-marker { margin-left: -22px; }
}
</style>
