<template>
  <div style="padding-top:32px;padding-bottom:80px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">
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
        <button class="btn" :class="tab==='music'?'btn-primary':'btn-outline'" @click="tab='music'">音乐管理</button>
        <button class="btn" :class="tab==='write'?'btn-primary':'btn-outline'" @click="openWrite()">写文章</button>
        <button class="btn" :class="tab==='password'?'btn-primary':'btn-outline'" @click="tab='password'">改密码</button>
      </div>

      <!-- Articles list -->
      <div v-if="tab==='articles'">
        <div v-if="articles.length===0" style="text-align:center;color:var(--text-muted);padding:40px;">暂无文章</div>
        <div v-for="a in articles" :key="a.id" class="card" style="margin-bottom:12px;display:flex;justify-content:space-between;align-items:center;">
          <div>
            <strong>{{ a.title }}</strong>
            <span style="color:var(--text-muted);font-size:0.8rem;margin-left:12px;">{{ formatDate(a.created_at) }}</span>
          </div>
          <div style="display:flex;gap:8px;">
            <button class="btn btn-outline btn-sm" @click="editArticle(a)">编辑</button>
            <button class="btn btn-danger btn-sm" @click="deleteArticle(a.id)">删除</button>
          </div>
        </div>
      </div>

      <!-- Music list -->
      <div v-if="tab==='music'">
        <div class="card" style="margin-bottom:20px;">
          <h3 style="margin-bottom:12px;">添加音乐</h3>
          <div class="form-group">
            <label>歌名</label>
            <input v-model="musicForm.title" placeholder="歌名" />
          </div>
          <div class="form-group">
            <label>艺术家</label>
            <input v-model="musicForm.artist" placeholder="艺术家" />
          </div>
          <div class="form-group">
            <label>专辑</label>
            <input v-model="musicForm.album" placeholder="专辑（可选）" />
          </div>
          <div class="form-group">
            <label>平台链接</label>
            <input v-model="musicForm.platform_url" placeholder="网易云/Spotify/QQ音乐链接" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="musicForm.description" placeholder="为什么喜欢这首歌？" rows="2"></textarea>
          </div>
          <button class="btn btn-primary" @click="addMusic">{{ editingMusicId ? '更新' : '添加' }}</button>
        </div>

        <div v-if="musicList.length===0" style="text-align:center;color:var(--text-muted);padding:20px;">暂无音乐</div>
        <div v-for="m in musicList" :key="m.id" class="card" style="margin-bottom:12px;display:flex;justify-content:space-between;align-items:center;">
          <div>
            <strong>{{ m.title }}</strong>
            <span style="color:var(--text-muted);margin-left:8px;">— {{ m.artist }}</span>
          </div>
          <button class="btn btn-outline btn-sm" @click="editMusic(m)">编辑</button>`n              <button class="btn btn-danger btn-sm" @click="deleteMusic(m.id)">删除</button>
        </div>
      </div>

      <!-- Write / Edit -->
      <div v-if="tab==='write'">
        <div class="form-group">
          <label>标题</label>
          <input v-model="editorForm.title" placeholder="文章标题" />
        </div>
        <div class="form-group">
          <label>摘要</label>
          <input v-model="editorForm.summary" placeholder="简短摘要（可选）" />
        </div>
        <div class="form-group"><label>专栏</label><select v-model="editorForm.category"><option value="blog">📝 技术文章</option><option value="leetcode">💡 算法笔记</option><option value="projects">🚀 项目</option><option value="notes">💬 碎碎念</option></select></div><div class="form-group">
          <label>标签（逗号分隔）</label>
          <input v-model="editorForm.tags" placeholder="如: Go,Vue,音乐" />
        </div>
        <div class="form-group">
          <label>内容（Markdown）</label>
          <div class="editor-wrap">
            <textarea v-model="editorForm.content" placeholder="用 Markdown 写文章..."></textarea>
            <div class="preview-pane article-content" v-html="previewContent"></div>
          </div>
        </div>
        <div style="display:flex;gap:12px;">
          <button class="btn btn-primary" @click="saveArticle">{{ editingId ? '更新' : '发布' }}</button>
          <button class="btn btn-outline" @click="tab='articles'">取消</button>
        </div>
      </div>

            <!-- Change Password -->
      <div v-if="tab==='password'" style="max-width:400px;">
        <div class="card">
          <h3 style="margin-bottom:16px;">修改密码</h3>
          <div class="form-group">
            <label>原密码</label>
            <input v-model="pwForm.old_password" type="password" placeholder="输入原密码" />
          </div>
          <div class="form-group">
            <label>新密码</label>
            <input v-model="pwForm.new_password" type="password" placeholder="至少6位" />
          </div>
          <div class="form-group">
            <label>确认新密码</label>
            <input v-model="pwForm.confirm" type="password" placeholder="再次输入新密码" />
          </div>
          <button class="btn btn-primary" @click="changePassword">确认修改</button>
        </div>
      </div>

      <!-- Toast -->
      <div v-if="toast" :class="['toast', 'toast-'+toast.type]">{{ toast.msg }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'
import { adminAPI } from '../api'
import { marked } from 'marked'

const router = useRouter()
const auth = useAuthStore()
const isLoggedIn = computed(() => auth.isLoggedIn)

const tab = ref('articles')
const articles = ref([])
const musicList = ref([])
const editingId = ref(null); const editingMusicId = ref(null)
const toast = ref(null)

const loginForm = ref({ username: '', password: '' })
const loginErr = ref('')
const loginLoading = ref(false)

const pwForm = ref({ old_password: '', new_password: '', confirm: '' });
const editorForm = ref({ title: '', content: '', summary: '', tags: '', category: 'blog' })
const musicForm = ref({ title: '', artist: '', album: '', platform_url: '', platform: 'netease', song_id: '', description: '' })

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
    fetchMusic()
  } catch (e) {
    loginErr.value = e?.message || '登录失败'
  } finally {
    loginLoading.value = false
  }
}

// Articles
async function fetchArticles() {
  try {
    const res = await adminAPI.listArticles({ page: 1, limit: 100 })
    articles.value = res.data.articles || []
  } catch (e) { console.error(e) }
}

function openWrite() {
  tab.value = 'write'
  editingId.value = null
  editorForm.value = { title: '', content: '', summary: '', tags: '', category: 'blog' }
}

function editArticle(a) {
  tab.value = 'write'
  editingId.value = a.id
  editorForm.value = { title: a.title, content: a.content, summary: a.summary, tags: a.tags, category: a.category || 'blog' }
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

// Music
async function fetchMusic() {
  try {
    const res = await adminAPI.listMusic()
    musicList.value = res.data || []
  } catch (e) { console.error(e) }
}

async function addMusic() {
  try {
    if (editingMusicId.value) {
      await adminAPI.updateMusic(editingMusicId.value, musicForm.value)
      showToast('音乐已更新')
      editingMusicId.value = null
    } else {
      await adminAPI.addMusic(musicForm.value)
      showToast('音乐已添加')
    }
    musicForm.value = { title: '', artist: '', album: '', platform_url: '', platform: 'netease', song_id: '', description: '' }
    fetchMusic()
  } catch (e) {
    showToast(e?.message || '添加失败', 'error')
  }
}

function editMusic(m) { editingMusicId.value = m.id; musicForm.value = { title: m.title, artist: m.artist, album: m.album, platform_url: m.platform_url, platform: m.platform || 'netease', song_id: m.song_id, description: m.description }; }

async function deleteMusic(id) {
  if (!confirm('确定删除吗？')) return
  try {
    await adminAPI.deleteMusic(id)
    showToast('已删除')
    fetchMusic()
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

onMounted(() => {
  if (isLoggedIn.value) {
    fetchArticles()
    fetchMusic()
  }
})
</script>
