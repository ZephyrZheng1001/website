<template>
  <div style="padding-top:32px;padding-bottom:80px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">
    <div style="display:flex;align-items:baseline;gap:12px;margin-bottom:8px;">
      <h1 style="font-size:2rem;">🎧 音乐</h1>
      <span style="color:var(--text-muted);font-size:0.9rem;">私藏歌单</span>
    </div>
    <p style="color:var(--text-muted);margin-bottom:40px;">好歌值得被听见。</p>

    <div v-if="loading" style="text-align:center;color:var(--text-muted);padding:60px;">加载中...</div>

    <div v-else-if="list.length === 0" style="text-align:center;color:var(--text-muted);padding:60px;">
      <p style="font-size:1.1rem;">还没有分享的音乐</p>
    </div>

    <div v-else class="music-grid">
      <a
        v-for="item in list"
        :key="item.id"
        :href="item.platform_url || '#'"
        target="_blank"
        class="card music-card-wrapper"
        style="display:block;"
      >
        <div class="music-card-row">
          <div :class="['cover', item.platform || 'netease']">
            <span v-if="item.platform==='qq'">🎵</span>
            <span v-else>🎶</span>
          </div>
          <div style="flex:1;min-width:0;">
            <h3 style="font-size:1.05rem;margin-bottom:2px;">{{ item.title }}</h3>
            <p style="color:var(--text-muted);font-size:0.85rem;">{{ item.artist }} <span v-if="item.album">· {{ item.album }}</span></p>
            <p v-if="item.description" style="color:var(--text-muted);font-size:0.82rem;margin-top:4px;line-height:1.4;">{{ item.description }}</p>
          </div>
          <div style="display:flex;align-items:center;gap:8px;flex-shrink:0;">
            <span v-if="item.platform==='qq'" style="font-size:0.75rem;color:var(--text-muted);background:#e8f5e9;padding:2px 8px;border-radius:10px;">QQ音乐</span>
            <span v-else-if="item.platform==='netease'" style="font-size:0.75rem;color:var(--text-muted);background:#fce4ec;padding:2px 8px;border-radius:10px;">网易云</span>
            <span class="music-link">▶ 去听</span>
          </div>
        </div>
      </a>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { musicAPI } from '../api'

const list = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await musicAPI.list()
    list.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.music-card-wrapper {
  padding: 0;
  overflow: hidden;
  cursor: pointer;
}
.music-card-wrapper:hover { text-decoration: none; }
.music-card-row {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 20px 24px;
}
.cover {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  flex-shrink: 0;
}
.cover.qq {
  background: linear-gradient(135deg, #31c27c 0%, #14b853 100%);
  box-shadow: 0 4px 16px rgba(49,194,124,0.25);
}
.cover.netease {
  background: linear-gradient(135deg, #c62f2f 0%, #a82828 100%);
  box-shadow: 0 4px 16px rgba(198,47,47,0.25);
}
.music-link {
  color: var(--accent);
  font-size: 0.85rem;
  font-weight: 500;
}
</style>
