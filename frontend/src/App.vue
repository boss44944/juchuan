<template>
  <el-container class="layout">
    <el-aside width="220px">
      <el-menu router>
        <el-menu-item index="/devices">{{ t('menu.devices') }}</el-menu-item>
        <el-menu-item index="/messages">{{ t('menu.messages') }}</el-menu-item>
        <el-menu-item index="/send">{{ t('menu.send') }}</el-menu-item>
        <el-menu-item index="/config">{{ t('menu.config') }}</el-menu-item>
      </el-menu>
    </el-aside>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from './stores/device'
import { useMessageStore } from './stores/message'
import { heartbeatDevice, registerDevice } from './api'
import { connectWebSocket, onWebSocketMessage } from './websocket/client'

const route = useRoute()
const { t } = useI18n()
const deviceStore = useDeviceStore()
const messageStore = useMessageStore()
let wsBound = false
let heartbeatTimer: number | null = null

watch(
  () => route.path,
  async (path) => {
    if (path === '/login') {
      return
    }

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
        heartbeatTimer = window.setInterval(() => {
          void heartbeatDevice(id)
        }, 30000)
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
</script>

<style scoped>
.layout {
  min-height: 100vh;
}
</style>
