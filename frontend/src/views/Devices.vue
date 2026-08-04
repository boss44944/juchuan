<template>
  <section class="page-surface devices-view" aria-labelledby="devices-title">
    <header class="panel-header">
      <div>
        <p class="panel-subtitle">JUCHUAN / NEARBY DEVICES</p>
        <h2 id="devices-title" class="panel-title">{{ t('menu.devices') }}</h2>
      </div>
      <div class="chip-group">
        <span class="stat-chip"><strong>{{ devices.length }}</strong>{{ t('devices.columns.name') }}</span>
        <span class="stat-chip stat-chip--online"><strong>{{ onlineCount }}</strong>{{ t('devices.status.online') }}</span>
      </div>
    </header>

    <Card padding="lg" class="entry-panel">
      <div class="entry-meta">
        <span class="step-marker">ACCESS</span>
        <h3>{{ t('devices.entryTitle') }}</h3>
        <a :href="entryURL" target="_blank" rel="noreferrer">{{ entryURL }}</a>
        <Button variant="outline" size="sm" @click="copyAddress">
          <Copy :size="17" aria-hidden="true" />{{ t('devices.copyAddress') }}
        </Button>
      </div>
      <div class="entry-qr-wrap">
        <img :src="qrImage" :alt="t('devices.qrAlt')" class="entry-qr" />
      </div>
    </Card>

    <div v-if="loading" class="loading-state" role="status">{{ t('messagesPage.status.created') }}…</div>
    <div v-else-if="loadError" class="empty-state" role="alert">{{ loadError }}</div>
    <div v-else-if="devices.length === 0" class="empty-state">{{ t('send.noDevices') }}</div>

    <div v-else class="device-grid">
      <Card v-for="device in devices" :key="device.id" padding="lg" class="device-card">
        <div class="device-card__head">
          <span class="device-icon"><MonitorSmartphone :size="24" aria-hidden="true" /></span>
          <Badge :variant="device.status === 'online' ? 'success' : 'outline'" size="sm" dot>
            {{ device.status === 'online' ? t('devices.status.online') : t('devices.status.offline') }}
          </Badge>
        </div>
        <div>
          <h3>{{ device.display_name }}</h3>
          <p>{{ device.platform || '—' }}</p>
          <code>{{ device.id }}</code>
        </div>
        <div class="device-actions">
          <Button variant="outline" size="sm" @click="openRename(device.id)">
            <Pencil :size="16" aria-hidden="true" />{{ t('devices.actions.rename') }}
          </Button>
          <Button variant="danger" size="sm" @click="openRemove(device.id)">
            <Trash2 :size="16" aria-hidden="true" />{{ t('devices.actions.remove') }}
          </Button>
        </div>
      </Card>
    </div>

    <div v-if="dialog.kind" class="modal-backdrop" @click.self="closeDialog">
      <section ref="dialogElement" class="modal-card" role="dialog" aria-modal="true" :aria-labelledby="dialogTitleID" @keydown.esc="closeDialog">
        <h3 :id="dialogTitleID">{{ dialog.kind === 'rename' ? t('devices.dialog.renameTitle') : t('devices.dialog.removeTitle') }}</h3>
        <template v-if="dialog.kind === 'rename'">
          <p>{{ t('devices.dialog.renameInput') }}</p>
          <Input ref="renameInput" v-model="dialog.name" size="lg" :aria-label="t('devices.dialog.renameInput')" @keyup.enter="confirmDialog" />
        </template>
        <p v-else>{{ t('devices.dialog.removeConfirm') }}</p>
        <div class="modal-actions">
          <Button variant="outline" @click="closeDialog">{{ t('messagesPage.filters.reset') }}</Button>
          <Button :variant="dialog.kind === 'remove' ? 'danger' : 'primary'" :loading="dialog.busy" @click="confirmDialog">
            {{ dialog.kind === 'remove' ? t('devices.actions.remove') : t('devices.actions.rename') }}
          </Button>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Copy, MonitorSmartphone, Pencil, Trash2 } from '@lucide/vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useToast } from '@/composables/useToast'
import { useDeviceStore } from '../stores/device'
import { qrCodeURL, renameDevice, removeDevice, resolveApiErrorMessage } from '../api'

const { t } = useI18n()
const toast = useToast()
const store = useDeviceStore()
const devices = computed(() => store.devices)
const onlineCount = computed(() => devices.value.filter((device) => device.status === 'online').length)
const entryURL = `${window.location.origin}/client/inbox`
const qrImage = computed(() => qrCodeURL(entryURL))
const loading = ref(true)
const loadError = ref('')
const dialogElement = ref<HTMLElement | null>(null)
const renameInput = ref<{ focus: () => void } | null>(null)
const dialogTitleID = 'device-dialog-title'
const dialog = reactive({ kind: '' as '' | 'rename' | 'remove', id: '', name: '', busy: false })

onMounted(async () => {
  try {
    await store.load()
  } catch (error) {
    loadError.value = resolveApiErrorMessage(error)
  } finally {
    loading.value = false
  }
})

async function openRename(id: string) {
  const target = devices.value.find((device) => device.id === id)
  if (!target) return
  Object.assign(dialog, { kind: 'rename', id, name: target.display_name, busy: false })
  await nextTick()
  renameInput.value?.focus()
}

async function openRemove(id: string) {
  Object.assign(dialog, { kind: 'remove', id, name: '', busy: false })
  await nextTick()
  dialogElement.value?.focus()
}

function closeDialog() {
  if (dialog.busy) return
  Object.assign(dialog, { kind: '', id: '', name: '', busy: false })
}

async function confirmDialog() {
  if (!dialog.kind || !dialog.id || dialog.busy) return
  dialog.busy = true
  try {
    if (dialog.kind === 'rename') {
      const displayName = dialog.name.trim()
      if (!displayName) return
      await renameDevice({ id: dialog.id, display_name: displayName })
      await store.load()
      toast.success(t('devices.toast.renamed'))
    } else {
      await removeDevice({ id: dialog.id })
      store.removeDevice(dialog.id)
      toast.success(t('devices.toast.removed'))
    }
    closeDialogAfterSuccess()
  } catch (error) {
    toast.error(resolveApiErrorMessage(error))
  } finally {
    dialog.busy = false
  }
}

function closeDialogAfterSuccess() {
  Object.assign(dialog, { kind: '', id: '', name: '', busy: false })
}

async function copyAddress() {
  try {
    await navigator.clipboard.writeText(entryURL)
    toast.success(t('devices.copySuccess'))
  } catch {
    toast.warning(entryURL)
  }
}
</script>

<style scoped>
.chip-group { display: flex; flex-wrap: wrap; gap: 8px; }
.stat-chip--online { background: #dbe9b8; }
.entry-panel { display: flex; align-items: center; justify-content: space-between; gap: 20px; margin-bottom: 22px; background: var(--brutal-accent); }
.entry-meta { display: grid; justify-items: start; gap: 8px; min-width: 0; }
.entry-meta h3 { margin: 0; font-size: 22px; }
.entry-meta a { word-break: break-all; font-weight: 800; }
.step-marker { padding: 3px 7px; border: 2px solid var(--brutal-border-color); background: var(--brutal-bg); font-size: 10px; font-weight: 900; letter-spacing: .12em; }
.entry-qr-wrap { flex: 0 0 auto; padding: 8px; border: 3px solid var(--brutal-border-color); background: #fff; }
.entry-qr { display: block; width: 112px; height: 112px; padding: 5px; background: #fff; }
.device-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 18px; }
.device-card { display: flex; flex-direction: column; gap: 18px; min-width: 0; }
.device-card__head { display: flex; align-items: center; justify-content: space-between; }
.device-icon { display: grid; width: 44px; height: 44px; place-items: center; border: 3px solid var(--brutal-border-color); background: var(--brutal-primary); box-shadow: 3px 3px 0 var(--brutal-shadow-color); }
.device-card h3 { margin: 0 0 5px; font-size: 21px; overflow-wrap: anywhere; }
.device-card p { margin: 0 0 8px; color: var(--brutal-muted-foreground); }
.device-card code { display: block; overflow: hidden; color: var(--brutal-muted-foreground); font-size: 11px; text-overflow: ellipsis; white-space: nowrap; }
.device-actions { display: flex; flex-wrap: wrap; gap: 9px; margin-top: auto; }

@media (max-width: 700px) {
  .panel-header, .entry-panel { flex-direction: column; }
  .entry-panel { align-items: stretch; }
  .entry-qr-wrap { align-self: center; }
}

@media (max-width: 520px) {
  .chip-group, .device-actions { display: grid; grid-template-columns: 1fr 1fr; width: 100%; }
  .device-grid { grid-template-columns: minmax(0, 1fr); gap: 14px; }
  .device-actions :deep(button) { width: 100%; min-height: 44px; }
  .entry-meta :deep(button) { width: 100%; min-height: 44px; }
}
</style>
