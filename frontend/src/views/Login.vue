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
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { login } from '../api'

const { t } = useI18n()
const deviceId = ref('')
const password = ref('')

async function loginSubmit() {
  const res = await login({
    device_id: deviceId.value,
    password: password.value,
  })

  localStorage.setItem('device_id', deviceId.value)
  localStorage.setItem('session', JSON.stringify(res.data))
  location.href = '/devices'
}
</script>

<style scoped>
.login-card {
  width: 360px;
  margin: 120px auto;
}
</style>
