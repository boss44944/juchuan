import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authStatus, login, logout } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const authenticated = ref(false)
  const deviceId = ref('')

  async function check() {
    const res = await authStatus()
    authenticated.value = !res.data.requires_password || res.data.authenticated
    return authenticated.value
  }

  async function signIn(id: string, password?: string) {
    await login({
      device_id: id,
      password
    })
    deviceId.value = id
    authenticated.value = true
  }

  async function signOut() {
    await logout()
    authenticated.value = false
    deviceId.value = ''
  }

  return { authenticated, deviceId, check, signIn, signOut }
})
