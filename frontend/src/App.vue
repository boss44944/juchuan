<template>
  <div class="app-root theme-juchuan">
    <ToastHost />
    <main v-if="isLoginRoute" class="login-stage">
      <router-view />
    </main>

    <div v-else class="shell">
      <aside class="shell__aside">
        <section class="brand-card">
          <img class="brand-logo" src="/app-logo.png" alt="Juchuan Logo" />
          <div>
            <p>JUCHUAN</p>
            <strong>菊传</strong>
          </div>
        </section>

        <nav class="nav-menu" :aria-label="t('menu.home')">
          <router-link v-for="item in navigation" :key="item.to" :to="item.to" class="nav-item">
            <component :is="item.icon" :size="21" :stroke-width="2.8" aria-hidden="true" />
            <span>{{ item.label }}</span>
          </router-link>
        </nav>

        <div class="aside-note">
          <span class="aside-note__dot" />
          <span>LOCAL ONLY</span>
        </div>
      </aside>

      <main class="shell__main">
        <header class="topbar">
          <div>
            <p class="topbar__eyebrow">JUCHUAN / {{ route.path.slice(1).toUpperCase() }}</p>
            <h1>{{ activeTitle }}</h1>
          </div>
          <div class="topbar-actions">
            <label class="select-field select-field--compact">
              <Languages :size="17" aria-hidden="true" />
              <span class="visually-hidden">Language</span>
              <select v-model="language" @change="changeLanguage(language)">
                <option value="zh-CN">中文</option>
                <option value="en-US">English</option>
                <option value="ja-JP">日本語</option>
              </select>
            </label>
            <Button variant="outline" size="sm" @click="handleLogout">
              <LogOut :size="17" aria-hidden="true" />
              {{ t('menu.logout') }}
            </Button>
          </div>
        </header>
        <section class="content-stage">
          <router-view />
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Languages, LogOut, MessageSquareText, Monitor, Send, Settings } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import ToastHost from '@/components/app/ToastHost.vue'
import { useToast } from '@/composables/useToast'
import { useDeviceStore } from './stores/device'
import { useMessageStore } from './stores/message'
import { useAuthStore } from './stores/auth'
import { heartbeatDevice, registerDevice } from './api'
import { connectWebSocket, onWebSocketMessage } from './websocket/client'

const route = useRoute()
const { t, locale } = useI18n()
const toast = useToast()
const deviceStore = useDeviceStore()
const messageStore = useMessageStore()
const authStore = useAuthStore()
let wsBound = false
let heartbeatTimer: number | null = null
const language = ref(String(locale.value || 'zh-CN'))

const isLoginRoute = computed(() => route.path === '/login')
const navigation = computed(() => [
  { to: '/devices', label: t('menu.devices'), icon: Monitor },
  { to: '/messages', label: t('menu.messages'), icon: MessageSquareText },
  { to: '/send', label: t('menu.send'), icon: Send },
  { to: '/config', label: t('menu.config'), icon: Settings },
])

const activeTitle = computed(() => {
  if (route.path.startsWith('/devices')) return t('menu.devices')
  if (route.path.startsWith('/messages')) return t('menu.messages')
  if (route.path.startsWith('/send')) return t('menu.send')
  if (route.path.startsWith('/config')) return t('menu.config')
  return 'Juchuan'
})

watch(
  () => route.path,
  async (path) => {
    if (path === '/login') return
    await deviceStore.load()
    const id = (localStorage.getItem('device_id') || '').trim()
    if (id) {
      await registerDevice({
        id,
        display_name: id,
        role: 'client',
        platform: navigator.platform,
        browser: navigator.userAgent,
      })
      if (heartbeatTimer == null) {
        heartbeatTimer = window.setInterval(() => void heartbeatDevice(id), 30000)
      }
    }
    if (!wsBound) {
      onWebSocketMessage((event) => {
        deviceStore.handleEvent(event)
        messageStore.handleEvent(event)
      })
      wsBound = true
    }
    connectWebSocket()
  },
  { immediate: true }
)

function changeLanguage(value: string) {
  locale.value = value
  language.value = value
  localStorage.setItem('juchuan_locale', value)
}

async function handleLogout() {
  try {
    await authStore.signOut()
    if (heartbeatTimer != null) window.clearInterval(heartbeatTimer)
    heartbeatTimer = null
    window.location.href = '/login'
  } catch {
    toast.error(t('error.UNKNOWN'))
  }
}
</script>

<style scoped>
.app-root { min-height: 100vh; padding: 18px; }
.shell { display: grid; grid-template-columns: 240px minmax(0, 1fr); min-height: calc(100vh - 36px); border: 3px solid var(--brutal-border-color); border-radius: 10px; background: var(--brutal-bg); box-shadow: 8px 8px 0 var(--brutal-shadow-color); overflow: hidden; }
.shell__aside { display: flex; flex-direction: column; gap: 20px; padding: 18px; border-right: 3px solid var(--brutal-border-color); background: #f3b63f; }
.brand-card { display: flex; align-items: center; gap: 12px; padding: 12px; border: 3px solid var(--brutal-border-color); border-radius: 7px; background: #fff8e7; box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.brand-logo { width: 46px; height: 46px; object-fit: cover; border: 2px solid var(--brutal-border-color); border-radius: 5px; }
.brand-card p { margin: 0 0 2px; color: #9a4b1e; font-size: 10px; font-weight: 900; letter-spacing: .15em; }
.brand-card strong { font-size: 23px; line-height: 1; }
.nav-menu { display: grid; gap: 10px; }
.nav-item { display: flex; align-items: center; gap: 12px; min-height: 48px; padding: 10px 12px; border: 3px solid transparent; border-radius: 6px; color: var(--brutal-fg); font-weight: 900; text-decoration: none; transition: transform 140ms ease, box-shadow 140ms ease, background 140ms ease; }
.nav-item:hover { border-color: var(--brutal-border-color); background: #fff0bd; transform: translate(-2px, -2px); box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.nav-item.router-link-active { border-color: var(--brutal-border-color); background: #fff8e7; box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.nav-item:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 2px; }
.aside-note { display: flex; align-items: center; gap: 8px; margin-top: auto; font-size: 10px; font-weight: 900; letter-spacing: .16em; }
.aside-note__dot { width: 10px; height: 10px; border: 2px solid var(--brutal-border-color); border-radius: 50%; background: #86a95b; }
.shell__main { min-width: 0; padding: 22px; background: #fff8e7; }
.topbar { display: flex; align-items: center; justify-content: space-between; gap: 18px; padding-bottom: 18px; border-bottom: 3px solid var(--brutal-border-color); }
.topbar h1 { margin: 0; font-size: clamp(28px, 3vw, 42px); line-height: 1; font-weight: 950; letter-spacing: -.035em; }
.topbar__eyebrow { margin: 0 0 7px; color: #9a4b1e; font-size: 10px; font-weight: 900; letter-spacing: .16em; }
.topbar-actions { display: flex; align-items: center; gap: 12px; }
.content-stage { padding-top: 24px; }
.login-stage { display: grid; min-height: calc(100vh - 36px); place-items: center; }

@media (max-width: 780px) {
  .app-root { padding: 10px 10px calc(88px + env(safe-area-inset-bottom)); }
  .shell { display: block; min-height: calc(100vh - 20px); overflow: visible; }
  .shell__aside { position: fixed; z-index: 50; right: 10px; bottom: calc(10px + env(safe-area-inset-bottom)); left: 10px; display: block; padding: 8px; border: 3px solid var(--brutal-border-color); border-radius: 8px; box-shadow: 5px 5px 0 var(--brutal-shadow-color); }
  .brand-card, .aside-note { display: none; }
  .nav-menu { grid-template-columns: repeat(4, 1fr); gap: 5px; }
  .nav-item { flex-direction: column; gap: 2px; min-height: 52px; padding: 5px 2px; font-size: 10px; }
  .nav-item:hover { transform: none; box-shadow: none; }
  .nav-item.router-link-active { box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
  .shell__main { padding: 16px; }
  .topbar { align-items: flex-start; }
  .topbar-actions { flex-direction: column; align-items: stretch; }
  .topbar-actions :deep(button) { font-size: 0; min-width: 42px; padding-inline: 10px; }
  .content-stage { padding-top: 18px; }
}

@media (max-width: 440px) {
  .topbar__eyebrow { max-width: 180px; overflow: hidden; text-overflow: ellipsis; }
  .topbar h1 { font-size: 28px; }
}
</style>
