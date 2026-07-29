<template>
  <el-table :data="devices">
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
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from '../stores/device'
import { ElMessage, ElMessageBox } from 'element-plus'
import { renameDevice, removeDevice, resolveApiErrorMessage } from '../api'

const { t } = useI18n()
const store = useDeviceStore()
const devices = computed(() => store.devices)

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
</script>
