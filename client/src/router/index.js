import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  { path: '/elephant', name: 'Elephant', component: { template: '<div></div>', beforeMount() { window.location.href = '/elephant.html' } } },
  { path: '/', name: 'Home', component: () => import('../views/Home.vue') },
  { path: '/blog', name: 'Blog', component: () => import('../views/Blog.vue') },
  { path: '/blog/:id', name: 'Article', component: () => import('../views/Article.vue') },
  { path: '/leetcode', name: 'Leetcode', component: () => import('../views/Blog.vue') },
  { path: '/leetcode/:id', name: 'LeetcodeArticle', component: () => import('../views/Article.vue') },
  { path: '/projects', name: 'Projects', component: () => import('../views/Blog.vue') },
  { path: '/projects/:id', name: 'ProjectArticle', component: () => import('../views/Article.vue') },
  { path: '/notes', name: 'Notes', component: () => import('../views/Blog.vue') },
  { path: '/notes/:id', name: 'NoteArticle', component: () => import('../views/Article.vue') },
  { path: '/about', name: 'About', component: () => import('../views/About.vue') },
  { path: '/admin', name: 'Admin', component: () => import('../views/Admin.vue') },
]

const router = createRouter({ history: createWebHashHistory(), routes })
export default router
