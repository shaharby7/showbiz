import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@showbiz/sdk'
import { useApi } from '@/composables/useApi'

export const useAuthStore = defineStore('auth', () => {
  const api = useApi()
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('accessToken'))

  const isAuthenticated = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const res = await api.auth.login({ email, password })
    token.value = res.accessToken
    localStorage.setItem('accessToken', res.accessToken)
    localStorage.setItem('refreshToken', res.refreshToken)
    user.value = res.user
  }

  async function register(email: string, password: string, displayName: string) {
    const res = await api.auth.register({ email, password, displayName })
    token.value = res.accessToken
    localStorage.setItem('accessToken', res.accessToken)
    localStorage.setItem('refreshToken', res.refreshToken)
    user.value = res.user
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
  }

  async function fetchCurrentUser() {
    user.value = await api.auth.me()
  }

  async function initialize() {
    const savedToken = localStorage.getItem('accessToken')
    if (savedToken) {
      token.value = savedToken
      try {
        await fetchCurrentUser()
      } catch {
        logout()
      }
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    register,
    logout,
    fetchCurrentUser,
    initialize,
  }
})
