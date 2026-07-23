import { defineStore } from 'pinia'
import { ref } from 'vue'
import { adminAPI } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const isLoggedIn = ref(!!token.value)

  async function login(username, password) {
    const res = await adminAPI.login({ username, password })
    token.value = res.data.token
    localStorage.setItem('token', token.value)
    isLoggedIn.value = true
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
    isLoggedIn.value = false
  }

  return { token, isLoggedIn, login, logout }
})
