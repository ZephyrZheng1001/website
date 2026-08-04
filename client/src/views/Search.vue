<template>
  <div style="padding-top:24px;padding-bottom:60px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">
    <h1 style="font-size:2rem;margin-bottom:24px;">🔍 搜索</h1>
    <div style="display:flex;gap:10px;margin-bottom:32px;">
      <div class="form-group" style="flex:1;margin-bottom:0;">
        <input v-model="query" type="text" placeholder="搜索文章标题、内容、标签..." @keyup.enter="search" style="padding:12px 16px;font-size:1rem;" autofocus />
      </div>
      <button class="btn btn-primary" @click="search" :disabled="!query.trim()" style="padding:12px 24px;font-size:0.95rem;">搜索</button>
    </div>
    <div v-if="loading" style="text-align:center;color:var(--text-muted);padding:40px;">搜索中...</div>
    <div v-else-if="searched && results.length === 0" style="text-align:center;color:var(--text-muted);padding:40px;">
      <p style="font-size:1.1rem;margin-bottom:8px;">未找到相关内容</p>
      <p>试试换个关键词？</p>
    </div>
    <div v-else-if="results.length > 0">
      <p style="color:var(--text-muted);margin-bottom:16px;font-size:0.9rem;">找到 {{ results.length }} 篇相关文章</p>
      <div style="display:flex;flex-direction:column;gap:12px;">
        <router-link v-for="article in results" :key="article.id" :to="articleLink(article)" class="card" style="display:block;">
          <div style="display:flex;align-items:flex-start;gap:12px;">
            <span style="font-size:1.1rem;flex-shrink:0;">{{ catIcon(article.category) }}</span>
            <div style="flex:1;min-width:0;">
              <h3 style="font-size:1rem;color:var(--text);margin-bottom:4px;">{{ article.title }}</h3>
              <p style="color:var(--text-muted);font-size:0.85rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ article.summary || (article.content || '').slice(0, 100) + '...' }}</p>
              <div style="display:flex;align-items:center;gap:8px;margin-top:6px;flex-wrap:wrap;">
                <span style="font-size:0.75rem;color:var(--accent);background:var(--accent-light);padding:1px 8px;border-radius:8px;">{{ catName(article.category) }}</span>
                <span v-for="tag in parseTags(article.tags)" :key="tag" style="font-size:0.75rem;color:var(--text-muted);">#{{ tag }}</span>
                <span style="font-size:0.75rem;color:var(--text-muted);margin-left:auto;">{{ formatDate(article.created_at) }}</span>
              </div>
            </div>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { articleAPI } from '../api'

const query = ref('')
const results = ref([])
const loading = ref(false)
const searched = ref(false)

function parseTags(tags) { if (!tags) return []; return tags.split(',').map(t => t.trim()).filter(Boolean).slice(0, 8) }
function formatDate(d) { return d ? new Date(d).toLocaleDateString('zh-CN') : '' }
function catIcon(cat) { const icons = { blog: '📝', leetcode: '💡', projects: '🚀', study: '📖', notes: '💬' }; return icons[cat] || '📄' }
function catName(cat) { const names = { blog: '技术文章', leetcode: '算法笔记', projects: '项目', study: '学习笔记', notes: '碎碎念' }; return names[cat] || cat }
function articleLink(article) { const cat = article.category; if (cat === 'blog') return '/blog/' + article.id; return '/' + cat + '/' + article.id }

async function search() {
  const q = query.value.trim()
  if (!q) return
  loading.value = true; searched.value = true
  try { const res = await articleAPI.search(q); results.value = res.data || [] }
  catch (e) { console.error(e); results.value = [] }
  finally { loading.value = false }
}
</script>