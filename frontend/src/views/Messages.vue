<template>
  <section class="view-panel messages-view">
    <header class="panel-header">
      <div>
        <h2 class="panel-title">{{ t('menu.messages') }}</h2>
        <p class="panel-subtitle">DELIVERY TELEMETRY STREAM</p>
      </div>
      <span class="stat-chip"><strong>{{ total }}</strong>{{ t('messagesPage.table.content') }}</span>
    </header>

    <el-space wrap class="filters-row">
      <el-select v-model="query.type" style="width: 160px" @change="refresh(1)">
        <el-option :label="t('messagesPage.filters.typeAll')" value="" />
        <el-option :label="t('messagesPage.filters.typeText')" value="TEXT" />
        <el-option :label="t('messagesPage.filters.typeFile')" value="FILE" />
      </el-select>

      <el-select v-model="query.status" style="width: 160px" @change="refresh(1)">
        <el-option :label="t('messagesPage.filters.statusAll')" value="" />
        <el-option :label="t('messagesPage.status.created')" value="CREATED" />
        <el-option :label="t('messagesPage.status.delivered')" value="DELIVERED" />
        <el-option :label="t('messagesPage.status.read')" value="READ" />
      </el-select>

      <el-input
        v-model="query.device_id"
        style="width: 220px"
        :placeholder="t('messagesPage.filters.deviceId')"
        @keyup.enter="refresh(1)"
      />

      <el-select v-model="query.sender_device_id" clearable style="width: 180px" :placeholder="t('messagesPage.filters.sender')" @change="refresh(1)">
        <el-option v-for="d in devices" :key="`s-${d.id}`" :label="d.display_name" :value="d.id" />
      </el-select>

      <el-select v-model="query.target_device_id" clearable style="width: 180px" :placeholder="t('messagesPage.filters.target')" @change="refresh(1)">
        <el-option v-for="d in devices" :key="`t-${d.id}`" :label="d.display_name" :value="d.id" />
      </el-select>

      <el-button @click="refresh(1)">{{ t('messagesPage.filters.apply') }}</el-button>
      <el-button @click="reset">{{ t('messagesPage.filters.reset') }}</el-button>
      <el-divider direction="vertical" />
      <el-button :disabled="selectedRows.length === 0" @click="markReadBatch">{{ t('messagesPage.actions.batchRead') }}</el-button>
      <el-button :disabled="selectedRows.length === 0" @click="retryBatch">{{ t('messagesPage.actions.batchRetry') }}</el-button>
    </el-space>

    <div class="table-scroll">
      <el-table
        v-loading="loading"
        :data="messages"
        row-key="row_key"
        :row-class-name="rowClassName"
        style="margin-top: 16px"
        @selection-change="onSelectionChange"
      >
        <el-table-column type="selection" width="44" />
        <el-table-column :label="t('messagesPage.table.time')" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('messagesPage.table.type')" width="100">
          <template #default="scope">
            {{ typeLabel(scope.row.type) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('messagesPage.table.content')">
          <template #default="scope">
            <span v-if="scope.row.type === 'TEXT'">{{ scope.row.content }}</span>
            <a v-else :href="fileURL(scope.row)">{{ t('messagesPage.table.download') }}</a>
          </template>
        </el-table-column>
        <el-table-column prop="sender_device_id" :label="t('messagesPage.table.sender')" width="140" />
        <el-table-column prop="target_device_id" :label="t('messagesPage.table.target')" width="140" />
        <el-table-column :label="t('messagesPage.table.status')" width="140">
          <template #default="scope">
            <el-tag :type="statusType(scope.row.status)">{{ statusLabel(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('messagesPage.table.operation')" width="220">
          <template #default="scope">
            <el-button
              v-if="scope.row.status !== 'READ' && canMarkRead(scope.row)"
              size="small"
              @click="markRead(scope.row)"
            >
              {{ t('messagesPage.actions.markRead') }}
            </el-button>
            <el-button size="small" @click="retry(scope.row)">{{ t('messagesPage.actions.retry') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-pagination
      style="margin-top: 16px"
      layout="total, prev, pager, next"
      :total="total"
      :current-page="query.page"
      :page-size="query.size"
      @current-change="refresh"
    />
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useMessageStore } from '../stores/message'
import { useDeviceStore } from '../stores/device'
import { downloadFileURL, getMessages, resolveApiErrorMessage, sendFileMessage, sendTextMessage, updateMessageStatus } from '../api'
import type { MessageItem } from '../stores/message'

const store = useMessageStore()
const { t } = useI18n()
const deviceStore = useDeviceStore()
const messages = computed(() => store.messages)
const devices = computed(() => deviceStore.devices)
const total = ref(0)
const loading = ref(false)
const selectedRows = ref<MessageItem[]>([])
const highlightRowKey = ref('')

const query = reactive({
  page: 1,
  size: 20,
  type: '' as '' | 'TEXT' | 'FILE',
  status: '' as '' | 'CREATED' | 'DELIVERED' | 'READ',
  device_id: '',
  sender_device_id: '',
  target_device_id: '',
})

onMounted(async () => {
  await deviceStore.load()
  await refresh(1)
})

function fileURL(item: MessageItem) {
  if (item.file_id) {
    return downloadFileURL(item.file_id)
  }
  return `/download/${item.history_id || ''}`
}

function statusType(status?: string) {
  if (status === 'READ') return 'success'
  if (status === 'DELIVERED') return 'warning'
  return 'info'
}

function statusLabel(status?: string) {
  if (status === 'READ') return t('messagesPage.status.read')
  if (status === 'DELIVERED') return t('messagesPage.status.delivered')
  return t('messagesPage.status.created')
}

function typeLabel(type?: string) {
  if (type === 'FILE') return t('messagesPage.filters.typeFile')
  return t('messagesPage.filters.typeText')
}

function formatTime(value?: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString()
}

function canMarkRead(item: MessageItem) {
  const localID = localStorage.getItem('device_id') || ''
  return !!localID && item.target_device_id === localID
}

function onSelectionChange(rows: MessageItem[]) {
  selectedRows.value = rows
}

function rowClassName(arg: { row: MessageItem }) {
  const key = arg.row.row_key || `${arg.row.id}:${arg.row.target_device_id || ''}`
  return key !== '' && key === highlightRowKey.value ? 'row-highlight' : ''
}

async function markRead(item: MessageItem) {
  try {
    const localID = localStorage.getItem('device_id') || ''
    if (!localID || !item.id) return
    await updateMessageStatus({
      message_id: item.id,
      device_id: localID,
      status: 'READ',
    })
    store.updateStatus(item.id, 'READ', localID)
    ElMessage.success(t('messagesPage.toast.markReadSuccess'))
  } catch (err) {
    ElMessage.error(resolveApiErrorMessage(err, 'messagesPage.toast.markReadFailed'))
  }
}

async function retry(item: MessageItem) {
  try {
    const rowKey = await retryOne(item)
    await refresh(query.page)
    flashHighlight(rowKey)
    ElMessage.success(t('messagesPage.toast.retrySuccess'))
  } catch (err) {
    ElMessage.error(resolveApiErrorMessage(err, 'messagesPage.toast.retryFailed'))
  }
}

async function refresh(page: number) {
  loading.value = true
  try {
    query.page = page
    const res = await getMessages({
      page: query.page,
      size: query.size,
      type: query.type,
      status: query.status,
      device_id: query.device_id.trim(),
      sender_device_id: query.sender_device_id,
      target_device_id: query.target_device_id,
    })

    const data = res.data?.data || { items: [], total: 0 }
    total.value = data.total || 0
    const rows = Array.isArray(data.items) ? data.items : []
    store.setMessages(rows.map((item: any) => ({
      row_key: `${item.message_id}:${item.target_device_id || ''}`,
      id: item.message_id,
      type: item.type,
      content: item.content,
      file_id: item.file_id,
      sender_device_id: item.sender_device_id,
      target_device_id: item.target_device_id,
      created_at: item.created_at,
      status: item.status,
    })))
    selectedRows.value = []
  } catch (err) {
    ElMessage.error(resolveApiErrorMessage(err, 'messagesPage.toast.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function reset() {
  query.type = ''
  query.status = ''
  query.device_id = ''
  query.sender_device_id = ''
  query.target_device_id = ''
  await refresh(1)
}

async function markReadBatch() {
  const rows = selectedRows.value.filter((row) => row.status !== 'READ' && canMarkRead(row))
  if (rows.length === 0) {
    ElMessage.warning(t('messagesPage.toast.noneMarkable'))
    return
  }

  const localID = localStorage.getItem('device_id') || ''
  let success = 0
  let failed = 0
  let lastError: unknown = null
  for (const row of rows) {
    if (!row.id) continue
    try {
      await updateMessageStatus({
        message_id: row.id,
        device_id: localID,
        status: 'READ',
      })
      store.updateStatus(row.id, 'READ', localID)
      success += 1
    } catch (err) {
      lastError = err
      failed += 1
    }
  }

  if (success > 0 && failed === 0) {
    ElMessage.success(t('messagesPage.toast.batchReadSuccess', { count: success }))
  } else if (success > 0 && failed > 0) {
    ElMessage.warning(t('messagesPage.toast.batchPartial', { success, failed }))
  } else {
    ElMessage.error(resolveApiErrorMessage(lastError, 'messagesPage.toast.batchReadFailed'))
  }
}

async function retryBatch() {
  const rows = selectedRows.value
  if (rows.length === 0) {
    ElMessage.warning(t('messagesPage.toast.selectFirst'))
    return
  }

  let success = 0
  let failed = 0
  let lastError: unknown = null
  const rowKeys: string[] = []
  for (const row of rows) {
    try {
      const rowKey = await retryOne(row)
      rowKeys.push(rowKey)
      success += 1
    } catch (err) {
      lastError = err
      failed += 1
    }
  }

  if (success > 0 && failed === 0) {
    await refresh(query.page)
    flashHighlight(rowKeys[rowKeys.length - 1] || '')
    ElMessage.success(t('messagesPage.toast.batchRetrySuccess', { count: success }))
  } else if (success > 0 && failed > 0) {
    await refresh(query.page)
    flashHighlight(rowKeys[rowKeys.length - 1] || '')
    ElMessage.warning(t('messagesPage.toast.batchPartial', { success, failed }))
  } else {
    ElMessage.error(resolveApiErrorMessage(lastError, 'messagesPage.toast.batchRetryFailed'))
  }
}

async function retryOne(item: MessageItem) {
  const senderID = localStorage.getItem('device_id') || ''
  const target = item.target_device_id
  if (!senderID || !target) {
    throw new Error('missing device info')
  }

  if (item.type === 'TEXT') {
    const res: any = await sendTextMessage({
      content: item.content || '',
      sender_device_id: senderID,
      targets: [target],
    })
    const nextID = res.data?.data?.id || res.data?.id || item.id
    return `${nextID}:${target}`
  }

  if (item.type === 'FILE' && item.file_id) {
    const res: any = await sendFileMessage({
      file_id: item.file_id,
      sender_device_id: senderID,
      targets: [target],
    })
    const nextID = res.data?.data?.id || res.data?.id || item.id
    return `${nextID}:${target}`
  }

  throw new Error('unsupported message')
}

function flashHighlight(rowKey: string) {
  if (!rowKey) return
  highlightRowKey.value = rowKey
  window.setTimeout(() => {
    if (highlightRowKey.value === rowKey) {
      highlightRowKey.value = ''
    }
  }, 2200)
}
</script>

<style scoped>
:deep(.row-highlight > td) {
  background: rgba(255, 202, 113, 0.2) !important;
}

.filters-row {
  margin-bottom: 6px;
}

.table-scroll {
  overflow-x: auto;
}

@media (max-width: 860px) {
  .table-scroll :deep(.el-table) {
    min-width: 980px;
  }
}
</style>
