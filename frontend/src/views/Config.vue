<template>
  <section class="page-surface config-view">
    <form class="config-form" @submit.prevent="save">
      <Card padding="lg" class="config-card">
        <div class="setting-heading">
          <span>01</span>
          <div><h3>{{ t('configPage.labels.port') }}</h3><p>HTTP / LAN</p></div>
        </div>
        <Input v-model="config.port" type="number" :aria-label="t('configPage.labels.port')" />
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
        <Input v-model="config.password" type="password" show-password :aria-label="t('configPage.labels.password')" />
      </Card>

      <Card padding="lg" class="config-card">
        <div class="setting-heading">
          <span>04</span>
          <div><h3>{{ t('configPage.labels.language') }}</h3><p>INTERFACE</p></div>
        </div>
        <Select
          v-model="config.language"
          :options="languageOptions"
          :aria-label="t('configPage.labels.language')"
          class="w-full"
        />
      </Card>

      <Button type="submit" variant="primary" size="lg" class="save-btn" :loading="saving">
        <Save :size="20" aria-hidden="true" />{{ t('configPage.save') }}
      </Button>
    </form>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Save } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
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
const languageOptions = computed(() => [
  { label: t('configPage.languages.zhCN'), value: 'zh-CN' },
  { label: t('configPage.languages.enUS'), value: 'en-US' },
  { label: t('configPage.languages.jaJP'), value: 'ja-JP' },
])
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
.config-form { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; max-width: 920px; }
.config-card { display: grid; align-content: start; gap: 20px; }
.setting-heading { display: flex; align-items: center; gap: 12px; }
.setting-heading > span { display: grid; width: 34px; height: 34px; place-items: center; border: 2px solid var(--brutal-border-color); background: var(--brutal-accent); font-weight: 500; }
.setting-heading h3 { margin: 0; font-size: 16px; }
.setting-heading p { margin: 2px 0 0; color: #9a4b1e; font-size: 9px; font-weight: 600; letter-spacing: .15em; }
.save-btn { grid-column: 1 / -1; justify-self: start; min-width: 190px; }

@media (max-width: 780px) {
  .config-form { grid-template-columns: 1fr; }
  .save-btn { width: 100%; min-height: 48px; }
}
</style>
