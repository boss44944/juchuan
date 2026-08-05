<template>
  <main class="client-login theme-juchuan">
    <Card padding="lg" class="client-login__card">
      <header class="client-login__header">
        <img class="client-login__logo" src="/app-logo.png" alt="Juchuan Logo" />
        <div>
          <p>JUCHUAN / MOBILE</p>
          <h1>菊传</h1>
        </div>
      </header>

      <form class="client-login__form" @submit.prevent="loginSubmit">
        <h2>{{ t('login.title') }}</h2>

        <label class="field-label">
          {{ t('configPage.labels.language') }}
          <Select
            v-model="language"
            :options="languageOptions"
            :aria-label="t('configPage.labels.language')"
            class="w-full"
            @update:model-value="(value) => changeLanguage(value || 'zh-CN')"
          />
        </label>

        <label class="field-label">
          {{ t('login.deviceName') }}
          <Input
            v-model="deviceId"
            :placeholder="t('login.deviceName')"
            :aria-label="t('login.deviceName')"
            :aria-invalid="Boolean(error)"
            :aria-describedby="error ? 'client-login-error' : undefined"
          />
        </label>

        <label class="field-label">
          {{ t('login.password') }}
          <Input
            v-model="password"
            type="password"
            show-password
            :placeholder="t('login.password')"
            :aria-label="t('login.password')"
            :aria-invalid="Boolean(error)"
            :aria-describedby="error ? 'client-login-error' : undefined"
          />
        </label>

        <div v-if="error" id="client-login-error" class="client-login__alert" role="alert">
          <CircleAlert :size="20" :stroke-width="3" aria-hidden="true" />
          {{ error }}
        </div>

        <Button
          type="submit"
          variant="primary"
          size="lg"
          class="client-login__submit"
          :loading="submitting"
          :disabled="submitting"
        >
          <LogIn :size="20" aria-hidden="true" />
          {{ t('login.submit') }}
        </Button>
      </form>
    </Card>
  </main>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { CircleAlert, LogIn } from '@lucide/vue'
import { Button } from 'brutx-ui-vue/button'
import { Card } from 'brutx-ui-vue/card'
import { Input } from 'brutx-ui-vue/input'
import { Select } from 'brutx-ui-vue/select'
import { registerDevice, resolveApiErrorMessage } from '../../api'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const authStore = useAuthStore()

const deviceId = ref('')
const password = ref('')
const error = ref('')
const submitting = ref(false)
const language = ref(String(locale.value || 'zh-CN'))
const languageOptions = computed(() => [
  { label: t('configPage.languages.zhCN'), value: 'zh-CN' },
  { label: t('configPage.languages.enUS'), value: 'en-US' },
  { label: t('configPage.languages.jaJP'), value: 'ja-JP' },
])

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
    submitting.value = true
    await authStore.signIn(id, password.value)
    await registerDevice({
      id,
      display_name: id,
      role: 'client',
      platform: navigator.platform,
      browser: navigator.userAgent,
    })

    const requested = typeof route.query.redirect === 'string' ? route.query.redirect : ''
    const safeRedirect = requested.startsWith('/client/') ? requested : ''
    await router.replace(safeRedirect || '/client/inbox')
  } catch (err) {
    error.value = resolveApiErrorMessage(err, 'login.failed')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.client-login {
  display: grid;
  min-height: 100dvh;
  place-items: center;
  padding: max(16px, env(safe-area-inset-top)) max(16px, env(safe-area-inset-right)) max(16px, env(safe-area-inset-bottom)) max(16px, env(safe-area-inset-left));
}

.client-login__card {
  width: min(100%, 440px);
  overflow: hidden;
}

.client-login__header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 2px solid var(--brutal-border-color);
}

.client-login__logo {
  width: 48px;
  height: 48px;
  object-fit: cover;
  border: 2px solid var(--brutal-border-color);
  border-radius: var(--brutal-radius);
  box-shadow: 3px 3px 0 var(--brutal-shadow-color);
}

.client-login__header p {
  margin: 0 0 3px;
  color: #9a4b1e;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: .16em;
}

.client-login__header h1 {
  margin: 0;
  font-size: 22px;
  line-height: 1.2;
}

.client-login__form {
  display: grid;
  gap: 16px;
  padding-top: 20px;
}

.client-login__form h2 {
  margin: 0;
  font-size: 24px;
}

.client-login__alert {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 10px 11px;
  border: 2px solid var(--brutal-border-color);
  border-radius: var(--brutal-radius);
  background: #f2b3aa;
  font-weight: 600;
}

.client-login__submit {
  width: 100%;
  min-height: 48px;
  margin-top: 4px;
}

@media (max-width: 480px) {
  .client-login {
    place-items: start center;
    padding-top: max(24px, env(safe-area-inset-top));
  }

  .client-login__card {
    padding: 16px;
  }
}
</style>
