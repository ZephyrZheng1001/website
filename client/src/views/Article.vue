<template>
  <div class="article-page">
    <div style="padding-top:40px;padding-bottom:80px;max-width:960px;margin:0 auto;padding-left:24px;padding-right:24px;">
      <div v-if="loading" style="text-align:center;color:var(--text-muted);padding:60px;">加载中...</div>

      <div v-else-if="!article" style="text-align:center;padding:60px;color:var(--text-muted);">
        <p>文章不存在</p>
        <router-link :to="backLink">返回</router-link>
      </div>

      <div v-else class="article-layout">
        <!-- TOC Sidebar -->
        <aside v-if="toc.length > 1" class="toc-sidebar">
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
            <span style="color:var(--text-muted);font-size:0.85rem;">{{ formatDate(article.created_at) }} · {{ readingTime(article.content) }} · {{ article.view_count || 0 }} 次阅读</span>
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
import katex from 'katex'
import 'katex/dist/katex.min.css'
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
hljs.registerLanguage('lua', function(hljs) {
  var LUA_KEYWORDS = {
    keyword: 'and break do else elseif end false for function goto if in local nil not or repeat return then true until while',
    literal: 'true false nil',
    built_in: '_G _VERSION assert collectgarbage dofile error getfenv getmetatable ipairs load loadfile loadstring module next pairs pcall print rawequal rawget rawlen rawset require select setfenv setmetatable tonumber tostring type unpack xpcall coroutine debug io math os package string table'
  };
  return {
    name: 'Lua',
    aliases: ['lua'],
    keywords: LUA_KEYWORDS,
    contains: [
      hljs.HASH_COMMENT_MODE,
      hljs.C_NUMBER_MODE,
      hljs.QUOTE_STRING_MODE,
      { className: 'string', begin: '"', end: '"' },
      { className: 'string', begin: "'", end: "'" }
    ]
  };
})

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
    return '<div class="code-block-wrapper"><span class="code-lang-tag">' + lang + '</span><pre><code class="hljs language-' + lang + '">' + result.value + '</code></pre></div>'
  }
  if (lang) {
    // Unknown language - still show tag, use highlightAuto
    const result = hljs.highlightAuto(code)
    return '<div class="code-block-wrapper"><span class="code-lang-tag">' + lang + '</span><pre><code class="hljs">' + result.value + '</code></pre></div>'
  }
  const result = hljs.highlightAuto(code)
  return '<div class="code-block-wrapper"><pre><code class="hljs">' + result.value + '</code></pre></div>'
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
  let html = marked(article.value.content)
  html = processInlineKatex(html)
  return html
})

function processInlineKatex(html) {
  // display math $$
  html = html.replace(/\$\$([\s\S]*?)\$\$/g, function(_, formula) {
    try {
      return '<p class="katex-block">' + katex.renderToString(formula.trim(), { displayMode: true, throwOnError: false }) + '</p>'
    } catch (e) {
      return '<code class="katex-error">' + _ + '</code>'
    }
  })
  // inline math $
  html = html.replace(/(?<![`$])\$(?!\$)([^$]+?)\$(?![`$])/g, function(_, formula) {
    try {
      return katex.renderToString(formula.trim(), { displayMode: false, throwOnError: false })
    } catch (e) {
      return '<code class="katex-error">' + _ + '</code>'
    }
  })
  return html
}

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
function readingTime(content) {
  if (!content) return '1 min'
  const text = content.replace(/<[^>]*>/g, '').replace(/[#*_`~\[\]()>\-!|]/g, '')
  const cnChars = (text.match(/[\u4e00-\u9fff]/g) || []).length
  const enWords = text.replace(/[\u4e00-\u9fff\u3400-\u4dbf\uf900-\ufaff]/g, '').split(/\s+/).filter(Boolean).length
  const mins = Math.ceil((cnChars / 400) + (enWords / 200))
  return (mins || 1) + ' min'
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

function injectImageZoom() {
  if (!contentRef.value) return
  contentRef.value.querySelectorAll('img').forEach(img => {
    if (img.closest('.katex-block') || img.closest('.mermaid-container') || img.classList.contains('no-zoom')) return
    img.loading = 'lazy'
    img.style.cursor = 'zoom-in'
    img.addEventListener('click', () => {
      const overlay = document.createElement('div')
      overlay.className = 'img-overlay'
      overlay.innerHTML = '<img src="' + img.src + '" style="max-width:90vw;max-height:90vh;object-fit:contain;border-radius:8px;" />'
      overlay.addEventListener('click', () => overlay.remove())
      document.body.appendChild(overlay)
    })
  })
}

function injectCopyButtons() {
  if (!contentRef.value) return
  var wrappers = contentRef.value.querySelectorAll('.code-block-wrapper')
  for (var i = 0; i < wrappers.length; i++) {
    var w = wrappers[i]
    if (w.querySelector('.copy-btn')) continue
    var btn = document.createElement('button')
    btn.className = 'copy-btn'
    btn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>'
    btn.onclick = function() {
      var code = this.parentNode.querySelector('code')
      var txt = code ? code.textContent : this.parentNode.textContent || ''
      navigator.clipboard.writeText(txt).then((function(btn) {
        return function() { btn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>'; setTimeout(function() { btn.innerHTML = '<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>' }, 2000) }
      })(this))
    }
    w.appendChild(btn)
  }
}
function _injectMermaid_old() {

  if (!contentRef.value) return
  const blocks = contentRef.value.querySelectorAll('.mermaid')
  if (blocks.length === 0) return

  function doIt() {
    window.mermaid.initialize({ startOnLoad: false, theme: 'default' })
    window.mermaid.run({ nodes: Array.from(blocks) })
  }

  if (mermaidLoaded) {
    doIt()
  } else {
    const s = document.createElement('script')
    s.src = 'https://unpkg.com/mermaid@11/dist/mermaid.min.js'
    s.onload = () => { mermaidLoaded = true; doIt() }
    document.head.appendChild(s)
  }
}

watch(renderedContent, async () => {
  await nextTick()
  await nextTick()
  extractTOC()
  injectCopyButtons()
  injectMermaid()
  injectImageZoom()
  // Fallback: sometimes v-html DOM isn't ready after nextTick
  setTimeout(function() { injectCopyButtons() }, 100)
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
    injectMermaid()
    injectImageZoom()
    window.addEventListener('scroll', function(e) { onScroll() }, { passive: true })
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})

// Watch for route param changes (same component, different article)
watch(() => route.params.id, async (newId) => {
  if (!newId) return
  loading.value = true
  article.value = null
  window.scrollTo({ top: 0 })
  try {
    const res = await articleAPI.get(newId)
    article.value = res.data.article || res.data
    prevNext.value = { prev: res.data.prev || null, next: res.data.next || null }
    document.title = (article.value.title || '') + ' - Zephyr'
    await nextTick()
    await nextTick()
    extractTOC()
    injectCopyButtons()
    injectMermaid()
    injectImageZoom()
    setTimeout(function() { injectCopyButtons() }, 100)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
<style scoped>
.article-layout {
  position: relative;
}
.toc-sidebar {
  position: fixed;
  left: max(24px, calc((100vw - 960px) / 2 - 180px));
  top: 80px;
  width: 160px;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  z-index: 50;
}
.toc-sticky {
  /* noop */
}
.toc-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 8px;
}
.toc-nav {
  display: flex;
  flex-direction: column;
  gap: 1px;
  border-left: 2px solid var(--border);
  padding-left: 10px;
}
.toc-item {
  font-size: 0.88rem;
  color: var(--text-muted);
  padding: 2px 0;
  transition: color 0.15s;
  display: block;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.toc-item:hover { color: var(--accent); text-decoration: none; }
.toc-item.toc-active { color: var(--accent); font-weight: 500; }
.toc-level-3 { padding-left: 10px; font-size: 0.85rem; }


:deep(.code-block-wrapper) .copy-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 32px;
  height: 32px;
  padding: 7px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-muted);
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s ease, background 0.2s, color 0.2s, border-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}
:deep(.code-block-wrapper:hover) .copy-btn {
  opacity: 1;
}
:deep(.copy-btn:hover) {
  background: var(--accent-light);
  color: var(--accent);
  border-color: var(--accent);
}
[data-theme="dark"] :deep(.code-block-wrapper) .copy-btn { background: var(--bg-card); border-color: var(--border); }
[data-theme="dark"] :deep(.copy-btn:hover) { background: rgba(77,184,165,0.15); color: var(--accent); }

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
  padding: 14px 18px;
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
  
  
}

/* Image zoom overlay */
:deep(.img-overlay) {
  position: fixed; inset: 0;
  background: rgba(0,0,0,0.85);
  display: flex; align-items: center; justify-content: center;
  z-index: 9999; cursor: zoom-out;
}
:deep(.img-overlay img) { box-shadow: 0 8px 40px rgba(0,0,0,0.5); }
[data-theme="dark"] :deep(.img-overlay) { background: rgba(0,0,0,0.92); }

:deep(.code-block-wrapper) { position: relative; margin: 1em 0; }
:deep(.code-lang-tag) {
  position: absolute; top: 0; left: 12px;
  background: var(--accent); color: #fff;
  font-size: 0.7rem; font-weight: 600;
  padding: 2px 10px; border-radius: 0 0 4px 4px;
  text-transform: uppercase;
  z-index: 1;
}
[data-theme="dark"] :deep(.code-lang-tag) { background: #4db8a5; }
@media (max-width: 768px) {
  .toc-sidebar { display: none !important; }
}
</style>
@media (max-width: 768px) {
  .toc-sidebar { display: none; }
}
