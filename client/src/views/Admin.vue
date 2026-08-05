<template>
  <div style="padding-top:32px;padding-bottom:80px;max-width:1100px;margin:0 auto;padding-left:24px;padding-right:24px;">
    <!-- Login -->
    <div v-if="!isLoggedIn" style="max-width:400px;margin:80px auto;">
      <h1 style="font-size:1.8rem;margin-bottom:24px;text-align:center;">管理员登录</h1>
      <div class="card">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="loginForm.username" type="text" placeholder="admin" @keyup.enter="doLogin" />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="loginForm.password" type="password" placeholder="密码" @keyup.enter="doLogin" />
        </div>
        <p v-if="loginErr" style="color:#e05555;font-size:0.85rem;margin-bottom:12px;">{{ loginErr }}</p>
        <button class="btn btn-primary" style="width:100%;" @click="doLogin" :disabled="loginLoading">
          {{ loginLoading ? '登录中...' : '登录' }}
        </button>
      </div>
    </div>

    <!-- Admin Panel -->
    <div v-else>
      <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:32px;">
        <h1 style="font-size:1.8rem;">管理后台</h1>
        <button class="btn btn-outline btn-sm" @click="auth.logout();router.push('/')">退出登录</button>
      </div>

      <!-- Tabs -->
      <div style="display:flex;gap:8px;margin-bottom:24px;">
        <button class="btn" :class="tab==='articles'?'btn-primary':'btn-outline'" @click="tab='articles'">文章管理</button>
        <button class="btn" :class="tab==='write'?'btn-primary':'btn-outline'" @click="openWrite()">写文章</button>
        <button class="btn" :class="tab==='studyCats'?'btn-primary':'btn-outline'" @click="tab='studyCats'">子分类</button>
        <button class="btn" :class="tab==='password'?'btn-primary':'btn-outline'" @click="tab='password'">改密码</button>
      </div>

      <!-- Articles list -->
      <div v-if="tab==='articles'">
        <!-- Filters -->
        <div style="display:flex;gap:12px;align-items:center;margin-bottom:16px;flex-wrap:wrap;">
          <input
            v-model="articleFilter.search"
            type="text"
            placeholder="搜索标题/内容/摘要..."
            style="flex:1;min-width:200px;padding:8px 12px;border:1px solid var(--border);border-radius:6px;background:var(--bg-card);color:var(--text);font-size:0.9rem;"
            @input="applyArticleFilter"
          />
          <select
            v-model="articleFilter.category"
            style="padding:8px 12px;border:1px solid var(--border);border-radius:6px;background:var(--bg-card);color:var(--text);font-size:0.9rem;"
            @change="applyArticleFilter"
          >
            <option value="">全部分类</option>
            <option value="blog">技术文章</option>
            <option value="leetcode">算法笔记</option>
            <option value="projects">项目</option>
            <option value="study">学习笔记</option>
            <option value="notes">碎碎念</option>
          </select>
          <select
            v-model="articleFilter.tag"
            style="padding:8px 12px;border:1px solid var(--border);border-radius:6px;background:var(--bg-card);color:var(--text);font-size:0.9rem;"
            @change="applyArticleFilter"
          >
            <option value="">全部标签</option>
            <option v-for="t in allTagsList" :key="t" :value="t">{{ t }}</option>
          </select>
          <span style="font-size:0.85rem;color:var(--text-muted);white-space:nowrap;">
            共 {{ filteredArticles.length }} 篇
          </span>
        </div>
        <div v-if="filteredArticles.length===0" style="text-align:center;color:var(--text-muted);padding:40px;">暂无文章</div>
        <div v-for="a in filteredArticles" :key="a.id" class="card" style="margin-bottom:10px;padding:14px 18px;">
          <div style="display:flex;justify-content:space-between;align-items:center;flex-wrap:wrap;gap:8px;">
            <div style="flex:1;min-width:0;">
              <div style="display:flex;align-items:center;gap:8px;margin-bottom:4px;">
                <span v-if="a.is_pinned" title="置顶" style="font-size:0.8rem;">📌</span><strong style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">{{ a.title }}</strong>
                <span class="tag" style="font-size:0.7rem;padding:1px 7px;">{{ catLabel(a.category) }}</span>
              </div>
              <div style="display:flex;flex-wrap:wrap;gap:4px;">
                <span v-for="t in parseTagsArr(a.tags)" :key="t" class="tag" style="font-size:0.7rem;padding:1px 7px;cursor:pointer;" @click="articleFilter.tag = t; applyArticleFilter()">{{ t }}</span>
                <span style="color:var(--text-muted);font-size:0.75rem;margin-left:8px;">{{ formatDate(a.created_at) }}</span>
              </div>
            </div>
            <div style="display:flex;gap:8px;flex-shrink:0;">
              <button class="btn btn-outline btn-sm" @click="editArticle(a)">编辑</button>
              <button class="btn btn-danger btn-sm" @click="deleteArticle(a.id)">删除</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Write / Edit -->
      <div v-if="tab==='write'">
        <div class="form-group">
          <label>标题</label>
          <input v-model="editorForm.title" placeholder="文章标题" />
        </div>
        <div class="form-group">
          <label>分类</label>
          <select v-model="editorForm.category">
            <option value="blog">技术文章</option>
            <option value="leetcode">算法笔记</option>
            <option value="projects">项目</option>
            <option value="study">学习笔记</option>
            <option value="notes">碎碎念</option>
          </select>
        </div>
        <div v-if="editorForm.category === 'study'" style="display:grid;grid-template-columns:1fr 1fr;gap:12px;">
          <div class="form-group">
            <label>子分类</label>
            <select v-model="editorForm.subcategory">
              <option value="">未分类</option>
              <option v-for="sc in studyColumns" :key="sc" :value="sc">{{ sc }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>学习状态</label>
            <select v-model="editorForm.study_status">
              <option value="">学习中</option>
              <option value="todo">待开始</option>
              <option value="doing">学习中</option>
              <option value="done">已完成</option>
            </select>
          </div>
        </div>
        <div class="form-group">
          <label>摘要</label>
          <textarea v-model="editorForm.summary" placeholder="简短摘要（可选）" rows="2"></textarea>
        </div>
        <div style="display:flex;align-items:center;gap:8px;margin-bottom:16px;">
          <label style="display:flex;align-items:center;gap:6px;cursor:pointer;">
            <input type="checkbox" v-model="editorForm.is_pinned" style="width:auto;accent-color:var(--accent);" />
            <span style="font-size:0.85rem;color:var(--text-muted);user-select:none;">📌 置顶文章</span>
          </label>
          <span v-if="editorForm.is_pinned" style="font-size:0.75rem;color:var(--accent);">（将显示在各分类列表最前）</span>
        </div>
        <div class="form-group">
          <label>标签（逗号分隔）</label>
          <input v-model="editorForm.tags" placeholder="Go, 后端, 微服务" />
        </div>
        <div class="editor-wrap">
          <textarea v-model="editorForm.content" placeholder="Markdown 内容..."></textarea>
          <div class="preview-pane article-content" v-html="previewContent"></div>
        </div>
        <div style="margin-top:16px;">
          <button class="btn btn-primary" @click="saveArticle">{{ editingId ? '更新文章' : '发布文章' }}</button>
          <button class="btn btn-outline" style="margin-left:8px;" @click="tab='articles'">取消</button>
        </div>
      </div>

      <!-- Study Categories -->
      <div v-if="tab==='studyCats'" style="max-width:600px;">
        <div style="display:flex;gap:8px;margin-bottom:16px;">
          <input v-model="newCat.name" placeholder="分类名称" style="flex:1;" @keyup.enter="addStudyCat" />
          <input v-model="newCat.icon" placeholder="图标" style="width:80px;" @keyup.enter="addStudyCat" />
          <button class="btn btn-primary" @click="addStudyCat">添加</button>
        </div>
        <div v-if="studyCatList.length===0" style="text-align:center;color:var(--text-muted);padding:40px;">暂无子分类</div>
        <div v-for="(cat, i) in studyCatList" :key="cat.id" class="card" style="margin-bottom:8px;display:flex;justify-content:space-between;align-items:center;padding:12px 16px;">
          <div>
            <span v-if="cat.icon" style="margin-right:8px;">{{ cat.icon }}</span>
            <strong>{{ cat.name }}</strong>
          </div>
          <div style="display:flex;gap:4px;align-items:center;">
            <button class="btn btn-outline btn-sm" @click="moveCat(cat, -1)" :disabled="i===0">↑</button>
            <button class="btn btn-outline btn-sm" @click="moveCat(cat, 1)" :disabled="i===studyCatList.length-1">↓</button>
            <button class="btn btn-danger btn-sm" @click="deleteStudyCat(cat.id)">删除</button>
          </div>
        </div>
      </div>

      <!-- Change Password -->
      <div v-if="tab==='password'" style="max-width:400px;">
        <div class="card">
          <div class="form-group">
            <label>旧密码</label>
            <input v-model="pwForm.old_password" type="password" />
          </div>
          <div class="form-group">
            <label>新密码</label>
            <input v-model="pwForm.new_password" type="password" />
          </div>
          <div class="form-group">
            <label>确认新密码</label>
            <input v-model="pwForm.confirm" type="password" />
          </div>
          <button class="btn btn-primary" @click="changePassword">修改密码</button>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <Teleport to="body">
      <Transition name="toast-fade">
        <div v-if="toast" class="toast-overlay" @click.self="toast = null">
          <div :class="['toast-box', 'toast-' + toast.type]">
            <span class="toast-icon">{{ toast.type === 'success' ? '?' : toast.type === 'error' ? '?' : '!' }}</span>
            <span class="toast-msg">{{ toast.msg }}</span>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { adminAPI, studyCategoryAPI } from '../api'
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
import rust from 'highlight.js/lib/languages/rust'
import typescript from 'highlight.js/lib/languages/typescript'
import plaintext from 'highlight.js/lib/languages/plaintext'

hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('js', javascript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('py', python)
hljs.registerLanguage('go', go)
hljs.registerLanguage('golang', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('shell', bash)
hljs.registerLanguage('sh', bash)
hljs.registerLanguage('css', css)
hljs.registerLanguage('json', json)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('mysql', sql)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('yml', yaml)
hljs.registerLanguage('java', java)
hljs.registerLanguage('cpp', cpp)
hljs.registerLanguage('c', cpp)
hljs.registerLanguage('rust', rust)
hljs.registerLanguage('rs', rust)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('ts', typescript)
hljs.registerLanguage('text', plaintext)
hljs.registerLanguage('plaintext', plaintext)

marked.use({ async: false })

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

const router = useRouter()
const auth = useAuthStore()
const isLoggedIn = computed(() => auth.isLoggedIn)

const tab = ref('articles')
const articles = ref([])
const articleFilter = ref({ search: '', category: '', tag: '' })
const allTagsList = ref([])
const filteredArticles = computed(() => {
  let list = articles.value
  if (articleFilter.value.search) {
    const q = articleFilter.value.search.toLowerCase()
    list = list.filter(a =>
      (a.title || '').toLowerCase().includes(q) ||
      (a.content || '').toLowerCase().includes(q) ||
      (a.summary || '').toLowerCase().includes(q)
    )
  }
  if (articleFilter.value.category) {
    list = list.filter(a => a.category === articleFilter.value.category)
  }
  if (articleFilter.value.tag) {
    list = list.filter(a => parseTagsArr(a.tags).includes(articleFilter.value.tag))
  }
  return list
})
const editingId = ref(null)
const toast = ref(null)

const loginForm = ref({ username: '', password: '' })
const loginErr = ref('')
const loginLoading = ref(false)

const studyCatList = ref([])
const studyColumns = ref([])
const newCat = ref({ name: '', icon: '' })

const pwForm = ref({ old_password: '', new_password: '', confirm: '' });
const editorForm = ref({ title: '', content: '', summary: '', tags: '', category: 'blog', subcategory: '', study_status: '', is_pinned: false })

const previewContent = computed(() => {
  if (!editorForm.value.content) return '<span style="color:#888">预览...</span>'
  return marked(editorForm.value.content)
})

function showToast(msg, type = 'success') {
  toast.value = { msg, type }
  setTimeout(() => toast.value = null, 3000)
}

function formatDate(d) {
  return d ? new Date(d).toLocaleDateString('zh-CN') : ''
}

// Auth
async function doLogin() {
  loginLoading.value = true
  loginErr.value = ''
  try {
    await auth.login(loginForm.value.username, loginForm.value.password)
    fetchArticles()
    fetchStudyCats()
  } catch (e) {
    loginErr.value = e?.message || '登录失败'
  } finally {
    loginLoading.value = false
  }
}

// Articles
function catLabel(c) {
  const m = { blog: '技术文章', leetcode: '算法笔记', projects: '项目', study: '学习笔记', notes: '碎碎念' }
  return m[c] || c
}
function parseTagsArr(tags) {
  if (!tags) return []
  return tags.split(',').map(t => t.trim()).filter(Boolean)
}

async function fetchArticles() {
  try {
    const res = await adminAPI.listArticles({ page: 1, limit: 200 })
    articles.value = res.data.articles || res.data || []
    // Build all tags list
    const tagSet = new Set()
    articles.value.forEach(a => parseTagsArr(a.tags).forEach(t => tagSet.add(t)))
    allTagsList.value = [...tagSet].sort()
    // Reset filters
    articleFilter.value = { search: '', category: '', tag: '' }
  } catch (e) { console.error(e) }
}

function applyArticleFilter() {
  // Rebuild tag list from all articles (not filtered)
  const tagSet = new Set()
  articles.value.forEach(a => parseTagsArr(a.tags).forEach(t => tagSet.add(t)))
  allTagsList.value = [...tagSet].sort()
}

function openWrite() {
  tab.value = 'write'
  editingId.value = null
  editorForm.value = { title: '', content: '', summary: '', tags: '', category: 'blog', subcategory: '', study_status: '', is_pinned: false }
}

function editArticle(a) {
  tab.value = 'write'
  editingId.value = a.id
  editorForm.value = { title: a.title, content: a.content, summary: a.summary, tags: a.tags, category: a.category || 'blog', subcategory: a.subcategory || '', study_status: a.study_status || '', is_pinned: a.is_pinned || false }
}

async function saveArticle() {
  try {
    if (editingId.value) {
      await adminAPI.updateArticle(editingId.value, editorForm.value)
      showToast('文章已更新')
    } else {
      await adminAPI.createArticle(editorForm.value)
      showToast('文章已发布')
    }
    tab.value = 'articles'
    fetchArticles()
  } catch (e) {
    showToast(e?.message || '操作失败', 'error')
  }
}

async function deleteArticle(id) {
  if (!confirm('确定删除这篇文章吗？')) return
  try {
    await adminAPI.deleteArticle(id)
    showToast('已删除')
    fetchArticles()
  } catch (e) {
    showToast(e?.message || '删除失败', 'error')
  }
}

async function changePassword() {
  if (!pwForm.value.old_password || !pwForm.value.new_password) {
    showToast('请填写完整', 'error'); return
  }
  if (pwForm.value.new_password.length < 6) {
    showToast('新密码至少6位', 'error'); return
  }
  if (pwForm.value.new_password !== pwForm.value.confirm) {
    showToast('两次密码不一致', 'error'); return
  }
  try {
    await adminAPI.changePassword({
      old_password: pwForm.value.old_password,
      new_password: pwForm.value.new_password,
    })
    showToast('密码修改成功')
    pwForm.value = { old_password: '', new_password: '', confirm: '' }
    tab.value = 'articles'
  } catch(e) {
    showToast(e?.message || '修改失败', 'error')
  }
}

async function fetchStudyCats() {
  try {
    const res = await studyCategoryAPI.list();
    const cats = res.data.categories || res.data || [];
    studyCatList.value = cats;
    studyColumns.value = cats.map(c => c.name);
    console.log('fetchStudyCats result:', studyCatList.value, studyColumns.value);
  } catch (e) { console.error(e) }
}
async function addStudyCat() {
  if (!newCat.value.name) return
  try {
    await studyCategoryAPI.create({ name: newCat.value.name, icon: newCat.value.icon, sort_order: studyCatList.value.length })
    showToast('子分类已添加')
    newCat.value = { name: '', icon: '' }
    fetchStudyCats()
  } catch (e) { showToast(e?.message || '添加失败', 'error') }
}
async function deleteStudyCat(id) {
  if (!confirm('确定删除此子分类？')) return
  try { await studyCategoryAPI.delete(id); showToast('已删除'); fetchStudyCats() } catch (e) { showToast(e?.message || '删除失败', 'error') }
}
async function moveCat(cat, dir) {
  const idx = studyCatList.value.findIndex(c => c.id === cat.id)
  if (idx === -1) return
  const other = studyCatList.value[idx + dir]
  if (!other) return
  try {
    await studyCategoryAPI.update(cat.id, { ...cat, sort_order: other.sort_order })
    await studyCategoryAPI.update(other.id, { ...other, sort_order: cat.sort_order })
    fetchStudyCats()
  } catch (e) { showToast('排序失败', 'error') }
}

watch(tab, (val) => {
  if (val === 'studyCats') fetchStudyCats()
})

onMounted(() => {
  fetchArticles()
  fetchStudyCats()
})
</script>
