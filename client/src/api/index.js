import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
    }
    return Promise.reject(err.response?.data || err)
  }
)

// ---- Public APIs ----
export const articleAPI = {
  list: (params) => api.get('/articles', { params }),
  get: (id) => api.get(`/articles/${id}`),
}

// ---- Admin APIs (need auth) ----
// ---- Admin APIs (need auth) ----
export const studyCategoryAPI = {
  list: () => api.get('/study-categories'),
  create: (data) => api.post('/admin/study-categories', data),
  update: (id, data) => api.put(`/admin/study-categories/${id}`, data),
  delete: (id) => api.delete(`/admin/study-categories/${id}`),
}

export const adminAPI = {
  login: (data) => api.post('/admin/login', data),
  createArticle: (data) => api.post('/admin/articles', data),
  updateArticle: (id, data) => api.put(`/admin/articles/${id}`, data),
  deleteArticle: (id) => api.delete(`/admin/articles/${id}`),
  listArticles: (params) => api.get('/admin/articles', { params }),
  changePassword: (data) => api.put('/admin/password', data),
}

export default api
