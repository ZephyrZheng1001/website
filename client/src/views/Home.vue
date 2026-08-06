<template>
  <div style="padding-top:48px;padding-bottom:32px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">

    <section style="margin-bottom:40px;">
      <h1 class="hero-title">
        <span class="typed-text">{{ typedText }}</span>
      </h1>
      <p style="font-size:1.08rem;color:var(--text-muted);max-width:560px;line-height:1.8;">
        欢迎来到我的主页，我是郑智毅，就读于同济大学计算机专业。在这里你可以了解我所学的知识，也可以<a href="/#/about" style="color:var(--accent);text-decoration:underline;">了解我</a>
      </p>
    </section>

    <section style="display:flex;gap:32px;padding:14px 0;border-top:1px solid var(--border);border-bottom:1px solid var(--border);margin-bottom:24px;font-size:0.84rem;color:var(--text-muted);">
      <div><strong style="color:var(--accent);">{{ stats.articles }}</strong> 篇文章</div>
      <div><strong style="color:var(--accent);">{{ stats.visitors }}</strong> 位访客</div>
    </section>

    <section style="margin-bottom:28px;">
      <div style="background:var(--bg-card);border:1px solid var(--border);border-radius:10px;padding:14px 20px;box-shadow:var(--shadow);">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:6px;">
          <span style="font-weight:600;font-size:0.78rem;color:var(--text);">{{ currentYear }} 年度进度</span>
          <span style="font-size:0.7rem;color:var(--text-muted);">{{ daysElapsed }}/{{ totalDays }} 天 · {{ yearPercent }}% · {{ currentTime }}</span>
        </div>
        <div style="height:5px;background:var(--bg-hover);border-radius:5px;overflow:hidden;margin-bottom:10px;">
          <div :style="{width:yearPercent+'%'}" style="height:100%;border-radius:5px;background:linear-gradient(90deg,#2d8a7b,#5cc4ae);transition:width 1s ease;"></div>
        </div>
        <div style="display:flex;align-items:center;justify-content:space-between;">
          <a v-if="lc.totalSolved" href="https://leetcode.cn/u/amazing-joliotj1h" target="_blank" style="font-weight:600;font-size:0.78rem;color:var(--text);text-decoration:none;">
            LeetCode
            <span style="font-weight:400;color:var(--text-muted);font-size:0.72rem;">{{ lc.totalSolved }}/{{ lc.totalQuestions }} ↗</span>
          </a>
          <span v-else class="lc-skeleton" style="width:120px;height:14px;border-radius:3px;display:inline-block;"></span>
          <span v-if="lc.totalSolved" style="font-size:0.7rem;color:var(--text-muted);font-weight:500;">{{ Math.round(lc.totalSolved / lc.totalQuestions * 100) }}%</span>
        </div>
        <div v-if="lc.totalSolved" style="height:5px;background:var(--bg-hover);border-radius:5px;overflow:hidden;margin-top:6px;">
          <div :style="{width: lc.totalSolved / lc.totalQuestions * 100 + '%'}" style="height:100%;border-radius:5px;background:linear-gradient(90deg,#e65100,#ff9800);transition:width 0.6s ease;"></div>
        </div>
      </div>
    </section>

    <section v-if="timelineItems.length" style="margin-top:44px;">
      <h3 style="font-size:1rem;color:var(--text-muted);font-weight:500;margin-bottom:20px;">最近更新</h3>
      <div class="home-timeline">
        <div v-for="item in timelineItems" :key="item.key" class="home-tl-group">
          <div class="home-tl-marker">
            <span class="home-tl-year">{{ item.year }}</span>
            <span class="home-tl-month">{{ item.month }}月</span>
          </div>
          <div class="home-tl-cards">
            <router-link
              v-for="a in item.articles"
              :key="a.id"
              :to="articleLink(a)"
              class="card home-tl-card"
            >
              <div style="display:flex;align-items:center;gap:8px;">
                <span class="home-tl-dot" :class="dotClass(a.category)"></span>
                <strong class="home-tl-title">{{ a.is_pinned ? '\U0001f4cc ' : '' }}{{ a.title }}</strong>
                <div style="display:flex;align-items:center;gap:6px;flex-shrink:0;margin-left:auto;">
                  <span class="tag" style="font-size:0.68rem;padding:1px 6px;">{{ catLabel(a.category) }}</span>
                  <span style="font-size:0.7rem;color:var(--text-muted);white-space:nowrap;">{{ dayLabel(a.created_at) }}</span>
                </div>
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { articleAPI } from '../api'

const timelineItems = ref([])
const stats = ref({ articles: 0, visitors: 0 })
const lc = ref({ totalSolved: 0, totalQuestions: 0, ranking: 0 })
const todayDate = ref('')

const yp = calcYearProgress()
const currentYear = ref(yp.year)
const daysElapsed = ref(yp.daysElapsed)
const totalDays = ref(yp.totalDays)
const yearPercent = ref(yp.percent)
const currentTime = ref('')

const fullText = '你好，我是Zephyr！'
const typedText = ref('')
let typeTimer = null
function startTyping() {
  let i = 0
  typedText.value = ''
  typeTimer = setInterval(() => {
    if (i <= fullText.length) {
      typedText.value = fullText.slice(0, i)
      i++
    } else {
      clearInterval(typeTimer)
    }
  }, 120)
}
startTyping()

function catLabel(c) { const m = { blog: '技术文章', leetcode: '算法笔记', projects: '项目', study: '学习笔记', notes: '碎碎念' }; return m[c] || c }
function dotClass(c) { const m = { blog: 'dot-blog', leetcode: 'dot-leetcode', projects: 'dot-projects', study: 'dot-study', notes: 'dot-notes' }; return m[c] || 'dot-blog' }
function articleLink(a) { const cat = a.category || 'blog'; return cat === 'blog' ? '/blog/' + a.id : '/' + cat + '/' + a.id }
function dayLabel(d) { return d ? new Date(d).toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' }) : '' }
function formatNow() { const n = new Date(); return n.getFullYear() + '-' + String(n.getMonth()+1).padStart(2,'0') + '-' + String(n.getDate()).padStart(2,'0') + ' ' + String(n.getHours()).padStart(2,'0') + ':' + String(n.getMinutes()).padStart(2,'0') + ':' + String(n.getSeconds()).padStart(2,'0') }

function calcYearProgress() {
  const now = new Date()
  const year = now.getFullYear()
  const start = new Date(year, 0, 1)
  const end = new Date(year + 1, 0, 1)
  const elapsed = Math.floor((now - start) / (1000 * 60 * 60 * 24))
  const total = Math.floor((end - start) / (1000 * 60 * 60 * 24))
  return { year, daysElapsed: elapsed, totalDays: total, percent: Math.round((elapsed / total) * 100 * 10) / 10 }
}

function formatDate(d) { return d ? new Date(d).toLocaleDateString('zh-CN') : '' }

onMounted(async () => {
  currentTime.value = formatNow()
  setInterval(() => { currentTime.value = formatNow() }, 1000)
  try {
    const res = await articleAPI.list({ limit: 200 })
    const all = res.data.articles || []
    stats.value.articles = res.data.total || all.length
    try {
      const uvRes = await fetch('/api/visitors/count')
      const uvData = await uvRes.json()
      stats.value.visitors = uvData.data?.visitors || uvData.visitors || 0
    } catch(e) {}
    try {
      const lcRes = await fetch('/api/leetcode/stats')
      const lcData = await lcRes.json()
      if (lcData.success) lc.value = lcData.data
    } catch(e) {}
    const d = new Date()
    todayDate.value = d.getFullYear() + '-' + String(d.getMonth()+1).padStart(2,'0') + '-' + String(d.getDate()).padStart(2,'0')

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
.hero-title {
  font-size: 2.4rem;
  font-weight: 800;
  letter-spacing: -0.5px;
  line-height: 1.3;
  margin-bottom: 18px;
}
.typed-cursor {
  color: var(--accent);
  animation: blink 1s step-end infinite;
}
@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
.lc-skeleton {
  background: var(--bg-hover);
  border-radius: 3px;
  animation: lc-shimmer 1.5s ease-in-out infinite;
}
@keyframes lc-shimmer {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}
.home-timeline { position: relative; padding-left: 16px; border-left: 2px solid var(--border); }
.home-tl-group { margin-bottom: 16px; }
.home-tl-marker { display: flex; align-items: baseline; gap: 8px; margin-bottom: 6px; margin-left: -26px; }
.home-tl-year { font-size: 1rem; font-weight: 700; color: var(--accent); }
.home-tl-month { font-size: 0.8rem; color: var(--text-muted); font-weight: 500; }
.home-tl-cards { display: flex; flex-direction: column; gap: 4px; }
.home-tl-card { padding: 10px 14px !important; display: block; }
.home-tl-card:hover { text-decoration: none; }
.home-tl-title { font-size: 0.88rem; font-weight: 600; color: var(--text); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; min-width: 0; display: block; }
.home-tl-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.dot-blog { background: #2d8a7b; }
.dot-leetcode { background: #e65100; }
.dot-projects { background: #6c3fb5; }
.dot-study { background: #2d7dd2; }
.dot-notes { background: #d4a574; }
@media (max-width: 768px) {
  .hero-title { font-size: 1.8rem; }
  .home-timeline { padding-left: 10px; }
  .home-tl-group { margin-bottom: 12px; }
  .home-tl-marker { margin-left: -20px; }
  .home-tl-year { font-size: 0.9rem; }
  .home-tl-month { font-size: 0.72rem; }
  .home-tl-card { padding: 8px 10px !important; }
  .home-tl-title { font-size: 0.82rem; }
}
</style>