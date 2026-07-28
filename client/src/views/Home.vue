<template>
  <div style="padding-top:48px;padding-bottom:32px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">

    <!-- Hero -->
    <section style="margin-bottom:48px;">
      <h1 style="font-size:2.4rem;font-weight:800;letter-spacing:-0.8px;line-height:1.25;margin-bottom:14px;">
        后端开发爱好者，<br/>也喜欢技术和文字。
      </h1>
      <p style="font-size:1.05rem;color:var(--text-muted);max-width:520px;line-height:1.7;">
        主业后端开发，平时爱写博客、刷算法、做项目。
      </p>
    </section>

    <!-- Stats line -->
    <section style="display:flex;gap:40px;padding:16px 0;border-top:1px solid var(--border);border-bottom:1px solid var(--border);margin-bottom:40px;font-size:0.88rem;color:var(--text-muted);">
      <div><strong style="color:var(--accent);">{{ stats.articles }}</strong> 篇文章</div>
      <div><strong style="color:var(--accent);">Always</strong> 在线</div>
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
    <section v-if="latestList.length" style="margin-top:48px;">
      <h3 style="font-size:1rem;color:var(--text-muted);font-weight:500;margin-bottom:16px;">最近更新</h3>
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

const latestList = ref([])
const stats = ref({ articles: 0 })

const colIcons = { blog: '📝', leetcode: '💡', projects: '🚀', notes: '💬' }
const colNames = { blog: '技术文章', leetcode: '算法笔记', projects: '项目', notes: '碎碎念' }

function formatDate(d) {
  return d ? new Date(d).toLocaleDateString('zh-CN') : ''
}

onMounted(async () => {
  try {
    const results = await Promise.all([
      articleAPI.list({ page: 1, limit: 1, category: 'blog' }),
      articleAPI.list({ page: 1, limit: 1, category: 'leetcode' }),
      articleAPI.list({ page: 1, limit: 1, category: 'projects' }),
      articleAPI.list({ page: 1, limit: 1, category: 'notes' }),
    ])
    let total = 0
    const cats = ['blog', 'leetcode', 'projects', 'notes']
    cats.forEach((cat, i) => {
      total += results[i].data.total
      const a = results[i].data.articles?.[0]
      if (a) {
        latestList.value.push({
          ...a,
          column: colNames[cat],
          icon: colIcons[cat],
          link: cat === 'blog' ? `/blog/${a.id}` : `/${cat}/${a.id}`,
        })
      }
    })
    stats.value.articles = total
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
  padding: 16px;
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
</style>
