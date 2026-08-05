<template>
  <div class="client-shell">
    <header class="client-header">
      <router-link to="/client/inbox" class="client-brand" :aria-label="t('client.home')">
        <img src="/app-logo.png" alt="" /><span><strong>菊传</strong><small>JUCHUAN</small></span>
      </router-link>
      <div class="client-header__actions">
        <Select
          v-model="language"
          :options="languageOptions"
          size="sm"
          :aria-label="t('configPage.labels.language')"
          class="w-24"
          @update:model-value="changeLanguage"
        />
        <button type="button" class="icon-action" :aria-label="t('menu.logout')" @click="handleLogout"><LogOut :size="20" /></button>
      </div>
    </header>

    <main class="client-content"><router-view /></main>

    <nav class="client-nav" :aria-label="t('client.navigation')">
      <router-link to="/client/inbox" class="client-nav__item">
        <span class="client-nav__icon"><Inbox :size="23" /><b v-if="unreadCount">{{ unreadLabel }}</b></span>
        <span>{{ t('client.inbox.title') }}</span>
      </router-link>
      <router-link to="/client/send" class="client-nav__item"><Send :size="23" /><span>{{ t('menu.send') }}</span></router-link>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Inbox, LogOut, Send } from '@lucide/vue'
import { Select } from '@/components/ui/select'
import { useMessageStore } from '@/stores/message'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { stopAppRuntime, useAppRuntime } from '@/composables/useAppRuntime'

const router = useRouter()
const { t, locale } = useI18n()
const messageStore = useMessageStore()
const authStore = useAuthStore()
const toast = useToast()
const language = ref(String(locale.value || 'zh-CN'))
const languageOptions = computed(() => [
  { label: t('configPage.languages.zhCN'), value: 'zh-CN' },
  { label: t('configPage.languages.enUS'), value: 'en-US' },
  { label: t('configPage.languages.jaJP'), value: 'ja-JP' },
])
const localID = (localStorage.getItem('device_id') || '').trim()
useAppRuntime('client')

const unreadCount = computed(() => messageStore.messages.filter((item) => item.target_device_id === localID && item.status !== 'READ').length)
const unreadLabel = computed(() => unreadCount.value > 99 ? '99+' : String(unreadCount.value))

function changeLanguage(value: string | undefined) {
  const v = value ?? 'zh-CN'
  locale.value = v
  localStorage.setItem('juchuan_locale', v)
}
async function handleLogout() {
  try { await authStore.signOut(); stopAppRuntime(); await router.replace('/client/login') }
  catch { toast.error(t('error.UNKNOWN')) }
}
</script>

<style scoped>
.client-shell { min-height: 100dvh; padding: 0 0 calc(72px + env(safe-area-inset-bottom)); background: var(--brutal-bg); }
.client-header { position: sticky; z-index: 30; top: 0; display: flex; align-items: center; justify-content: space-between; min-height: 56px; padding: calc(8px + env(safe-area-inset-top)) 14px 8px; border-bottom: 2px solid var(--brutal-border-color); background: rgba(255, 248, 231, .98); }
.client-brand { display: flex; align-items: center; gap: 9px; color: var(--brutal-fg); text-decoration: none; }
.client-brand img { width: 36px; height: 36px; object-fit: cover; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
.client-brand span { display: grid; }
.client-brand strong { font-size: 18px; line-height: 1.05; }
.client-brand small { margin-top: 2px; color: #9a4b1e; font-size: 8px; font-weight: 600; letter-spacing: .14em; }
.client-header__actions { display: flex; align-items: center; gap: 7px; }
.client-language, .icon-action { display: flex; align-items: center; min-height: 38px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: #fff; box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
.client-language { gap: 3px; padding: 0 7px; }
.client-language select { width: 42px; border: 0; background: transparent; color: var(--brutal-fg); font-weight: 600; }
.icon-action { width: 38px; justify-content: center; color: var(--brutal-fg); cursor: pointer; }
.client-language:focus-within, .icon-action:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 2px; }
.client-content { width: min(100%, 680px); margin: 0 auto; padding: 16px 14px 20px; }
.client-nav { position: fixed; z-index: 40; right: 0; bottom: 0; left: 0; display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); padding: 6px 10px calc(6px + env(safe-area-inset-bottom)); border-top: 2px solid var(--brutal-border-color); background: var(--brutal-primary); }
.client-nav__item { position: relative; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 3px; min-height: 52px; border: 2px solid transparent; border-radius: var(--brutal-radius); color: var(--brutal-fg); font-size: 11px; font-weight: 600; text-decoration: none; }
.client-nav__item.router-link-active { border-color: var(--brutal-border-color); background: var(--brutal-bg); box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
.client-nav__item:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: -3px; }
.client-nav__icon { position: relative; display: flex; }
.client-nav__icon b { position: absolute; top: -8px; left: 16px; min-width: 19px; padding: 1px 4px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: var(--brutal-secondary); color: #fff; font-size: 9px; text-align: center; }
@media (min-width: 700px) {
  .client-shell { width: 500px; min-height: calc(100vh - 24px); margin: 12px auto; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); box-shadow: 4px 4px 0 var(--brutal-shadow-color); overflow: hidden; }
  .client-nav { right: auto; left: 50%; width: 494px; transform: translateX(-50%); }
}
</style>
