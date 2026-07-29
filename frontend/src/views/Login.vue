<template>
  <el-card class="login-card">
    <h2>{{ t('login.title') }}</h2>
    <el-form>
      <el-form-item>
        <el-input v-model="deviceName" :placeholder="t('login.deviceName')" />
      </el-form-item>
      <el-form-item>
        <el-input v-model="password" type="password" :placeholder="t('login.password')" />
      </el-form-item>
      <el-button type="primary" @click="login">{{ t('login.submit') }}</el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const deviceName = ref('')
const password = ref('')

function login() {
  localStorage.setItem('device_name', deviceName.value)
}
</script>

<style scoped>
.login-card {
  width: 360px;
  margin: 120px auto;
}
</style>
