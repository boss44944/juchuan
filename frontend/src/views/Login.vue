<template>
  <el-card class="login-card">
    <h2>{{ t('login.title') }}</h2>
    <el-form>
      <el-form-item>
        <el-input v-model="deviceId" :placeholder="t('login.deviceName')" />
      </el-form-item>
      <el-form-item>
        <el-input v-model="password" type="password" :placeholder="t('login.password')" />
      </el-form-item>
      <el-button type="primary" @click="loginSubmit">{{ t('login.submit') }}</el-button>
      <el-alert v-if="error" :title="error" type="error" :closable="false" style="margin-top: 12px;" />
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { registerDevice, resolveApiErrorMessage } from '../api'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const deviceId = ref('')
const password = ref('')
const error = ref('')

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
.login-card {
  width: 360px;
  margin: 120px auto;
}
</style>
