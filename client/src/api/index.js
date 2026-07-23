import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
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

export const musicAPI = {
  list: () => api.get('/music'),
}

// ---- Admin APIs (need auth) ----
export const adminAPI = {
  login: (data) => api.post('/admin/login', data),
  createArticle: (data) => api.post('/admin/articles', data),
  updateArticle: (id, data) => api.put(`/admin/articles/${id}`, data),
  deleteArticle: (id) => api.delete(`/admin/articles/${id}`),
  listArticles: (params) => api.get('/admin/articles', { params }),
  addMusic: (data) => api.post('/admin/music', data),
  changePassword: (data) => api.put('/admin/password', data),
  updateMusic: (id, data) => api.put('/admin/music/' + id, data),
  deleteMusic: (id) => api.delete(`/admin/music/${id}`),
  listMusic: () => api.get('/admin/music'),
}

export default api
