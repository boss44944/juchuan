<template>
  <section class="page-surface config-view" aria-labelledby="config-title">
    <header class="panel-header">
      <div>
        <p class="panel-subtitle">JUCHUAN / LOCAL SETTINGS</p>
        <h2 id="config-title" class="panel-title">{{ t('menu.config') }}</h2>
      </div>
      <span class="config-mark"><SlidersHorizontal :size="24" aria-hidden="true" /></span>
    </header>

    <form class="config-form" @submit.prevent="save">
      <Card padding="lg" class="config-card">
        <div class="setting-heading">
          <span>01</span>
          <div><h3>{{ t('configPage.labels.port') }}</h3><p>HTTP / LAN</p></div>
        </div>
        <Input v-model="config.port" type="number" size="lg" :aria-label="t('configPage.labels.port')" />
      </Card>

      <Card padding="lg" class="config-card">
        <div class="setting-heading">
          <span>02</span>
          <div><h3>{{ t('configPage.labels.autoOpen') }}</h3><p>BROWSER</p></div>
        </div>
        <label class="brutal-switch">
          <input v-model="config.auto_open" type="checkbox" />
          <span class="brutal-switch__track" aria-hidden="true" />
          <strong>{{ config.auto_open ? 'ON' : 'OFF' }}</strong>
        </label>
      </Card>

      <Card padding="lg" class="config-card">
        <div class="setting-heading">
          <span>03</span>
          <div><h3>{{ t('configPage.labels.password') }}</h3><p>ACCESS</p></div>
        </div>
        <Input v-model="config.password" type="password" show-password size="lg" :aria-label="t('configPage.labels.password')" />
      </Card>

      <Card padding="lg" class="config-card">
        <div class="setting-heading">
          <span>04</span>
          <div><h3>{{ t('configPage.labels.language') }}</h3><p>INTERFACE</p></div>
        </div>
        <select v-model="config.language" class="brutal-select" :aria-label="t('configPage.labels.language')">
          <option value="zh-CN">{{ t('configPage.languages.zhCN') }}</option>
          <option value="en-US">{{ t('configPage.languages.enUS') }}</option>
          <option value="ja-JP">{{ t('configPage.languages.jaJP') }}</option>
        </select>
      </Card>

      <Button type="submit" variant="primary" size="lg" class="save-btn" :loading="saving">
        <Save :size="20" aria-hidden="true" />{{ t('configPage.save') }}
      </Button>
    </form>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Save, SlidersHorizontal } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useToast } from '@/composables/useToast'
import { getConfig, resolveApiErrorMessage, updateConfig } from '../api'

interface ConfigForm {
  port: string
  auto_open: boolean
  password: string
  language: string
}

const { t } = useI18n()
const toast = useToast()
const saving = ref(false)
const config = reactive<ConfigForm>({ port: '8000', auto_open: false, password: '', language: 'zh-CN' })

onMounted(async () => {
  try {
    const response = await getConfig()
    const data = response.data ?? {}
    config.port = String(data.port ?? 8000)
    config.auto_open = Boolean(data.auto_open)
    config.password = String(data.password ?? '')
    config.language = String(data.language ?? 'zh-CN')
  } catch (error) {
    toast.error(resolveApiErrorMessage(error))
  }
})

async function save() {
  saving.value = true
  try {
    await updateConfig({
      port: Number(config.port),
      auto_open: config.auto_open,
      password: config.password,
    })
    toast.success(t('configPage.saved'))
  } catch (error) {
    toast.error(resolveApiErrorMessage(error))
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.config-mark { display: grid; width: 50px; height: 50px; place-items: center; border: 3px solid var(--brutal-border-color); background: var(--brutal-primary); box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.config-form { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; max-width: 920px; }
.config-card { display: grid; align-content: start; gap: 22px; }
.setting-heading { display: flex; align-items: center; gap: 12px; }
.setting-heading > span { display: grid; width: 38px; height: 38px; place-items: center; border: 3px solid var(--brutal-border-color); background: var(--brutal-accent); font-weight: 950; }
.setting-heading h3 { margin: 0; font-size: 18px; }
.setting-heading p { margin: 2px 0 0; color: #9a4b1e; font-size: 9px; font-weight: 900; letter-spacing: .15em; }
.save-btn { grid-column: 1 / -1; justify-self: start; min-width: 190px; }

@media (max-width: 680px) {
  .config-form { grid-template-columns: 1fr; }
  .save-btn { width: 100%; }
}
</style>
