<template>
  <div class="app-root">
    <div class="orb orb--left" />
    <div class="orb orb--right" />

    <main v-if="isLoginRoute" class="login-stage">
      <router-view />
    </main>

    <el-container v-else class="shell">
      <el-aside class="shell__aside" width="260px">
        <section class="brand-card">
          <div class="brand-head">
            <img class="brand-logo" src="/app-logo.png" alt="Juchuan Logo" />
            <div>
              <div class="brand-card__label">Juchuan</div>
              <h2 class="brand-card__title">菊传</h2>
            </div>
          </div>
          <p class="brand-card__hint">Secure Transfer Workspace</p>
        </section>

        <el-menu class="nav-menu" router :default-active="route.path">
          <el-menu-item index="/devices">
            <el-icon class="menu-icon"><Monitor /></el-icon>
            <span class="menu-label">{{ t('menu.devices') }}</span>
          </el-menu-item>
          <el-menu-item index="/messages">
            <el-icon class="menu-icon"><ChatLineRound /></el-icon>
            <span class="menu-label">{{ t('menu.messages') }}</span>
          </el-menu-item>
          <el-menu-item index="/send">
            <el-icon class="menu-icon"><Promotion /></el-icon>
            <span class="menu-label">{{ t('menu.send') }}</span>
          </el-menu-item>
          <el-menu-item index="/config">
            <el-icon class="menu-icon"><Setting /></el-icon>
            <span class="menu-label">{{ t('menu.config') }}</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-main class="shell__main">
        <header class="topbar">
          <div>
            <h1 class="topbar__title">{{ activeTitle }}</h1>
            <p class="topbar__path">{{ route.path }}</p>
          </div>
          <div class="topbar-actions">
            <el-select v-model="language" size="small" class="lang-switch" @change="changeLanguage">
              <el-option label="中文" value="zh-CN" />
              <el-option label="English" value="en-US" />
              <el-option label="日本語" value="ja-JP" />
            </el-select>
            <el-button size="small" @click="handleLogout">{{ t('menu.logout') }}</el-button>
            <div class="topbar__badge">Enterprise Console</div>
          </div>
        </header>
        <section class="content-stage">
          <router-view />
        </section>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ChatLineRound, Monitor, Promotion, Setting } from '@element-plus/icons-vue'
import { useDeviceStore } from './stores/device'
import { useMessageStore } from './stores/message'
import { useAuthStore } from './stores/auth'
import { heartbeatDevice, registerDevice } from './api'
import { connectWebSocket, onWebSocketMessage } from './websocket/client'
import { ElMessage } from 'element-plus'

const route = useRoute()
const { t, locale } = useI18n()
const deviceStore = useDeviceStore()
const messageStore = useMessageStore()
const authStore = useAuthStore()
let wsBound = false
let heartbeatTimer: number | null = null
const language = ref(String(locale.value || 'zh-CN'))

const isLoginRoute = computed(() => route.path === '/login')

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

function changeLanguage(value: string) {
  locale.value = value
  language.value = value
  localStorage.setItem('juchuan_locale', value)
}

async function handleLogout() {
  try {
    await authStore.signOut()
    if (heartbeatTimer != null) {
      window.clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
    window.location.href = '/login'
  } catch {
    ElMessage.error(t('error.UNKNOWN'))
  }
}
</script>

<style scoped>
.app-root {
  position: relative;
  min-height: 100vh;
  padding: 14px;
  overflow: hidden;
}

.app-root::before {
  content: '';
  position: fixed;
  right: -120px;
  top: -88px;
  width: 420px;
  height: 420px;
  pointer-events: none;
  border-radius: 50%;
  background:
    radial-gradient(circle at center, rgba(109, 60, 26, 0.7) 0 20%, rgba(255, 197, 79, 0.22) 21% 60%, transparent 61%),
    conic-gradient(from 0deg, rgba(255, 205, 94, 0.62), rgba(232, 138, 35, 0.38), rgba(255, 198, 71, 0.62), rgba(214, 110, 29, 0.4), rgba(255, 205, 94, 0.62));
  filter: blur(2px);
  opacity: 0.3;
}

.app-root::after {
  content: '';
  position: fixed;
  left: -140px;
  bottom: -120px;
  width: 430px;
  height: 430px;
  border-radius: 50%;
  pointer-events: none;
  background:
    radial-gradient(circle at center, rgba(100, 55, 24, 0.55) 0 22%, rgba(251, 190, 72, 0.18) 23% 58%, transparent 59%),
    conic-gradient(from 160deg, rgba(250, 201, 84, 0.48), rgba(198, 101, 28, 0.28), rgba(255, 205, 94, 0.52), rgba(250, 161, 48, 0.3), rgba(250, 201, 84, 0.48));
  opacity: 0.22;
}

.orb {
  position: fixed;
  width: 280px;
  height: 280px;
  border-radius: 50%;
  filter: blur(40px);
  pointer-events: none;
  opacity: 0.2;
}

.orb--left {
  left: -120px;
  top: 48px;
  background: radial-gradient(circle at center, rgba(255, 197, 83, 0.72), transparent 72%);
}

.orb--right {
  right: -110px;
  top: 26%;
  background: radial-gradient(circle at center, rgba(227, 121, 33, 0.7), transparent 70%);
}

.shell {
  position: relative;
  z-index: 1;
  min-height: calc(100vh - 28px);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(240, 174, 58, 0.26);
  background:
    linear-gradient(155deg, rgba(88, 49, 23, 0.28), rgba(39, 22, 12, 0.4)),
    repeating-linear-gradient(16deg, rgba(255, 202, 84, 0.02) 0 2px, rgba(0, 0, 0, 0) 2px 9px);
  backdrop-filter: blur(4px);
}

.shell__aside {
  padding: 12px;
  border-right: 1px solid rgba(239, 176, 78, 0.24);
  background: linear-gradient(180deg, rgba(87, 48, 21, 0.36), rgba(41, 24, 13, 0.32));
}

.shell__main {
  padding: 14px;
}

.brand-card {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(248, 187, 78, 0.28);
  border-radius: 10px;
  padding: 12px;
  margin-bottom: 10px;
  background:
    radial-gradient(circle at 20% 20%, rgba(251, 200, 95, 0.14), transparent 38%),
    linear-gradient(135deg, rgba(106, 59, 28, 0.26), rgba(42, 24, 14, 0.36));
}

.brand-card::before {
  content: '';
  position: absolute;
  right: -32px;
  top: -26px;
  width: 160px;
  height: 160px;
  border-radius: 50%;
  pointer-events: none;
  background:
    radial-gradient(circle at center, rgba(110, 59, 23, 0.66) 0 18%, transparent 19%),
    repeating-conic-gradient(
      from 0deg,
      rgba(255, 217, 130, 0.24) 0deg 8deg,
      rgba(172, 92, 31, 0.1) 8deg 16deg
    );
  opacity: 0.52;
  filter: blur(0.4px);
}

.brand-card::after {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    radial-gradient(circle at 79% 24%, rgba(255, 227, 160, 0.18), transparent 35%),
    repeating-linear-gradient(22deg, rgba(255, 213, 118, 0.05) 0 2px, transparent 2px 11px);
  opacity: 0.5;
}

.brand-head {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: cover;
  border: 1px solid rgba(255, 216, 150, 0.56);
  box-shadow: 0 0 0 3px rgba(239, 176, 78, 0.18);
}

.brand-card__label {
  font-size: 12px;
  letter-spacing: 0.03em;
  color: #d5b57b;
}

.brand-card__title {
  margin: 6px 0 2px;
  font-size: 22px;
  font-weight: 600;
  letter-spacing: 0.03em;
  color: #ffe6aa;
}

.brand-card__hint {
  margin: 0;
  font-size: 12px;
  letter-spacing: 0.02em;
  color: #e0c48e;
}

.nav-menu {
  border-right: none;
  background: transparent;
}

.nav-menu :deep(.el-menu-item) {
  margin-bottom: 4px;
  border-radius: 8px;
  color: #f4d595;
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-icon {
  font-size: 16px;
}

.menu-label {
  font-size: 14px;
}

.nav-menu :deep(.el-menu-item.is-active) {
  color: #fff6d8;
  background: linear-gradient(120deg, rgba(229, 146, 37, 0.5), rgba(253, 198, 80, 0.34));
  box-shadow: inset 0 0 0 1px rgba(245, 188, 82, 0.5);
}

.topbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.lang-switch {
  width: 112px;
}

.topbar__title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  letter-spacing: 0.03em;
}

.topbar__path {
  margin: 6px 0 0;
  color: #e4c98f;
  font-size: 12px;
}

.topbar__badge {
  padding: 6px 10px;
  border-radius: 999px;
  border: 1px solid rgba(255, 204, 102, 0.55);
  color: #ffdf9b;
  background: rgba(117, 67, 26, 0.12);
  font-size: 11px;
  letter-spacing: 0.04em;
}

.content-stage {
  animation: panel-in 420ms ease-out both;
}

.login-stage {
  position: relative;
  min-height: calc(100vh - 36px);
  display: grid;
  place-items: center;
}

@keyframes panel-in {
  from {
    opacity: 0;
    transform: translateY(14px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 1024px) {
  .app-root {
    padding: 8px;
  }

  .shell {
    min-height: calc(100vh - 16px);
  }

  .topbar {
    flex-direction: column;
    gap: 10px;
  }

  .topbar-actions {
    width: 100%;
    flex-wrap: wrap;
  }
}

@media (max-width: 860px) {
  .shell {
    display: block;
    padding-bottom: 58px;
  }

  .shell__aside {
    position: fixed;
    z-index: 20;
    left: 8px;
    right: 8px;
    bottom: 8px;
    width: auto !important;
    border-right: 1px solid rgba(239, 176, 78, 0.24);
    border-bottom: none;
    border-radius: 8px;
    background: rgba(54, 32, 16, 0.46);
    backdrop-filter: blur(12px);
    box-shadow: 0 -8px 24px rgba(0, 0, 0, 0.24);
    padding: 8px;
  }

  .brand-card {
    display: none;
  }

  .nav-menu {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 4px;
  }

  .nav-menu :deep(.el-menu-item) {
    margin: 0;
    justify-content: center;
    flex-direction: column;
    min-width: 0;
    padding: 0 4px;
    height: 46px;
    gap: 2px;
    line-height: 1;
  }

  .menu-icon {
    font-size: 15px;
  }

  .menu-label {
    font-size: 11px;
  }

  .shell__main {
    padding: 12px;
  }

  .topbar__title {
    font-size: 24px;
  }
}
</style>
