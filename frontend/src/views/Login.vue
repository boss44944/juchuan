<template>
  <section class="login-wrap theme-juchuan">
    <Card padding="lg" class="login-card">
      <header class="login-head">
        <img class="login-logo" src="/app-logo.png" alt="Juchuan Logo" />
        <div>
          <p>JUCHUAN / LOCAL ONLY</p>
          <h1>菊传</h1>
        </div>
      </header>

      <div class="login-grid" :class="{ 'login-grid--client': role === 'client' }">
        <section v-if="role === 'server'" class="qr-panel">
          <p>{{ t('login.qrTip') }}</p>
          <div class="login-qr-wrap">
            <img :src="qrImage" :alt="t('login.qrAlt')" class="login-qr" />
          </div>
        </section>

        <form class="login-form" @submit.prevent="loginSubmit">
          <h2>{{ t('login.title') }}</h2>
          <label class="field-label">
            {{ t('configPage.labels.language') }}
            <select v-model="language" class="brutal-select" @change="changeLanguage(language)">
              <option value="zh-CN">{{ t('configPage.languages.zhCN') }}</option>
              <option value="en-US">{{ t('configPage.languages.enUS') }}</option>
              <option value="ja-JP">{{ t('configPage.languages.jaJP') }}</option>
            </select>
          </label>
          <label class="field-label">
            {{ t('login.deviceName') }}
            <Input v-model="deviceId" size="lg" :placeholder="t('login.deviceName')" :aria-label="t('login.deviceName')" />
          </label>
          <label class="field-label">
            {{ t('login.password') }}
            <Input v-model="password" type="password" show-password size="lg" :placeholder="t('login.password')" :aria-label="t('login.password')" />
          </label>
          <div v-if="error" class="login-alert" role="alert">
            <CircleAlert :size="20" :stroke-width="3" aria-hidden="true" />
            {{ error }}
          </div>
          <Button type="submit" variant="primary" size="lg" class="login-submit">
            <LogIn :size="20" aria-hidden="true" />
            {{ t('login.submit') }}
          </Button>
        </form>
      </div>
    </Card>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { CircleAlert, LogIn } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useAuthStore } from '../stores/auth'
import { qrCodeURL, registerDevice, resolveApiErrorMessage } from '../api'
import { currentRole } from '../utils/role'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const deviceId = ref('')
const password = ref('')
const error = ref('')
const language = ref(String(locale.value || 'zh-CN'))
const role = computed(() => currentRole(route.path))
const entryURL = computed(() => `${window.location.origin}/client/inbox`)
const qrImage = computed(() => qrCodeURL(entryURL.value))

function changeLanguage(value: string) {
  locale.value = value
  language.value = value
  localStorage.setItem('juchuan_locale', value)
}

async function loginSubmit() {
  const id = deviceId.value.trim()
  if (!id) {
    error.value = t('login.deviceNameRequired')
    return
  }
  try {
    error.value = ''
    await authStore.signIn(id, password.value)
    await registerDevice({
      id,
      display_name: id,
      role: role.value,
      platform: navigator.platform,
      browser: navigator.userAgent,
    })
    const requested = typeof route.query.redirect === 'string' ? route.query.redirect : ''
    const safeRedirect = requested.startsWith(`/${role.value}/`) ? requested : ''
    await router.replace(safeRedirect || (role.value === 'server' ? '/server/devices' : '/client/inbox'))
  } catch (err) {
    error.value = resolveApiErrorMessage(err, 'login.failed')
  }
}
</script>

<style scoped>
.login-wrap { display: grid; min-height: 100dvh; place-items: center; padding: 18px; }
.login-card { width: min(860px, 94vw); }
.login-card { overflow: hidden; }
.login-head { display: flex; align-items: center; gap: 14px; padding-bottom: 18px; border-bottom: 3px solid var(--brutal-border-color); }
.login-logo { width: 64px; height: 64px; object-fit: cover; border: 3px solid var(--brutal-border-color); border-radius: 6px; box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.login-head p { margin: 0 0 3px; color: #9a4b1e; font-size: 10px; font-weight: 900; letter-spacing: .16em; }
.login-head h1 { margin: 0; font-size: 34px; line-height: 1; }
.login-grid { display: grid; grid-template-columns: .85fr 1.15fr; gap: 28px; padding-top: 24px; }
.login-grid--client { grid-template-columns: minmax(0, 460px); justify-content: center; }
.qr-panel { display: grid; align-content: start; gap: 12px; }
.qr-panel p { margin: 0; color: var(--brutal-muted-foreground); font-weight: 700; }
.login-qr-wrap { display: grid; place-items: center; padding: 18px; border: 3px solid var(--brutal-border-color); border-radius: 7px; background: #fff; box-shadow: 5px 5px 0 var(--brutal-shadow-color); }
.login-qr { width: min(240px, 100%); aspect-ratio: 1; background: #fff; padding: 8px; }
.login-form { display: grid; align-content: start; gap: 16px; }
.login-form h2 { margin: 0; font-size: 29px; }
.login-alert { display: flex; align-items: center; gap: 9px; padding: 11px 12px; border: 3px solid var(--brutal-border-color); border-radius: 6px; background: #f2b3aa; font-weight: 800; }
.login-submit { width: 100%; margin-top: 4px; }

@media (max-width: 680px) {
  .login-wrap { padding: 12px; }
  .login-card { width: 100%; padding: 18px; }
  .login-head { align-items: flex-start; }
  .login-logo { width: 52px; height: 52px; }
  .login-head h1 { font-size: 28px; }
  .login-grid { grid-template-columns: 1fr; gap: 22px; }
  .login-qr { width: 190px; }
  .login-submit { min-height: 48px; }
}

.diag-panel { margin-top: 14px; padding: 12px; border: 3px solid var(--brutal-border-color); border-radius: 7px; background: #0f172a; color: #e2e8f0; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 11px; }
.diag-title { margin: 0 0 8px; font-weight: 900; letter-spacing: .1em; color: #fbbf24; }
.diag-panel table { width: 100%; border-collapse: collapse; word-break: break-all; }
.diag-panel td { padding: 3px 4px; border-bottom: 1px solid #334155; vertical-align: top; }
.diag-panel td:first-child { color: #94a3b8; white-space: nowrap; padding-right: 8px; }
</style>
