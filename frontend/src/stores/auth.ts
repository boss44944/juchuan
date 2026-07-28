import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authStatus, login, logout } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const authenticated = ref(false)

  async function check() {
    const res = await authStatus()
    authenticated.value = !res.data.requires_password || res.data.authenticated
    return authenticated.value
  }

  async function signIn(password:string) {
    await login(password)
    authenticated.value = true
  }

  async function signOut() {
    await logout()
    authenticated.value = false
  }

  return { authenticated, check, signIn, signOut }
})
