<template>
  <div class="app-shell">
    <!-- Top bar -->
    <header class="top-bar">
      <button class="hamburger-btn" @click="sidebarOpen = !sidebarOpen" aria-label="菜单">
        <span :class="{ open: sidebarOpen }"></span>
        <span :class="{ open: sidebarOpen }"></span>
        <span :class="{ open: sidebarOpen }"></span>
      </button>
      <router-link to="/" class="top-logo">Zephyr</router-link>
      <div class="top-right">
        <router-link to="/search" class="social-icon" title="搜索">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
        </router-link>
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

    <!-- Overlay for mobile sidebar -->
    <div class="sidebar-overlay" :class="{ show: sidebarOpen }" @click="sidebarOpen = false"></div>

    <div class="shell-body">
      <!-- Desktop Sidebar -->
      <aside class="sidebar" :class="{ open: sidebarOpen }">
        <nav class="sidebar-nav">
          <router-link to="/" class="sidebar-item" exact-active-class="active" @click="sidebarOpen = false"><span>🏠</span> 首页</router-link>

          <div class="sidebar-group">
            <span class="sidebar-group-title">文章</span>
            <router-link to="/blog" class="sidebar-item sidebar-sub" @click="sidebarOpen = false"><span>📝</span> 技术文章</router-link>
            <router-link to="/leetcode" class="sidebar-item sidebar-sub" @click="sidebarOpen = false"><span>💡</span> 算法笔记</router-link>
            <router-link to="/projects" class="sidebar-item sidebar-sub" @click="sidebarOpen = false"><span>🚀</span> 项目</router-link>
            <router-link to="/notes" class="sidebar-item sidebar-sub" @click="sidebarOpen = false"><span>💬</span> 碎碎念</router-link>
          </div>

          <div class="sidebar-divider"></div>
          <router-link to="/study" class="sidebar-item" @click="sidebarOpen = false"><span>📖</span> 学习笔记</router-link>
          <router-link to="/about" class="sidebar-item" @click="sidebarOpen = false"><span>👤</span> 关于</router-link>
        </nav>
      </aside>

      <!-- Main content -->
      <main class="main-content">
        <router-view />
      </main>
      <button class="back-to-top" :class="{ visible: showBackTop }" @click="scrollToTop" title="返回顶部">↑</button>
    </div>

    <!-- Mobile Bottom Nav -->
    <nav class="bottom-nav">
      <router-link to="/" class="bottom-nav-item" exact-active-class="active">
        <span>🏠</span>
        <span>首页</span>
      </router-link>
      <router-link to="/blog" class="bottom-nav-item" active-class="active">
        <span>📝</span>
        <span>文章</span>
      </router-link>
      <router-link to="/study" class="bottom-nav-item" active-class="active">
        <span>📖</span>
        <span>学习</span>
      </router-link>
      <router-link to="/about" class="bottom-nav-item" active-class="active">
        <span>👤</span>
        <span>关于</span>
      </router-link>
      <button class="bottom-nav-item" @click="sidebarOpen = !sidebarOpen">
        <span>☰</span>
        <span>更多</span>
      </button>
    </nav>

    <footer class="site-footer">
      <p>
        © 2026 Zephyr &nbsp;·&nbsp;
        <router-link to="/admin">管理</router-link>
      </p>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useTheme } from './composables/theme'
const { isDark, toggle: toggleTheme } = useTheme()
const showBackTop = ref(false)
function onMainScroll() { showBackTop.value = window.scrollY > 400 }
function scrollToTop() { window.scrollTo({ top: 0, behavior: 'smooth' }) }
onMounted(() => window.addEventListener('scroll', onMainScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onMainScroll))
const sidebarOpen = ref(false)
</script>

<style scoped>
.app-shell { min-height: 100vh; display: flex; flex-direction: column; }

/* Top bar */
.top-bar {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 48px;
  background: rgba(255,255,255,0.92);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 200;
}
[data-theme="dark"] .top-bar { background: rgba(17,21,24,0.92); }
.top-logo {
  font-size: 1.1rem; font-weight: 700; color: var(--text); letter-spacing: -0.3px;
}
.top-logo:hover { color: var(--accent); text-decoration: none; }
.top-right { display: flex; align-items: center; gap: 2px; }
.social-icon {
  display: flex; align-items: center; justify-content: center;
  width: 36px; height: 36px; border-radius: 8px;
  color: var(--text-muted); transition: all 0.2s;
}
.social-icon:hover { background: var(--bg-hover); color: var(--accent); }
.theme-toggle {
  background: none; border: none; font-size: 1.05rem; cursor: pointer;
  padding: 6px 8px; border-radius: 8px; transition: background 0.2s; line-height: 1;
}
.theme-toggle:hover { background: var(--bg-hover); }

/* Hamburger */
.hamburger-btn {
  display: none;
  flex-direction: column; gap: 4px;
  background: none; border: none; cursor: pointer;
  padding: 8px; border-radius: 8px;
  transition: background 0.2s;
}
.hamburger-btn:hover { background: var(--bg-hover); }
.hamburger-btn span {
  display: block; width: 18px; height: 2px;
  background: var(--text); border-radius: 2px;
  transition: all 0.25s;
}
.hamburger-btn span.open:nth-child(1) { transform: translateY(6px) rotate(45deg); }
.hamburger-btn span.open:nth-child(2) { opacity: 0; }
.hamburger-btn span.open:nth-child(3) { transform: translateY(-6px) rotate(-45deg); }

/* Sidebar overlay */
.sidebar-overlay {
  display: none;
  position: fixed; inset: 0; top: 48px;
  background: rgba(0,0,0,0.3);
  z-index: 150; opacity: 0; transition: opacity 0.3s;
  pointer-events: none;
}
.sidebar-overlay.show { opacity: 1; pointer-events: auto; }

/* Shell body */
.shell-body { display: flex; margin-top: 48px; flex: 1; }

/* Sidebar */
.sidebar {
  position: fixed; top: 48px; left: 0; bottom: 0;
  width: 240px;
  border-right: 1px solid var(--border);
  padding: 20px 12px;
  overflow-y: auto;
  background: var(--bg);
  z-index: 160;
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.sidebar-nav { display: flex; flex-direction: column; gap: 1px; }
.sidebar-item {
  display: flex; align-items: center; gap: 10px;
  padding: 10px 16px; border-radius: 8px;
  font-size: 0.9rem; color: var(--text-muted);
  transition: all 0.15s; min-height: 40px;
}
.sidebar-item:hover { background: var(--bg-hover); color: var(--text); text-decoration: none; }
.sidebar-item.active {
  background: var(--accent-light); color: var(--accent); font-weight: 600;
}
.sidebar-item span { font-size: 0.9rem; }
.sidebar-sub { padding-left: 32px; font-size: 0.86rem; }
.sidebar-group { margin: 8px 0 4px; }
.sidebar-group-title {
  display: block; padding: 10px 16px 6px;
  font-size: 0.7rem; font-weight: 700;
  color: var(--text-muted); text-transform: uppercase;
  letter-spacing: 1.4px;
}
.sidebar-divider { height: 1px; background: var(--border); margin: 12px 16px 8px; }

/* Main content */
.main-content { margin-left: 240px; flex: 1; min-width: 0; padding-bottom: 20px; }

/* Footer */
.site-footer {
  margin-left: 240px; text-align: center;
  padding: 28px 0; color: var(--text-muted);
  font-size: 0.78rem; border-top: 1px solid var(--border);
}
.site-footer a { color: var(--text-muted); }
.site-footer a:hover { color: var(--accent); }

/* Bottom nav (mobile) */
.bottom-nav {
  display: none;
  position: fixed; bottom: 0; left: 0; right: 0;
  height: 60px;
  background: rgba(255,255,255,0.95);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-top: 1px solid var(--border);
  z-index: 200;
  padding-bottom: env(safe-area-inset-bottom);
}
[data-theme="dark"] .bottom-nav { background: rgba(17,21,24,0.95); }
.bottom-nav-item {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  gap: 3px; flex: 1;
  background: none; border: none; cursor: pointer;
  color: var(--text-muted); font-size: 0.65rem;
  padding: 6px 4px; font-family: inherit;
  transition: color 0.15s; text-decoration: none;
  -webkit-tap-highlight-color: transparent;
}
.bottom-nav-item span:first-child { font-size: 1.15rem; }
.bottom-nav-item.active { color: var(--accent); font-weight: 600; }

/* ===== RESPONSIVE ===== */
@media (max-width: 768px) {
  .hamburger-btn { display: flex; }
  .sidebar-overlay { display: block; }
  .sidebar { transform: translateX(-100%); }
  .sidebar.open { transform: translateX(0); box-shadow: 4px 0 24px rgba(0,0,0,0.1); }
  .main-content { margin-left: 0; padding-bottom: 80px; }
  .site-footer { margin-left: 0; padding-bottom: 80px; }
  .bottom-nav { display: flex; }
  .top-bar { padding: 0 12px; }
  .top-logo { font-size: 1rem; }
  .social-icon { width: 32px; height: 32px; }
}

@media (min-width: 769px) {
  .sidebar { transform: translateX(0) !important; }
  .sidebar-overlay { display: none !important; }
}

.back-to-top {
  position: fixed; bottom: 80px; right: 24px;
  width: 42px; height: 42px; border-radius: 50%;
  background: var(--accent); color: #fff; border: none;
  font-size: 1.2rem; cursor: pointer; z-index: 300;
  opacity: 0; transform: translateY(16px);
  transition: opacity 0.3s, transform 0.3s;
  pointer-events: none; box-shadow: 0 4px 16px rgba(45,138,123,0.3);
}
.back-to-top.visible { opacity: 1; transform: translateY(0); pointer-events: auto; }
.back-to-top:hover { background: #236b60; transform: translateY(-2px); }
@media (max-width: 768px) {
  .back-to-top { bottom: 76px; right: 16px; width: 38px; height: 38px; font-size: 1rem; }
}
</style>