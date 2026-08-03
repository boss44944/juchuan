<template>
  <section class="view-panel">
    <header class="panel-header">
      <div>
        <h2 class="panel-title">{{ t('menu.devices') }}</h2>
        <p class="panel-subtitle">NODE TOPOLOGY</p>
      </div>
      <div class="chip-group">
        <span class="stat-chip"><strong>{{ devices.length }}</strong>{{ t('devices.columns.name') }}</span>
        <span class="stat-chip"><strong>{{ onlineCount }}</strong>{{ t('devices.status.online') }}</span>
        <span class="stat-chip"><strong>{{ offlineCount }}</strong>{{ t('devices.status.offline') }}</span>
      </div>
    </header>

    <section class="entry-panel">
      <div class="entry-meta">
        <p class="entry-label">{{ t('devices.entryTitle') }}</p>
        <a :href="entryURL" target="_blank" rel="noreferrer" class="entry-link">{{ entryURL }}</a>
        <el-button size="small" @click="copyAddress">{{ t('devices.copyAddress') }}</el-button>
      </div>
      <div class="entry-qr-wrap">
        <img :src="qrImage" :alt="t('devices.qrAlt')" class="entry-qr" />
      </div>
    </section>

    <div class="table-scroll">
      <el-table :data="devices" min-width="640">
        <el-table-column prop="display_name" :label="t('devices.columns.name')" />
        <el-table-column prop="platform" :label="t('devices.columns.platform')" />
        <el-table-column :label="t('devices.columns.status')">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'online' ? 'success' : 'info'">
              {{ scope.row.status === 'online' ? t('devices.status.online') : t('devices.status.offline') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('devices.columns.actions')" width="220">
          <template #default="scope">
            <el-button size="small" @click="rename(scope.row.id)">{{ t('devices.actions.rename') }}</el-button>
            <el-button size="small" type="danger" @click="remove(scope.row.id)">{{ t('devices.actions.remove') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from '../stores/device'
import { ElMessage, ElMessageBox } from 'element-plus'
import { qrCodeURL, renameDevice, removeDevice, resolveApiErrorMessage } from '../api'

const { t } = useI18n()
const store = useDeviceStore()
const devices = computed(() => store.devices)
const onlineCount = computed(() => devices.value.filter((d) => d.status === 'online').length)
const offlineCount = computed(() => Math.max(devices.value.length - onlineCount.value, 0))
const entryURL = `${window.location.origin}/`
const qrImage = computed(() => qrCodeURL(entryURL))

onMounted(async () => {
  await store.load()
})

async function rename(id: string) {
  const target = devices.value.find((d) => d.id === id)
  if (!target) return

  try {
    const value = await ElMessageBox.prompt(t('devices.dialog.renameInput'), t('devices.dialog.renameTitle'), {
      inputValue: target.display_name,
    })
    const displayName = value.value.trim()
    if (!displayName) {
      return
    }
    await renameDevice({ id, display_name: displayName })
    await store.load()
    ElMessage.success(t('devices.toast.renamed'))
  } catch (err: any) {
    if (err === 'cancel' || err === 'close' || err?.action === 'cancel' || err?.action === 'close') {
      return
    }
    ElMessage.error(resolveApiErrorMessage(err))
  }
}

async function remove(id: string) {
  try {
    await ElMessageBox.confirm(t('devices.dialog.removeConfirm'), t('devices.dialog.removeTitle'))
    await removeDevice({ id })
    store.removeDevice(id)
    ElMessage.success(t('devices.toast.removed'))
  } catch (err: any) {
    if (err === 'cancel' || err === 'close' || err?.action === 'cancel' || err?.action === 'close') {
      return
    }
    ElMessage.error(resolveApiErrorMessage(err))
  }
}

async function copyAddress() {
  try {
    await navigator.clipboard.writeText(entryURL)
    ElMessage.success(t('devices.copySuccess'))
  } catch {
    ElMessage.warning(entryURL)
  }
}
</script>

<style scoped>
.chip-group {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.entry-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  border: 1px solid rgba(239, 176, 78, 0.22);
  border-radius: 6px;
  background: rgba(28, 21, 15, 0.5);
  padding: 10px;
  margin-bottom: 12px;
}

.entry-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.entry-label {
  margin: 0;
  color: #d0ba93;
  font-size: 12px;
}

.entry-link {
  font-size: 13px;
  color: #f5d79b;
  word-break: break-all;
}

.entry-qr-wrap {
  border: 1px solid rgba(239, 176, 78, 0.2);
  border-radius: 8px;
  background: rgba(15, 12, 9, 0.44);
  padding: 6px;
}

.entry-qr {
  width: 92px;
  height: 92px;
  border-radius: 6px;
  border: 1px solid rgba(223, 169, 73, 0.46);
  background:
    radial-gradient(circle at 24% 20%, rgba(255, 248, 230, 0.9), transparent 46%),
    repeating-linear-gradient(16deg, rgba(222, 187, 121, 0.08) 0 1px, rgba(0, 0, 0, 0) 1px 7px),
    #fff5e3;
  box-shadow:
    0 4px 10px rgba(35, 19, 8, 0.22),
    inset 0 0 0 1px rgba(255, 236, 196, 0.68);
  padding: 5px;
}

.table-scroll {
  overflow-x: auto;
}

@media (max-width: 860px) {
  .entry-panel {
    align-items: flex-start;
  }

  .entry-qr {
    width: 78px;
    height: 78px;
  }

  .table-scroll :deep(.el-table) {
    width: 640px;
  }
}
</style>
