<template>
  <div class="server-shell">
    <aside class="server-sidebar">
      <section class="brand-card">
        <img src="/app-logo.png" alt="Juchuan Logo" />
        <div><p>JUCHUAN</p><strong>菊传</strong></div>
      </section>
      <nav class="server-nav" :aria-label="t('menu.home')">
        <router-link v-for="item in navigation" :key="item.to" :to="item.to" class="server-nav__item">
          <component :is="item.icon" :size="21" :stroke-width="2.8" aria-hidden="true" />
          <span>{{ item.label }}</span>
        </router-link>
      </nav>
      <div class="local-note"><span />LOCAL ONLY</div>
    </aside>

    <main class="server-main">
      <header class="server-topbar">
        <div>
          <p>JUCHUAN / SERVER</p>
          <h1>{{ activeTitle }}</h1>
        </div>
        <div class="server-actions">
          <label class="select-field select-field--compact">
            <Languages :size="17" aria-hidden="true" />
            <span class="visually-hidden">Language</span>
            <select v-model="language" @change="changeLanguage(language)">
              <option value="zh-CN">中文</option><option value="en-US">English</option><option value="ja-JP">日本語</option>
            </select>
          </label>
          <Button variant="outline" size="sm" @click="handleLogout"><LogOut :size="17" />{{ t('menu.logout') }}</Button>
        </div>
      </header>
      <section class="server-content"><router-view /></section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Languages, LogOut, MessageSquareText, Monitor, Send, Settings } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { stopAppRuntime, useAppRuntime } from '@/composables/useAppRuntime'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const toast = useToast()
const language = ref(String(locale.value || 'zh-CN'))
useAppRuntime('server')

const navigation = computed(() => [
  { to: '/server/devices', label: t('menu.devices'), icon: Monitor },
  { to: '/server/messages', label: t('menu.messages'), icon: MessageSquareText },
  { to: '/server/send', label: t('menu.send'), icon: Send },
  { to: '/server/config', label: t('menu.config'), icon: Settings },
])
const activeTitle = computed(() => navigation.value.find((item) => route.path.startsWith(item.to))?.label || 'Juchuan')

function changeLanguage(value: string) {
  locale.value = value
  localStorage.setItem('juchuan_locale', value)
}
async function handleLogout() {
  try { await authStore.signOut(); stopAppRuntime(); await router.replace('/server/login') }
  catch { toast.error(t('error.UNKNOWN')) }
}
</script>

<style scoped>
.server-shell { display: grid; grid-template-columns: 240px minmax(0, 1fr); min-height: calc(100vh - 36px); margin: 18px; border: 3px solid var(--brutal-border-color); border-radius: 10px; background: var(--brutal-bg); box-shadow: 8px 8px 0 var(--brutal-shadow-color); overflow: hidden; }
.server-sidebar { display: flex; flex-direction: column; gap: 20px; padding: 18px; border-right: 3px solid var(--brutal-border-color); background: var(--brutal-primary); }
.brand-card { display: flex; align-items: center; gap: 12px; padding: 12px; border: 3px solid var(--brutal-border-color); border-radius: 7px; background: var(--brutal-bg); box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.brand-card img { width: 46px; height: 46px; object-fit: cover; border: 2px solid var(--brutal-border-color); border-radius: 5px; }
.brand-card p, .server-topbar p { margin: 0 0 3px; color: #9a4b1e; font-size: 10px; font-weight: 900; letter-spacing: .16em; }
.brand-card strong { font-size: 23px; }
.server-nav { display: grid; gap: 10px; }
.server-nav__item { display: flex; align-items: center; gap: 12px; min-height: 48px; padding: 10px 12px; border: 3px solid transparent; border-radius: 6px; color: var(--brutal-fg); font-weight: 900; text-decoration: none; transition: transform 140ms ease, box-shadow 140ms ease, background 140ms ease; }
.server-nav__item:hover, .server-nav__item.router-link-active { border-color: var(--brutal-border-color); background: var(--brutal-bg); box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.server-nav__item:hover { transform: translate(-2px, -2px); }
.server-nav__item:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 2px; }
.local-note { display: flex; align-items: center; gap: 8px; margin-top: auto; font-size: 10px; font-weight: 900; letter-spacing: .16em; }
.local-note span { width: 10px; height: 10px; border: 2px solid var(--brutal-border-color); border-radius: 50%; background: var(--brutal-success); }
.server-main { min-width: 0; padding: 22px; background: var(--brutal-bg); }
.server-topbar { display: flex; align-items: center; justify-content: space-between; gap: 18px; padding-bottom: 18px; border-bottom: 3px solid var(--brutal-border-color); }
.server-topbar h1 { margin: 0; font-size: clamp(28px, 3vw, 42px); line-height: 1; font-weight: 950; letter-spacing: -.035em; }
.server-actions { display: flex; align-items: center; gap: 12px; }
.server-content { padding-top: 24px; }
@media (max-width: 860px) {
  .server-shell { grid-template-columns: 190px minmax(0, 1fr); margin: 8px; min-height: calc(100vh - 16px); }
  .server-sidebar { padding: 12px; }
  .brand-card { padding: 9px; }
  .brand-card img { width: 38px; height: 38px; }
  .server-main { padding: 16px; }
  .server-topbar { align-items: flex-start; }
  .server-actions { flex-wrap: wrap; justify-content: flex-end; }
}
@media (prefers-reduced-motion: reduce) { .server-nav__item { transition: none; } }
</style>
