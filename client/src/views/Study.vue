<template>
  <div style="padding-top:24px;padding-bottom:60px;max-width:960px;margin:0 auto;padding-left:24px;padding-right:24px;">
    <div style="display:flex;align-items:center;gap:12px;margin-bottom:8px;">
      <h1 style="font-size:1.6rem;">📖 学习笔记</h1>
      <span style="font-size:0.85rem;color:var(--text-muted);background:var(--bg-hover);padding:4px 12px;border-radius:12px;">
        {{ totalArticles }} 篇 · {{ doneCount }} 已完成
      </span>
    </div>
    <p style="color:var(--text-muted);margin-bottom:8px;">系统化学习记录，告别碎片化</p>

    <div v-if="totalArticles > 0" style="margin-bottom:28px;">
      <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:6px;">
        <span style="font-size:0.85rem;color:var(--text-muted);">总体进度</span>
        <span style="font-size:0.85rem;font-weight:600;color:var(--accent);">{{ overallPct }}%</span>
      </div>
      <div style="height:8px;background:var(--bg-hover);border-radius:4px;overflow:hidden;">
        <div :style="{width:overallPct+'%',background:'linear-gradient(90deg,var(--accent),#5bb8a8)',borderRadius:'4px',height:'100%',transition:'width 0.5s ease'}"></div>
      </div>
    </div>

    <div class="study-layout">
      <aside class="study-sidebar">
        <nav class="sidebar-nav">
          <button v-for="sub in subcategories" :key="sub.key" :class="['sidebar-link', { active: activeSub === sub.key }]" @click="activeSub = sub.key">
            <span class="sidebar-icon">{{ sub.icon }}</span>
            <span>{{ sub.name }}</span>
            <span class="sidebar-count">{{ sub.count }}</span>
          </button>
        </nav>
      </aside>

      <div class="study-content">
        <div style="display:flex;gap:8px;margin-bottom:20px;">
          <button v-for="s in statusFilters" :key="s.key" :class="['btn', 'btn-sm', activeStatus === s.key ? 'btn-primary' : 'btn-outline']" @click="activeStatus = s.key">{{ s.label }}</button>
        </div>

        <div v-if="loading" style="text-align:center;color:var(--text-muted);padding:60px;">加载中...</div>

        <div v-else-if="filteredArticles.length === 0" style="text-align:center;color:var(--text-muted);padding:60px;">
          <p style="font-size:1.2rem;margin-bottom:12px;">📭 {{ activeSub === 'all' ? '还没有学习笔记' : '这个分类下还没有笔记' }}</p>
          <router-link to="/admin" class="btn btn-primary">去写一篇 →</router-link>
        </div>

        <div v-else style="display:flex;flex-direction:column;gap:16px;">
          <router-link v-for="article in filteredArticles" :key="article.id" :to="'/study/' + article.id" class="card study-card" style="display:flex;align-items:center;gap:16px;padding:18px 22px;">
            <div class="study-status-dot" :class="statusClass(article.study_status)" :title="statusLabel(article.study_status)"></div>
            <div style="flex:1;min-width:0;">
              <div style="display:flex;align-items:center;gap:8px;margin-bottom:4px;">
                <span v-if="article.is_pinned" style="font-size:0.75rem;" title="置顶">📌</span>
                <span style="font-size:0.95rem;font-weight:600;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ article.is_pinned ? '📌 ' : '' }}{{ article.title }}</span>
              </div>
              <div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap;">
                <span v-if="article.subcategory" style="font-size:0.75rem;color:var(--accent);background:var(--accent-light);padding:1px 8px;border-radius:8px;">{{ subNames[article.subcategory] || article.subcategory }}</span>
                <span v-for="tag in parseTags(article.tags)" :key="tag" style="font-size:0.75rem;color:var(--text-muted);">#{{ tag }}</span>
                <span style="font-size:0.75rem;color:var(--text-muted);margin-left:auto;">{{ formatDate(article.created_at) }} · {{ readingTime(article.content || article.summary) }}</span>
              </div>
            </div>
            <span style="color:var(--text-muted);font-size:0.85rem;">→</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { articleAPI, studyCategoryAPI } from '../api'

const articles = ref([])
const loading = ref(true)
const activeSub = ref('all')
const activeStatus = ref('all')

const subcategories = ref([
  { key: 'all', name: '全部', icon: '📚', count: 0 },
  { key: 'go', name: 'Go 语言', icon: '🔷', count: 0 },
  { key: 'rust', name: 'Rust', icon: '🦀', count: 0 },
  { key: 'python', name: 'Python', icon: '🐍', count: 0 },
  { key: 'system', name: '系统设计', icon: '🏗️', count: 0 },
  { key: 'frontend', name: '前端', icon: '🎨', count: 0 },
  { key: 'devops', name: 'DevOps', icon: '⚙️', count: 0 },
  { key: 'db', name: '数据库', icon: '🗄️', count: 0 },
  { key: 'cs', name: '计算机基础', icon: '💻', count: 0 },
  { key: 'other', name: '其他', icon: '📋', count: 0 },
])

const subNames = { go: 'Go 语言', rust: 'Rust', python: 'Python', system: '系统设计', frontend: '前端', devops: 'DevOps', db: '数据库', cs: '计算机基础', other: '其他' }

const statusFilters = [
  { key: 'all', label: '全部' },
  { key: 'doing', label: '🔄 学习中' },
  { key: 'todo', label: '📋 待开始' },
  { key: 'done', label: '✅ 已完成' },
]

const totalArticles = computed(() => articles.value.length)
const doneCount = computed(() => articles.value.filter(a => a.study_status === 'done').length)
const overallPct = computed(() => totalArticles.value > 0 ? Math.round((doneCount.value / totalArticles.value) * 100) : 0)

const filteredArticles = computed(() => {
  let list = articles.value
  if (activeSub.value !== 'all') list = list.filter(a => a.subcategory === activeSub.value)
  if (activeStatus.value !== 'all') list = list.filter(a => a.study_status === activeStatus.value || (activeStatus.value === 'doing' && !a.study_status))
  return list
})

function parseTags(tags) {
  if (!tags) return []
  return tags.split(',').map(t => t.trim()).filter(Boolean).slice(0, 5)
}
function readingTime(content) {
  if (!content) return '1 min'
  const text = content.replace(/<[^>]*>/g, '').replace(/[#*_`~[\]()>\-!|]/g, '')
  const cnChars = (text.match(/[\u4e00-\u9fff]/g) || []).length
  const enWords = text.replace(/[\u4e00-\u9fff]/g, '').split(/\s+/).filter(Boolean).length
  const mins = Math.max(1, Math.ceil((cnChars / 400) + (enWords / 200)))
  return mins + ' min'
}
function formatDate(d) { return d ? new Date(d).toLocaleDateString('zh-CN') : '' }
function statusClass(s) { if (s === 'done') return 'status-done'; if (s === 'todo') return 'status-todo'; return 'status-doing' }
function statusLabel(s) { if (s === 'done') return '已完成'; if (s === 'todo') return '待开始'; return '学习中' }

async function fetchData() {
  loading.value = true
  try {
    const [artRes, catRes] = await Promise.all([
      articleAPI.list({ category: 'study', limit: 200 }),
      studyCategoryAPI.list()
    ])
    articles.value = artRes.data.articles || []
    const cats = catRes.data.categories || []

    const counts = { all: articles.value.length }
    articles.value.forEach(a => {
      const sub = a.subcategory || 'other'
      counts[sub] = (counts[sub] || 0) + 1
    })

    const dynamicSubs = [{ key: 'all', name: '全部', icon: '📚', count: counts.all || 0 }]
    cats.forEach(cat => {
      dynamicSubs.push({ key: cat.name, name: cat.name, icon: cat.icon || '📋', count: counts[cat.name] || 0 })
    })
    // Add 'other' if there are articles with unknown subcategory
    if (counts['other'] > 0) {
      dynamicSubs.push({ key: 'other', name: '其他', icon: '📋', count: counts['other'] })
    }
    subcategories.value = dynamicSubs
  } catch (e) { console.error(e) } finally { loading.value = false }
}

onMounted(fetchData)
</script>

<style scoped>
.study-layout { display: flex; gap: 32px; }
.study-sidebar { width: 180px; flex-shrink: 0; }
.sidebar-nav { display: flex; flex-direction: column; gap: 2px; }
.sidebar-link { display: flex; align-items: center; gap: 8px; padding: 9px 14px; border-radius: 6px; font-size: 0.88rem; color: var(--text-muted); transition: all 0.15s; background: none; border: none; cursor: pointer; text-align: left; width: 100%; font-family: inherit; }
.sidebar-link:hover { background: var(--bg-hover); color: var(--text); }
.sidebar-link.active { background: var(--accent-light); color: var(--accent); font-weight: 600; }
.sidebar-icon { font-size: 0.9rem; }
.sidebar-count { margin-left: auto; font-size: 0.75rem; background: var(--bg-hover); padding: 1px 7px; border-radius: 10px; color: var(--text-muted); }
.sidebar-link.active .sidebar-count { background: rgba(45,138,123,0.15); color: var(--accent); }
.study-content { flex: 1; min-width: 0; }
.study-card { cursor: pointer; }
.study-card:hover { text-decoration: none; }
.study-status-dot { width: 12px; height: 12px; border-radius: 50%; flex-shrink: 0; border: 2px solid; }
.status-doing { background: #fff3e0; border-color: #e65100; }
.status-todo { background: #f5f5f5; border-color: #999; }
.status-done { background: #e8f5e9; border-color: #2e7d32; }
[data-theme="dark"] .status-doing { background: #3a2a10; }
[data-theme="dark"] .status-todo { background: #222; border-color: #666; }
[data-theme="dark"] .status-done { background: #1a3a1a; }
@media (max-width: 768px) { .study-layout { flex-direction: column; } .study-sidebar { width: 100%; } .sidebar-nav { flex-direction: row; flex-wrap: wrap; gap: 4px; margin-bottom: 20px; } .sidebar-link { padding: 6px 12px; font-size: 0.82rem; width: auto; } .sidebar-count { display: none; } }
</style>