<template>
  <div style="padding-top:24px;padding-bottom:60px;max-width:900px;margin:0 auto;padding-left:24px;padding-right:24px;">
    <h1 style="font-size:2rem;margin-bottom:8px;">{{ columnInfo.icon }} {{ columnInfo.title }}</h1>
    <p style="color:var(--text-muted);margin-bottom:32px;">{{ columnInfo.desc }}</p>

    <div class="blog-layout">
      <!-- Sidebar -->
      <aside class="blog-sidebar">
        <nav class="sidebar-nav">
          <router-link
            v-for="col in allColumns"
            :key="col.key"
            :to="'/' + col.key"
            :class="['sidebar-link', { active: category === col.key }]"
          >
            <span class="sidebar-icon">{{ col.icon }}</span>
            <span>{{ col.title }}</span>
            <span class="sidebar-count" v-if="col.count !== null">{{ col.count }}</span>
          </router-link>
        </nav>
        <div style="margin-top:20px;">
          <h4 style="font-size:0.8rem;font-weight:600;color:var(--text-muted);text-transform:uppercase;letter-spacing:1px;margin-bottom:10px;">标签</h4>
          <div style="display:flex;flex-wrap:wrap;gap:6px;">
            <router-link
              v-for="t in popularTags"
              :key="t"
              :to="`/${routeName}?tag=${encodeURIComponent(t)}`"
              :class="['tag', { 'tag-active': activeTag === t }]"
              style="cursor:pointer;"
            >{{ t }}</router-link>
            <router-link
              v-if="activeTag"
              :to="`/${routeName}`"
              class="tag tag-clear"
              style="cursor:pointer;"
            >清除筛选</router-link>
          </div>
        </div>
      </aside>

      <!-- Content -->
      <div class="blog-content">
        <div v-if="loading" style="text-align:center;color:var(--text-muted);padding:60px;">加载中...</div>

        <div v-else-if="articles.length === 0" style="text-align:center;color:var(--text-muted);padding:60px;">
          <p style="font-size:1.1rem;">还没有内容</p>
          <router-link to="/admin">去写一篇 →</router-link>
        </div>

        <div v-else style="display:flex;flex-direction:column;gap:16px;">
          <router-link
            v-for="article in articles"
            :key="article.id"
            :to="`/${routeName}/${article.id}`"
            class="card"
            style="display:block;"
          >
            <h2 style="font-size:1.2rem;margin-bottom:8px;color:var(--text);">{{ article.title }}</h2>
            <p style="color:var(--text-muted);font-size:0.9rem;margin-bottom:10px;line-height:1.5;">
              {{ article.summary || article.content?.slice(0, 150) + '...' }}
            </p>
            <div style="display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:8px;">
              <div>
                <span v-for="tag in parseTags(article.tags)" :key="tag" :class="getTagClass(tag)" style="cursor:pointer;" @click.prevent="goToTag(tag)">{{ cleanTag(tag) }}</span>
              </div>
              <span style="color:var(--text-muted);font-size:0.8rem;">{{ formatDate(article.created_at) }}</span>
            </div>
          </router-link>
        </div>

        <div v-if="totalPages > 1" class="pagination">
          <button :disabled="page <= 1" @click="page--;fetchArticles()">上一页</button>
          <button v-for="p in totalPages" :key="p" @click="page=p;fetchArticles()" :class="{active:p===page}">{{ p }}</button>
          <button :disabled="page >= totalPages" @click="page++;fetchArticles()">下一页</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { articleAPI } from '../api'

const route = useRoute()
const router = useRouter()

const columns = {
  blog: { title: '技术文章', desc: '技术分享与教程', icon: '📝' },
  leetcode: { title: '算法笔记', desc: '算法题解与刷题笔记', icon: '💡' },
  projects: { title: '项目', desc: '项目复盘与技术方案', icon: '🚀' },
  notes: { title: '碎碎念', desc: '日常随想与生活记录', icon: '💬' },
  study: { title: '学习笔记', desc: '系统化学习记录', icon: '📓' },
}

const allColumns = ref([
  { key: 'blog', title: '技术文章', icon: '📝', count: null },
  { key: 'leetcode', title: '算法笔记', icon: '💡', count: null },
  { key: 'projects', title: '项目', icon: '🚀', count: null },
  { key: 'notes', title: '碎碎念', icon: '💬', count: null },
])

const routeName = computed(() => route.name?.toLowerCase() || 'blog')
const category = computed(() => routeName.value)
const columnInfo = computed(() => columns[category.value] || columns.blog)
const activeTag = computed(() => route.query.tag || '')

const articles = ref([])
const loading = ref(true)
const page = ref(1)
const limit = 10
const total = ref(0)
const totalPages = ref(0)
const popularTags = ref([])

function getTagClass(tag) {
  const t = tag.toLowerCase()
  if (t === 'easy' || t === '简单') return 'tag-difficulty diff-easy'
  if (t === 'medium' || t === '中等') return 'tag-difficulty diff-medium'
  if (t === 'hard' || t === '困难') return 'tag-difficulty diff-hard'
  return 'tag'
}
function cleanTag(tag) {
  const map = { easy: '简单', medium: '中等', hard: '困难' }
  return map[tag.toLowerCase()] || tag
}
function parseTags(tags) {
  if (!tags) return []
  return tags.split(',').map(t => t.trim()).filter(Boolean)
}
function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('zh-CN')
}

function goToTag(tag) {
  // Navigate to the same column with that tag filter
  router.push(`/${routeName.value}?tag=${encodeURIComponent(tag)}`)
}

async function fetchArticles() {
  loading.value = true
  try {
    const params = { page: page.value, limit, category: category.value }
    if (route.query.tag) { params.tag = route.query.tag }
    const res = await articleAPI.list(params)
    articles.value = res.data.articles || []
    total.value = res.data.total
    totalPages.value = Math.ceil(total.value / limit)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

// Load ALL tags for the current category (unfiltered), so sidebar tags don't disappear
async function fetchAllTags() {
  try {
    const res = await articleAPI.list({ category: category.value, limit: 200 })
    const tagSet = new Set()
    ;(res.data.articles || []).forEach(a => parseTags(a.tags).forEach(t => tagSet.add(t)))
    popularTags.value = [...tagSet].slice(0, 20)
  } catch (e) { /* silent */ }
}

async function fetchColumnCounts() {
  const cats = ['blog', 'leetcode', 'projects', 'notes']
  const results = await Promise.allSettled(
    cats.map(c => articleAPI.list({ category: c, limit: 1 }))
  )
  // Trigger reactivity by replacing entire array
  const updated = allColumns.value.map((col, i) => {
    const r = results[i]
    const total = (r.status === 'fulfilled' && r.value?.data?.total !== undefined)
      ? r.value.data.total : col.count
    return { ...col, count: total }
  })
  allColumns.value = updated
}

// Watch both category and tag changes
watch(() => [category.value, route.query.tag], () => { page.value = 1; fetchArticles() })
onMounted(() => { fetchArticles(); fetchAllTags(); fetchColumnCounts() })
watch(category, () => { fetchAllTags(); fetchColumnCounts() })
watch(() => route.fullPath, () => { fetchColumnCounts() })
fetchColumnCounts()
</script>

<style scoped>
.blog-layout {
  display: flex;
  gap: 32px;
}
.blog-sidebar {
  width: 180px;
  flex-shrink: 0;
}
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.sidebar-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 0.88rem;
  color: var(--text-muted);
  transition: all 0.15s;
}
.sidebar-link:hover {
  background: var(--bg-hover);
  color: var(--text);
  text-decoration: none;
}
.sidebar-link.active {
  background: var(--accent-light);
  color: var(--accent);
  font-weight: 600;
}
.sidebar-icon { font-size: 0.9rem; }
.sidebar-count {
  margin-left: auto;
  font-size: 0.75rem;
  background: var(--bg-hover);
  padding: 1px 7px;
  border-radius: 10px;
  color: var(--text-muted);
}
.sidebar-link.active .sidebar-count {
  background: rgba(45,138,123,0.15);
  color: var(--accent);
}
.blog-content {
  flex: 1;
  min-width: 0;
}
.tag-active {
  background: var(--accent) !important;
  color: #fff !important;
}
.tag-clear {
  background: var(--bg-hover);
  color: var(--text-muted);
  border: 1px dashed var(--border);
}
.tag-clear:hover {
  border-color: var(--accent);
  color: var(--accent);
}
@media (max-width: 768px) {
  .blog-layout { flex-direction: column; }
  .blog-sidebar { width: 100%; }
  .sidebar-nav { flex-direction: row; flex-wrap: wrap; gap: 4px; margin-bottom: 20px; }
  .sidebar-link { padding: 6px 12px; font-size: 0.82rem; }
  .sidebar-count { display: none; }
}
</style>