<template>
  <div style="padding-top:40px;padding-bottom:80px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">
    <h1 style="font-size:1.8rem;margin-bottom:8px;">🕐 时光轴</h1>
    <p style="color:var(--text-muted);margin-bottom:36px;">所有文章按时间排列</p>

    <div v-if="loading" style="text-align:center;color:var(--text-muted);padding:60px;">加载中...</div>

    <div v-else-if="timelineItems.length === 0" style="text-align:center;color:var(--text-muted);padding:60px;">
      还没有文章
    </div>

    <div v-else class="timeline">
      <div v-for="item in timelineItems" :key="item.year + item.month" class="timeline-group">
        <div class="timeline-marker">
          <span class="timeline-year">{{ item.year }}</span>
          <span class="timeline-month">{{ item.month }}月</span>
          <span class="timeline-count">{{ item.articles.length }}篇</span>
        </div>
        <div class="timeline-cards">
          <router-link
            v-for="a in item.articles"
            :key="a.id"
            :to="articleLink(a)"
            class="card timeline-card"
          >
            <div style="display:flex;align-items:flex-start;gap:12px;">
              <span class="timeline-dot" :class="catDotClass(a.category)"></span>
              <div style="flex:1;min-width:0;">
                <div style="display:flex;align-items:center;gap:8px;margin-bottom:4px;">
                  <span class="tag" style="font-size:0.7rem;padding:1px 7px;">{{ catLabel(a.category) }}</span>
                  <span v-if="a.is_pinned" style="font-size:0.7rem;">📌</span>
                </div>
                <strong style="font-size:0.95rem;display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ a.title }}</strong>
                <p v-if="a.summary" style="font-size:0.82rem;color:var(--text-muted);margin-top:4px;line-height:1.5;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ a.summary }}</p>
                <div style="display:flex;gap:4px;margin-top:6px;flex-wrap:wrap;">
                  <span v-for="t in parseTags(a.tags)" :key="t" class="tag" style="font-size:0.68rem;padding:1px 6px;margin:0;">{{ t }}</span>
                </div>
              </div>
              <span style="font-size:0.75rem;color:var(--text-muted);white-space:nowrap;">{{ dayLabel(a.created_at) }}</span>
            </div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { articleAPI } from '../api'

const loading = ref(true)
const timelineItems = ref([])

const catLabel = (c) => {
  const m = { blog: '技术文章', leetcode: '算法笔记', projects: '项目', study: '学习笔记', notes: '碎碎念' }
  return m[c] || c
}
const catDotClass = (c) => {
  const m = { blog: 'dot-blog', leetcode: 'dot-leetcode', projects: 'dot-projects', study: 'dot-study', notes: 'dot-notes' }
  return m[c] || 'dot-blog'
}
const articleLink = (a) => {
  const cat = a.category || 'blog'
  return cat === 'blog' ? '/blog/' + a.id : '/' + cat + '/' + a.id
}
const dayLabel = (d) => d ? new Date(d).toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' }) : ''
const parseTags = (tags) => {
  if (!tags) return []
  return tags.split(',').map(t => t.trim()).filter(Boolean).slice(0, 5)
}

onMounted(async () => {
  try {
    const res = await articleAPI.list({ limit: 200 })
    const articles = res.data.articles || []
    // Group by year-month
    const groups = {}
    articles.forEach(a => {
      const d = new Date(a.created_at)
      const key = d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0')
      if (!groups[key]) groups[key] = { year: d.getFullYear(), month: d.getMonth() + 1, articles: [] }
      groups[key].articles.push(a)
    })
    // Sort by key descending
    timelineItems.value = Object.values(groups).sort((a, b) => {
      if (a.year !== b.year) return b.year - a.year
      return b.month - a.month
    })
  } catch (e) { console.error(e) }
  finally { loading.value = false }
})
</script>

<style scoped>
.timeline {
  position: relative;
  padding-left: 20px;
  border-left: 2px solid var(--border);
}
.timeline-group {
  position: relative;
  margin-bottom: 32px;
}
.timeline-marker {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
  margin-left: -32px;
}
.timeline-year {
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--accent);
}
.timeline-month {
  font-size: 0.95rem;
  color: var(--text);
  font-weight: 600;
}
.timeline-count {
  font-size: 0.78rem;
  color: var(--text-muted);
  background: var(--bg-hover);
  padding: 1px 10px;
  border-radius: 10px;
}
.timeline-cards {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.timeline-card {
  padding: 14px 18px !important;
  display: block;
}
.timeline-card:hover {
  text-decoration: none;
  transform: translateY(-1px);
}
.timeline-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 6px;
}
.dot-blog { background: #2d8a7b; }
.dot-leetcode { background: #e65100; }
.dot-projects { background: #6c3fb5; }
.dot-study { background: #2d7dd2; }
.dot-notes { background: #d4a574; }

@media (max-width: 768px) {
  .timeline { padding-left: 14px; }
  .timeline-marker { margin-left: -26px; }
  .timeline-card { padding: 12px 14px !important; }
}
</style>