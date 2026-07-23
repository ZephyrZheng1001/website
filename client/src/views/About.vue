<template>
  <div style="padding-top:48px;padding-bottom:80px;max-width:720px;margin:0 auto;padding-left:24px;padding-right:24px;">

    <!-- Header -->
    <section style="text-align:center;margin-bottom:48px;">
      <div class="avatar">Z</div>
      <h1 style="font-size:2rem;margin-top:16px;">Zephyr</h1>
      <p style="color:var(--text-muted);margin-top:6px;">开发者 · 音乐爱好者</p>
    </section>

    <!-- Bio -->
    <section class="card" style="margin-bottom:24px;text-align:center;padding:32px;">
      <p style="font-size:1.05rem;line-height:1.8;color:var(--text-muted);">
        喜欢干净代码和好旋律。<br/>
        用 Go 构建后端服务，用音乐填满生活。<br/>
        永远在学习，永远在探索。
      </p>
    </section>

    <!-- Contact -->
    <section style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px;margin-bottom:24px;">
      <a :href="'mailto:2453374@tongji.edu.cn'" class="card contact-card">
        <div class="contact-icon">📧</div>
        <div>
          <h3>学校邮箱</h3>
          <p>2453374@tongji.edu.cn</p>
        </div>
      </a>
      <a href="https://github.com/ZephyrZheng1001" target="_blank" class="card contact-card">
        <div class="contact-icon">🐙</div>
        <div>
          <h3>GitHub</h3>
          <p>ZephyrZheng1001</p>
        </div>
      </a>
    </section>

    <section style="display:grid;grid-template-columns:repeat(auto-fit,minmax(240px,1fr));gap:14px;margin-bottom:24px;">
      <div class="card contact-card" style="cursor:default;">
        <div class="contact-icon">💬</div>
        <div>
          <h3>QQ</h3>
          <p>2531391463</p>
        </div>
      </div>
      <div class="card contact-card" style="cursor:default;">
        <div class="contact-icon">💚</div>
        <div>
          <h3>微信</h3>
          <p>ZephyrZheng1001</p>
        </div>
      </div>
    </section>

    <!-- Site info -->
    <section class="card" style="text-align:center;padding:24px;">
      <p style="color:var(--text-muted);font-size:0.88rem;">
        本站由 <strong>Go</strong> 驱动 &nbsp;·&nbsp; 部署在我自己的服务器上
      </p>
    </section>
  </div>
</template>

<style scoped>
.avatar {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, var(--accent) 0%, #5bb8a8 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  font-weight: 700;
  color: #fff;
  margin: 0 auto;
  box-shadow: 0 4px 20px rgba(45,138,123,0.3);
}
.contact-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  cursor: pointer;
}
.contact-card:hover { text-decoration: none; }
.contact-icon {
  font-size: 1.5rem;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent-light);
  border-radius: 12px;
  flex-shrink: 0;
}
.contact-card h3 {
  font-size: 0.9rem;
  color: var(--text);
  margin-bottom: 2px;
}
.contact-card p {
  font-size: 0.82rem;
  color: var(--text-muted);
}
</style>
