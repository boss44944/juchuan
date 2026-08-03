<template>
  <section class="login-wrap">
    <el-card class="login-card" shadow="never">
      <div class="login-stack">
        <header class="login-head">
          <img class="login-logo" src="/app-logo.png" alt="Juchuan Logo" />
          <div class="login-brand-block" aria-label="juchuan 菊传">
            <div class="login-brand-line">
              <p class="login-kicker">JuChuan</p>
              <span class="login-brand-char">菊</span>
            </div>
            <div class="login-brand-line login-brand-line--second">
              <span class="login-brand-char">传</span>
            </div>
          </div>
        </header>

        <p class="login-sub">{{ t('login.qrTip') }}</p>

        <div class="login-qr-wrap">
          <img :src="qrImage" :alt="t('login.qrAlt')" class="login-qr" />
        </div>

        <el-select v-model="language" size="default" class="login-lang-switch" @change="changeLanguage">
          <el-option :label="t('configPage.languages.zhCN')" value="zh-CN" />
          <el-option :label="t('configPage.languages.enUS')" value="en-US" />
          <el-option :label="t('configPage.languages.jaJP')" value="ja-JP" />
        </el-select>

        <div class="login-form-zone">
          <h2>{{ t('login.title') }}</h2>
          <el-form @submit.prevent>
            <el-form-item>
              <el-input v-model="deviceId" :placeholder="t('login.deviceName')" />
            </el-form-item>
            <el-form-item>
              <el-input v-model="password" type="password" :placeholder="t('login.password')" show-password />
            </el-form-item>
            <el-button type="primary" class="login-submit" @click="loginSubmit">{{ t('login.submit') }}</el-button>
            <el-alert v-if="error" :title="error" type="error" :closable="false" style="margin-top: 12px;" />
          </el-form>
        </div>
      </div>
    </el-card>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { qrCodeURL, registerDevice, resolveApiErrorMessage } from '../api'

const router = useRouter()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const deviceId = ref('')
const password = ref('')
const error = ref('')
const language = ref(String(locale.value || 'zh-CN'))
const entryURL = computed(() => `${window.location.origin}/`)
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
      role: 'client',
      platform: navigator.platform,
      browser: navigator.userAgent,
    })
    await router.replace('/devices')
  } catch (err) {
    error.value = resolveApiErrorMessage(err, 'login.failed')
  }
}
</script>

<style scoped>
.login-wrap {
  width: min(420px, 92vw);
}

.login-card {
  overflow: hidden;
  border-radius: 16px;
}

.login-stack {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 20px;
}

.login-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.login-logo {
  width: 64px;
  height: 64px;
  object-fit: cover;
  border-radius: 10px;
  border: 1px solid rgba(255, 221, 162, 0.58);
  box-shadow: 0 0 0 4px rgba(239, 176, 78, 0.18);
}

.login-brand-block {
  margin-left: auto;
  display: inline-flex;
  flex-direction: column;
}

.login-brand-line {
  display: flex;
  align-items: baseline;
  justify-content: flex-end;
  gap: 12px;
}

.login-brand-line--second {
  margin-top: 2px;
}

.login-kicker {
  margin: 0;
  font-size: 30px;
  letter-spacing: 0.08em;
  color: #d7be8f;
  text-transform: lowercase;
}

.login-brand-char {
  display: inline-block;
  font-size: 32px;
  line-height: 1;
  letter-spacing: 0.02em;
  color: #fff2d7;
}

.login-sub {
  margin: 0;
  color: #e2cfaa;
  font-size: 13px;
  text-align: center;
}

.login-qr-wrap {
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  border: 1px solid rgba(239, 176, 78, 0.25);
  border-radius: 12px;
  background: rgba(20, 16, 12, 0.45);
}

.login-qr-wrap::before {
  content: '';
  position: absolute;
  inset: -36px;
  border-radius: 50%;
  pointer-events: none;
  background: repeating-conic-gradient(
    from 0deg,
    rgba(255, 212, 116, 0.28) 0deg 9deg,
    rgba(189, 104, 31, 0.1) 9deg 18deg
  );
  mask-image: radial-gradient(circle at center, transparent 39%, rgba(0, 0, 0, 0.96) 52%, transparent 64%);
  opacity: 0.8;
}

.login-qr-wrap::after {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at center, transparent 54%, rgba(255, 229, 169, 0.16) 60%, transparent 72%),
    repeating-linear-gradient(14deg, rgba(255, 211, 116, 0.06) 0 1px, transparent 1px 8px);
}

.login-qr {
  position: relative;
  z-index: 1;
  width: 220px;
  height: 220px;
  border-radius: 8px;
  border: 1px solid rgba(223, 169, 73, 0.5);
  background:
    radial-gradient(circle at 22% 20%, rgba(255, 248, 230, 0.9), transparent 44%),
    repeating-linear-gradient(16deg, rgba(222, 187, 121, 0.08) 0 1px, rgba(0, 0, 0, 0) 1px 7px),
    #fff5e3;
  box-shadow:
    0 6px 14px rgba(35, 19, 8, 0.26),
    inset 0 0 0 1px rgba(255, 236, 196, 0.7);
  padding: 8px;
}

.login-form-zone {
  padding-top: 2px;
}

.login-form-zone h2 {
  margin: 0 0 12px;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.login-lang-switch {
  width: 100%;
}

.login-submit {
  width: 100%;
  margin-top: 6px;
}

@media (max-width: 860px) {
  .login-wrap {
    width: min(100%, 430px);
  }

  .login-stack {
    padding: 16px;
    gap: 12px;
  }

  .login-logo {
    width: 56px;
    height: 56px;
  }

  .login-brand-line {
    gap: 10px;
  }

  .login-brand-char {
    font-size: 28px;
  }

  .login-qr {
    width: 180px;
    height: 180px;
  }
}
</style>
