<template>
  <div class="app-shell">
    <!-- Top bar -->
    <header class="top-bar">
      <router-link to="/" class="top-logo">Zephyr</router-link>
      <div class="top-right">
        <a href="https://github.com/ZephyrZheng1001" target="_blank" class="social-icon" title="GitHub">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.3 3.438 9.8 8.205 11.387.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.73.083-.73 1.205.085 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.418-1.305.762-1.604-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 21.795 24 17.295 24 12c0-6.63-5.37-12-12-12z"/></svg>
        </a>
        <a href="tencent://message/?uin=2531391463" class="social-icon" title="QQ">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M21.395 15.035a40 40 0 0 0-.803-2.264l-1.079-2.695c.001-.032.014-.562.014-.836C19.526 4.632 17.351 0 12 0S4.474 4.632 4.474 9.241c0 .274.013.804.014.836l-1.08 2.695a39 39 0 0 0-.802 2.264c-1.021 3.283-1.045 4.643 4.493 5.308.654.08 1.092.442 1.521.999.636.823 1.259 2.031 3.38 2.657 2.122-.626 2.745-1.834 3.38-2.657.429-.557.867-.919 1.521-.999 5.538-.665 5.514-2.025 4.493-5.308z"/></svg>
        </a>
        <button class="theme-toggle" @click="toggleTheme" :title="isDark ? '切换亮色' : '切换暗色'">
          {{ isDark ? '☀️' : '🌙' }}
        </button>
      </div>
    </header>

    <div class="shell-body">
      <!-- Sidebar -->
      <aside class="sidebar">
        <nav class="sidebar-nav">
          <router-link to="/" class="sidebar-item" exact-active-class="active"><span>🏠</span> 首页</router-link>

          <div class="sidebar-group">
            <span class="sidebar-group-title">文章</span>
            <router-link to="/blog" class="sidebar-item sidebar-sub"><span>📝</span> 技术文章</router-link>
            <router-link to="/leetcode" class="sidebar-item sidebar-sub"><span>💡</span> 算法笔记</router-link>
            <router-link to="/projects" class="sidebar-item sidebar-sub"><span>🚀</span> 项目</router-link>
            <router-link to="/notes" class="sidebar-item sidebar-sub"><span>💬</span> 碎碎念</router-link>
          </div>

          <router-link to="/music" class="sidebar-item"><span>🎧</span> 音乐</router-link>
          <router-link to="/about" class="sidebar-item"><span>👤</span> 关于</router-link>
        </nav>
      </aside>

      <!-- Main content -->
      <main class="main-content">
        <router-view />
      </main>
    </div>

    <footer class="site-footer">
      <p>
        © 2026 Zephyr &nbsp;·&nbsp;
        <router-link to="/admin">管理</router-link>
      </p>
    </footer>
  </div>
</template>

<script setup>
import { useTheme } from './composables/theme'
const { isDark, toggle: toggleTheme } = useTheme()
</script>

<style scoped>
.app-shell { min-height: 100vh; display: flex; flex-direction: column; }

/* Top bar */
.top-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 48px;
  background: rgba(255,255,255,0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  z-index: 100;
}
[data-theme="dark"] .top-bar { background: rgba(17,21,24,0.88); }
.top-logo {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text);
  letter-spacing: -0.3px;
}
.top-logo:hover { color: var(--accent); text-decoration: none; }
.top-right { display: flex; align-items: center; gap: 6px; }
.social-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 6px;
  color: var(--text-muted);
  transition: all 0.2s;
}
.social-icon:hover { background: var(--bg-hover); color: var(--accent); }
.theme-toggle {
  background: none; border: none;
  font-size: 1rem; cursor: pointer;
  padding: 4px 6px; border-radius: 6px;
  transition: background 0.2s; line-height: 1;
}
.theme-toggle:hover { background: var(--bg-hover); }

/* Shell body */
.shell-body {
  display: flex;
  margin-top: 48px;
  flex: 1;
}

/* Sidebar */
.sidebar {
  position: fixed;
  top: 48px;
  left: 0;
  bottom: 0;
  width: 220px;
  border-right: 1px solid var(--border);
  padding: 24px 16px;
  overflow-y: auto;
  background: var(--bg);
  z-index: 50;
}
.sidebar-nav { display: flex; flex-direction: column; gap: 2px; }
.sidebar-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 12px;
  border-radius: 6px;
  font-size: 0.88rem;
  color: var(--text-muted);
  transition: all 0.15s;
}
.sidebar-item:hover { background: var(--bg-hover); color: var(--text); text-decoration: none; }
.sidebar-item.active {
  background: var(--accent-light);
  color: var(--accent);
  font-weight: 600;
}
.sidebar-item span { font-size: 0.85rem; }
.sidebar-sub { padding-left: 28px; font-size: 0.84rem; }
.sidebar-group { margin: 8px 0 4px; }
.sidebar-group-title {
  display: block;
  padding: 8px 12px 4px;
  font-size: 0.72rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 1.2px;
}

/* Main content */
.main-content {
  margin-left: 220px;
  flex: 1;
  min-width: 0;
  padding-bottom: 60px;
}

/* Footer */
.site-footer {
  margin-left: 220px;
  text-align: center;
  padding: 32px 0;
  color: var(--text-muted);
  font-size: 0.8rem;
  border-top: 1px solid var(--border);
}
.site-footer a { color: var(--text-muted); }
.site-footer a:hover { color: var(--accent); }

@media (max-width: 768px) {
  .sidebar { display: none; }
  .main-content { margin-left: 0; }
  .site-footer { margin-left: 0; }
}
</style>
