import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { login as apiLogin, sendSmsCode as apiSendSms, logout as apiLogout } from '@/api/user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(getToken())
  const user = ref<{ user_id: number; phone: string; nickname: string; avatar: string } | null>(null)
  const isLoggedIn = computed(() => !!token.value)

  async function sendSmsCode(phone: string) {
    await apiSendSms(phone)
  }

  async function login(phone: string, code: string) {
    const res = await apiLogin(phone, code)
    token.value = res.token
    setToken(res.token)
    user.value = {
      user_id: res.user.id,
      phone: res.user.phone,
      nickname: res.user.nickname,
      avatar: res.user.avatar,
    }
    return res
  }

  async function logout() {
    try { await apiLogout() } catch { /* ignore */ }
    token.value = null
    user.value = null
    removeToken()
  }

  return { token, user, isLoggedIn, sendSmsCode, login, logout }
})
