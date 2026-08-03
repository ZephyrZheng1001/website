<template>
  <div class="article-page">
    <div style="padding-top:40px;padding-bottom:80px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">
      <div v-if="loading" style="text-align:center;color:var(--text-muted);padding:60px;">加载中...</div>

      <div v-else-if="!article" style="text-align:center;padding:60px;color:var(--text-muted);">
        <p>文章不存在</p>
        <router-link :to="backLink">返回</router-link>
      </div>

      <div v-else class="article-layout">
        <!-- TOC Sidebar -->
        <aside v-if="toc.length > 0" class="toc-sidebar">
          <div class="toc-sticky">
            <h4 class="toc-title">目录</h4>
            <nav class="toc-nav">
              <a
                v-for="(item, i) in toc"
                :key="i"
                :href="'#' + item.id"
                :class="['toc-item', 'toc-level-' + item.level, { 'toc-active': activeTocId === item.id }]"
                @click.prevent="scrollToHeading(item.id)"
              >{{ item.text }}</a>
            </nav>
          </div>
        </aside>

        <!-- Main content -->
        <article class="article-main">
          <router-link :to="backLink" style="display:inline-block;margin-bottom:24px;font-size:0.9rem;">
            ← 返回{{ backLabel }}
          </router-link>

          <h1 style="font-size:2rem;margin-bottom:16px;">{{ article.title }}</h1>

          <div style="display:flex;gap:12px;align-items:center;margin-bottom:32px;flex-wrap:wrap;">
            <span v-for="tag in parseTags(article.tags)" :key="tag" :class="getTagClass(tag)">{{ cleanTag(tag) }}</span>
            <span style="color:var(--text-muted);font-size:0.85rem;">{{ formatDate(article.created_at) }}</span>
          </div>

          <div ref="contentRef" class="article-content" v-html="renderedContent"></div>

          <!-- Prev / Next -->
          <nav v-if="prevNext.prev || prevNext.next" class="prev-next-nav">
            <router-link
              v-if="prevNext.prev"
              :to="prevLink(prevNext.prev)"
              class="prev-next-link"
            >
              <span class="prev-next-label">← 上一篇</span>
              <span class="prev-next-title">{{ prevNext.prev.title }}</span>
            </router-link>
            <div v-else></div>
            <router-link
              v-if="prevNext.next"
              :to="prevLink(prevNext.next)"
              class="prev-next-link prev-next-right"
            >
              <span class="prev-next-label">下一篇 →</span>
              <span class="prev-next-title">{{ prevNext.next.title }}</span>
            </router-link>
            <div v-else></div>
          </nav>
        </article>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import { articleAPI } from '../api'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import 'highlight.js/styles/github.css'

import javascript from 'highlight.js/lib/languages/javascript'
import python from 'highlight.js/lib/languages/python'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'
import css from 'highlight.js/lib/languages/css'
import json from 'highlight.js/lib/languages/json'
import xml from 'highlight.js/lib/languages/xml'
import sql from 'highlight.js/lib/languages/sql'
import yaml from 'highlight.js/lib/languages/yaml'
import java from 'highlight.js/lib/languages/java'
import cpp from 'highlight.js/lib/languages/cpp'

hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('js', javascript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('py', python)
hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('sh', bash)
hljs.registerLanguage('css', css)
hljs.registerLanguage('json', json)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('yml', yaml)
hljs.registerLanguage('java', java)
hljs.registerLanguage('cpp', cpp)
hljs.registerLanguage('c', cpp)
hljs.registerLanguage('mysql', sql)

const renderer = new marked.Renderer()
let headingCounter = 0
renderer.heading = function({ text, depth }) {
  const id = 'h-' + headingCounter++
  return '<h' + depth + ' id="' + id + '">' + text + '</h' + depth + '>'
}

marked.use({
  renderer,
  async: false,
  hooks: {
    preprocess(markdown) { return markdown }
  },
  extensions: [],
})

const origCode = marked.Renderer.prototype.code
marked.Renderer.prototype.code = function(token) {
  const lang = token.lang || ''
  const code = token.text || ''
  if (lang && hljs.getLanguage(lang)) {
    const result = hljs.highlight(code, { language: lang })
    return '<pre><code class="hljs language-' + lang + '">' + result.value + '</code></pre>'
  }
  const result = hljs.highlightAuto(code)
  return '<pre><code class="hljs">' + result.value + '</code></pre>'
}

const route = useRoute()
const article = ref(null)
const loading = ref(true)
const contentRef = ref(null)
const toc = ref([])
const activeTocId = ref('')
const prevNext = ref({ prev: null, next: null })

const columnLabels = { blog: '技术文章', leetcode: '算法笔记', projects: '项目', notes: '碎碎念' }

const backLink = computed(() => {
  const cat = article.value?.category || 'blog'
  return cat === 'blog' ? '/blog' : '/' + cat
})
const backLabel = computed(() => {
  const cat = article.value?.category || 'blog'
  return columnLabels[cat] || ''
})

function prevLink(item) {
  if (!item || !article.value) return '/'
  const cat = article.value.category || 'blog'
  const base = cat === 'blog' ? '/blog' : '/' + cat
  return base + '/' + item.id
}

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  headingCounter = 0
  return marked(article.value.content)
})

function getTagClass(tag) {
  const t = tag.toLowerCase()
  if (t === 'easy' || t === '\u7b80\u5355') return 'tag-difficulty diff-easy'
  if (t === 'medium' || t === '\u4e2d\u7b49') return 'tag-difficulty diff-medium'
  if (t === 'hard' || t === '\u56f0\u96be') return 'tag-difficulty diff-hard'
  return 'tag'
}
function cleanTag(tag) {
  if (!tag) return ''
  const map = { easy: '\u7b80\u5355', medium: '\u4e2d\u7b49', hard: '\u56f0\u96be' }
  return map[tag.toLowerCase()] || tag
}
function parseTags(tags) {
  if (!tags) return []
  return tags.split(',').map(t => t.trim()).filter(Boolean)
}
function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

function extractTOC() {
  if (!contentRef.value) return
  const headings = contentRef.value.querySelectorAll('h2, h3')
  toc.value = Array.from(headings).map(h => ({
    id: h.id,
    text: h.textContent,
    level: parseInt(h.tagName.charAt(1)),
  }))
}

function scrollToHeading(id) {
  const el = document.getElementById(id)
  if (el) {
    activeTocId.value = id
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

function onScroll() {
  if (!contentRef.value || toc.value.length === 0) return
  const headings = contentRef.value.querySelectorAll('h2, h3')
  let current = ''
  headings.forEach(h => {
    const rect = h.getBoundingClientRect()
    if (rect.top <= 100) current = h.id
  })
  if (current) activeTocId.value = current
}

function injectCopyButtons() {
  if (!contentRef.value) return
  contentRef.value.querySelectorAll('pre').forEach(pre => {
    if (pre.querySelector('.copy-btn')) return
    const wrapper = document.createElement('div')
    wrapper.style.position = 'relative'
    const btn = document.createElement('button')
    btn.className = 'copy-btn'
    btn.textContent = '\u590d\u5236'
    btn.onclick = () => {
      const code = pre.querySelector('code')?.textContent || pre.textContent || ''
      navigator.clipboard.writeText(code).then(() => {
        btn.textContent = '\u5df2\u590d\u5236!'
        setTimeout(() => btn.textContent = '\u590d\u5236', 2000)
      })
    }
    pre.parentNode.insertBefore(wrapper, pre)
    wrapper.appendChild(pre)
    wrapper.appendChild(btn)
  })
}

watch(renderedContent, async () => {
  await nextTick()
  extractTOC()
  injectCopyButtons()
})

onMounted(async () => {
  try {
    const res = await articleAPI.get(route.params.id)
    article.value = res.data.article || res.data
    prevNext.value = { prev: res.data.prev || null, next: res.data.next || null }
    document.title = (article.value.title || '') + ' - Zephyr'
    await nextTick()
    extractTOC()
    injectCopyButtons()
    window.addEventListener('scroll', onScroll, { passive: true })
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
<style scoped>
.article-layout {
  display: flex;
  gap: 40px;
  position: relative;
}
.article-main {
  flex: 1;
  min-width: 0;
}
.toc-sidebar {
  width: 180px;
  flex-shrink: 0;
}
.toc-sticky {
  position: sticky;
  top: 80px;
}
.toc-title {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 12px;
}
.toc-nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-left: 2px solid var(--border);
  padding-left: 12px;
}
.toc-item {
  font-size: 0.8rem;
  color: var(--text-muted);
  padding: 3px 0;
  transition: color 0.15s;
  display: block;
  line-height: 1.5;
}
.toc-item:hover { color: var(--accent); text-decoration: none; }
.toc-item.toc-active { color: var(--accent); font-weight: 500; }
.toc-level-3 { padding-left: 14px; font-size: 0.76rem; }

:deep(.copy-btn) {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 12px;
  background: rgba(0,0,0,0.06);
  border: 1px solid var(--border);
  border-radius: 4px;
  font-size: 0.75rem;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}
:deep(.copy-btn:hover) {
  background: var(--accent-light);
  color: var(--accent);
  border-color: var(--accent);
}
[data-theme="dark"] :deep(.copy-btn) { background: rgba(255,255,255,0.06); }

/* Prev/Next Nav */
.prev-next-nav {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
}
.prev-next-link {
  display: block;
  padding: 12px 16px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  transition: all 0.2s;
}
.prev-next-link:hover {
  border-color: var(--accent);
  box-shadow: var(--shadow-md);
  text-decoration: none;
}
.prev-next-right { text-align: right; }
.prev-next-label {
  display: block;
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-bottom: 4px;
}
.prev-next-title {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text);
}

@media (max-width: 900px) {
  .article-layout { flex-direction: column; }
  .toc-sidebar { width: 100%; order: -1; }
  .toc-sticky { position: static; }
  .toc-nav { flex-direction: row; flex-wrap: wrap; gap: 8px; border-left: none; padding-left: 0; border-bottom: 1px solid var(--border); padding-bottom: 12px; margin-bottom: 20px; }
  .toc-level-3 { padding-left: 0; }
}
</style>