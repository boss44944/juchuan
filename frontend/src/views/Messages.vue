<template>
  <section class="page-surface messages-view" aria-labelledby="messages-title">
    <div class="messages-total">
      <span class="stat-chip"><strong>{{ total }}</strong>{{ t('messagesPage.table.content') }}</span>
    </div>

    <Card padding="lg" class="filter-panel">
      <div class="filter-grid">
        <Select
          :model-value="query.type"
          :options="typeOptions"
          :aria-label="t('messagesPage.table.type')"
          class="w-full"
          @update:model-value="(v) => { query.type = (v || '') as typeof query.type; refresh(1) }"
        />
        <Select
          :model-value="query.status"
          :options="statusOptions"
          :aria-label="t('messagesPage.table.status')"
          class="w-full"
          @update:model-value="(v) => { query.status = (v || '') as typeof query.status; refresh(1) }"
        />
        <Input v-model="query.device_id" :placeholder="t('messagesPage.filters.deviceId')" :aria-label="t('messagesPage.filters.deviceId')" @keyup.enter="refresh(1)" />
        <Select
          :model-value="query.sender_device_id"
          :options="senderOptions"
          :aria-label="t('messagesPage.filters.sender')"
          class="w-full"
          @update:model-value="(v) => { query.sender_device_id = v || ''; refresh(1) }"
        />
        <Select
          :model-value="query.target_device_id"
          :options="targetOptions"
          :aria-label="t('messagesPage.filters.target')"
          class="w-full"
          @update:model-value="(v) => { query.target_device_id = v || ''; refresh(1) }"
        />
      </div>
      <div class="filter-actions">
        <Button variant="primary" size="sm" @click="refresh(1)"><Search :size="16" aria-hidden="true" />{{ t('messagesPage.filters.apply') }}</Button>
        <Button variant="outline" size="sm" @click="reset"><RotateCcw :size="16" aria-hidden="true" />{{ t('messagesPage.filters.reset') }}</Button>
        <span class="action-divider" />
        <Button variant="outline" size="sm" :disabled="selectedRows.length === 0" @click="markReadBatch">{{ t('messagesPage.actions.batchRead') }}</Button>
        <Button variant="secondary" size="sm" :disabled="selectedRows.length === 0" @click="retryBatch">{{ t('messagesPage.actions.batchRetry') }}</Button>
        <Button variant="danger" size="sm" :disabled="selectedRows.length === 0" @click="deleteBatch">{{ t('messagesPage.actions.batchDelete') }}</Button>
        <span class="action-divider" />
        <Button variant="danger" size="sm" @click="openClearAllConfirm">{{ t('messagesPage.actions.clearAll') }}</Button>
      </div>
    </Card>

    <div v-if="loading" class="loading-state" role="status">{{ t('messagesPage.status.created') }}…</div>
    <div v-else-if="messages.length === 0" class="empty-state">{{ t('messagesPage.filters.typeAll') }} · 0</div>

    <template v-else>
      <div class="message-table-wrap">
        <table class="message-table">
          <thead><tr>
            <th><span class="visually-hidden">Select</span></th>
            <th>{{ t('messagesPage.table.time') }}</th><th>{{ t('messagesPage.table.type') }}</th><th>{{ t('messagesPage.table.content') }}</th>
            <th>{{ t('messagesPage.table.sender') }}</th><th>{{ t('messagesPage.table.target') }}</th><th>{{ t('messagesPage.table.status') }}</th><th>{{ t('messagesPage.table.operation') }}</th>
          </tr></thead>
          <tbody>
            <tr v-for="item in messages" :key="messageKey(item)" :class="{ 'row-highlight': messageKey(item) === highlightRowKey }">
              <td><input class="brutal-checkbox" type="checkbox" :checked="isSelected(item)" :aria-label="t('messagesPage.toast.selectFirst')" @change="toggleSelection(item)" /></td>
              <td>{{ formatTime(item.created_at) }}</td>
              <td><strong>{{ typeLabel(item.type) }}</strong></td>
              <td class="content-cell"><span v-if="item.type === 'TEXT'">{{ item.content }}</span><a v-else :href="fileURL(item)"><Download :size="15" aria-hidden="true" />{{ t('messagesPage.table.download') }}</a></td>
              <td><code>{{ item.sender_device_id || '—' }}</code></td><td><code>{{ item.target_device_id || '—' }}</code></td>
              <td><Badge :variant="statusType(item.status)" size="sm" dot>{{ statusLabel(item.status) }}</Badge></td>
              <td><div class="row-actions">
                <Button v-if="item.status !== 'READ' && canMarkRead(item)" variant="outline" size="sm" @click="markRead(item)">{{ t('messagesPage.actions.markRead') }}</Button>
                <Button variant="secondary" size="sm" @click="retry(item)"><RefreshCw :size="15" aria-hidden="true" />{{ t('messagesPage.actions.retry') }}</Button>
                <Button variant="danger" size="sm" @click="openDeleteConfirm(item)"><Trash2 :size="15" aria-hidden="true" /></Button>
              </div></td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="message-cards">
        <Card v-for="item in messages" :key="`mobile-${messageKey(item)}`" padding="lg" :class="messageKey(item) === highlightRowKey ? 'message-card row-highlight' : 'message-card'">
          <div class="message-card__head">
            <label><input class="brutal-checkbox" type="checkbox" :checked="isSelected(item)" @change="toggleSelection(item)" /> <strong>{{ typeLabel(item.type) }}</strong></label>
            <Badge :variant="statusType(item.status)" size="sm" dot>{{ statusLabel(item.status) }}</Badge>
          </div>
          <p class="message-time">{{ formatTime(item.created_at) }}</p>
          <p v-if="item.type === 'TEXT'" class="message-content">{{ item.content }}</p>
          <a v-else :href="fileURL(item)" class="download-link"><Download :size="16" />{{ t('messagesPage.table.download') }}</a>
          <dl><div><dt>{{ t('messagesPage.table.sender') }}</dt><dd>{{ item.sender_device_id || '—' }}</dd></div><div><dt>{{ t('messagesPage.table.target') }}</dt><dd>{{ item.target_device_id || '—' }}</dd></div></dl>
          <div class="row-actions">
            <Button v-if="item.status !== 'READ' && canMarkRead(item)" variant="outline" size="sm" @click="markRead(item)">{{ t('messagesPage.actions.markRead') }}</Button>
            <Button variant="secondary" size="sm" @click="retry(item)">{{ t('messagesPage.actions.retry') }}</Button>
            <Button variant="danger" size="sm" @click="openDeleteConfirm(item)"><Trash2 :size="15" /></Button>
          </div>
        </Card>
      </div>
    </template>

    <nav class="pagination" aria-label="Pagination">
      <Button variant="outline" size="sm" :disabled="query.page <= 1" @click="refresh(query.page - 1)"><ChevronLeft :size="17" /></Button>
      <Button v-for="page in visiblePages" :key="page" :variant="page === query.page ? 'primary' : 'outline'" size="sm" @click="refresh(page)">{{ page }}</Button>
      <Button variant="outline" size="sm" :disabled="query.page >= pageCount" @click="refresh(query.page + 1)"><ChevronRight :size="17" /></Button>
    </nav>

    <Teleport to="body">
      <div v-if="deleteConfirm.show" class="confirm-overlay" @click.self="closeDeleteConfirm">
        <div class="confirm-dialog">
          <strong>{{ t('messagesPage.dialog.deleteTitle') }}</strong>
          <p>{{ t('messagesPage.dialog.deleteConfirm') }}</p>
          <div class="confirm-actions">
            <Button variant="outline" size="sm" @click="closeDeleteConfirm">{{ t('messagesPage.dialog.cancel') }}</Button>
            <Button variant="danger" size="sm" @click="confirmDelete">{{ t('messagesPage.dialog.confirmDelete') }}</Button>
          </div>
        </div>
      </div>
      <div v-if="clearAllConfirm" class="confirm-overlay" @click.self="clearAllConfirm = false">
        <div class="confirm-dialog">
          <strong>{{ t('messagesPage.dialog.clearAllTitle') }}</strong>
          <p>{{ t('messagesPage.dialog.clearAllConfirm') }}</p>
          <div class="confirm-actions">
            <Button variant="outline" size="sm" @click="clearAllConfirm = false">{{ t('messagesPage.dialog.cancel') }}</Button>
            <Button variant="danger" size="sm" :loading="clearing" @click="confirmClearAll">{{ t('messagesPage.dialog.confirmClearAll') }}</Button>
          </div>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ChevronLeft, ChevronRight, Download, RefreshCw, RotateCcw, Search, Trash2 } from '@lucide/vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { useToast } from '@/composables/useToast'
import { useMessageStore, type MessageItem } from '../stores/message'
import { useDeviceStore } from '../stores/device'
import { clearMessages, deleteMessage, downloadFileURL, getMessages, resolveApiErrorMessage, sendFileMessage, sendTextMessage, updateMessageStatus } from '../api'

interface MessageListRow {
  message_id: string
  type: string
  content?: string
  file_id?: string
  sender_device_id?: string
  target_device_id?: string
  created_at?: string
  status?: 'CREATED' | 'DELIVERED' | 'READ'
}

const store = useMessageStore()
const deviceStore = useDeviceStore()
const { t } = useI18n()
const toast = useToast()
const messages = computed(() => store.messages)
const devices = computed(() => deviceStore.devices)
const typeOptions = computed(() => [
  { label: t('messagesPage.filters.typeAll'), value: '' },
  { label: t('messagesPage.filters.typeText'), value: 'TEXT' },
  { label: t('messagesPage.filters.typeFile'), value: 'FILE' },
])
const statusOptions = computed(() => [
  { label: t('messagesPage.filters.statusAll'), value: '' },
  { label: t('messagesPage.status.created'), value: 'CREATED' },
  { label: t('messagesPage.status.delivered'), value: 'DELIVERED' },
  { label: t('messagesPage.status.read'), value: 'READ' },
])
const senderOptions = computed(() => [
  { label: t('messagesPage.filters.sender'), value: '' },
  ...devices.value.map((device) => ({ label: device.display_name, value: device.id })),
])
const targetOptions = computed(() => [
  { label: t('messagesPage.filters.target'), value: '' },
  ...devices.value.map((device) => ({ label: device.display_name, value: device.id })),
])
const total = ref(0)
const loading = ref(false)
const selectedRows = ref<MessageItem[]>([])
const highlightRowKey = ref('')
const query = reactive({ page: 1, size: 20, type: '' as '' | 'TEXT' | 'FILE', status: '' as '' | 'CREATED' | 'DELIVERED' | 'READ', device_id: '', sender_device_id: '', target_device_id: '' })
const pageCount = computed(() => Math.max(1, Math.ceil(total.value / query.size)))
const visiblePages = computed(() => {
  const start = Math.max(1, Math.min(query.page - 2, pageCount.value - 4))
  return Array.from({ length: Math.min(5, pageCount.value) }, (_, index) => start + index)
})

onMounted(async () => { await deviceStore.load(); await refresh(1) })

function messageKey(item: MessageItem) { return item.row_key || `${item.id}:${item.target_device_id || ''}` }
function fileURL(item: MessageItem) { return item.file_id ? downloadFileURL(item.file_id) : `/download/${item.history_id || ''}` }
function statusType(status?: string): 'success' | 'accent' | 'outline' { return status === 'READ' ? 'success' : status === 'DELIVERED' ? 'accent' : 'outline' }
function statusLabel(status?: string) { return status === 'READ' ? t('messagesPage.status.read') : status === 'DELIVERED' ? t('messagesPage.status.delivered') : t('messagesPage.status.created') }
function typeLabel(type?: string) { return type === 'FILE' ? t('messagesPage.filters.typeFile') : t('messagesPage.filters.typeText') }
function formatTime(value?: string) {
  if (!value) return '—'
  // 兼容 Go time.String() 长格式（如 "2026-08-04 00:34:20.140165 +0800 CST m=+..."），
  // 只提取「日期 + 时:分」，避免把超长时间戳原样显示在列表里。
  const match = value.match(/^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2})/)
  if (match) {
    const [, y, m, d, h, min] = match
    return `${y}-${m}-${d} ${h}:${min}`
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${d} ${h}:${min}`
}
function canMarkRead(item: MessageItem) { const localID = localStorage.getItem('device_id') || ''; return !!localID && item.target_device_id === localID }
function isSelected(item: MessageItem) { return selectedRows.value.some((row) => messageKey(row) === messageKey(item)) }
function toggleSelection(item: MessageItem) { selectedRows.value = isSelected(item) ? selectedRows.value.filter((row) => messageKey(row) !== messageKey(item)) : [...selectedRows.value, item] }

async function markRead(item: MessageItem) {
  try { const localID = localStorage.getItem('device_id') || ''; if (!localID || !item.id) return; await updateMessageStatus({ message_id: item.id, device_id: localID, status: 'READ' }); store.updateStatus(item.id, 'READ', localID); toast.success(t('messagesPage.toast.markReadSuccess')) }
  catch (error) { toast.error(resolveApiErrorMessage(error, 'messagesPage.toast.markReadFailed')) }
}

async function retry(item: MessageItem) {
  try { const rowKey = await retryOne(item); await refresh(query.page); flashHighlight(rowKey); toast.success(t('messagesPage.toast.retrySuccess')) }
  catch (error) { toast.error(resolveApiErrorMessage(error, 'messagesPage.toast.retryFailed')) }
}

async function refresh(page: number) {
  loading.value = true
  try {
    query.page = page
    const response = await getMessages({ page: query.page, size: query.size, type: query.type, status: query.status, device_id: query.device_id.trim(), sender_device_id: query.sender_device_id, target_device_id: query.target_device_id })
    const data = response.data?.data || { items: [], total: 0 }
    total.value = Number(data.total || 0)
    const rows: MessageListRow[] = Array.isArray(data.items) ? data.items : []
    store.setMessages(rows.map((item) => ({ row_key: `${item.message_id}:${item.target_device_id || ''}`, id: item.message_id, type: item.type, content: item.content, file_id: item.file_id, sender_device_id: item.sender_device_id, target_device_id: item.target_device_id, created_at: item.created_at, status: item.status })))
    selectedRows.value = []
  } catch (error) { toast.error(resolveApiErrorMessage(error, 'messagesPage.toast.loadFailed')) }
  finally { loading.value = false }
}

async function reset() { Object.assign(query, { type: '', status: '', device_id: '', sender_device_id: '', target_device_id: '' }); await refresh(1) }

async function markReadBatch() {
  const rows = selectedRows.value.filter((row) => row.status !== 'READ' && canMarkRead(row))
  if (!rows.length) { toast.warning(t('messagesPage.toast.noneMarkable')); return }
  const localID = localStorage.getItem('device_id') || ''; let success = 0; let failed = 0; let lastError: unknown = null
  for (const row of rows) { try { await updateMessageStatus({ message_id: row.id, device_id: localID, status: 'READ' }); store.updateStatus(row.id, 'READ', localID); success += 1 } catch (error) { lastError = error; failed += 1 } }
  if (success && !failed) toast.success(t('messagesPage.toast.batchReadSuccess', { count: success }))
  else if (success) toast.warning(t('messagesPage.toast.batchPartial', { success, failed }))
  else toast.error(resolveApiErrorMessage(lastError, 'messagesPage.toast.batchReadFailed'))
}

async function retryBatch() {
  if (!selectedRows.value.length) { toast.warning(t('messagesPage.toast.selectFirst')); return }
  let success = 0; let failed = 0; let lastError: unknown = null; const rowKeys: string[] = []
  for (const row of selectedRows.value) { try { rowKeys.push(await retryOne(row)); success += 1 } catch (error) { lastError = error; failed += 1 } }
  if (success) { await refresh(query.page); flashHighlight(rowKeys.at(-1) || '') }
  if (success && !failed) toast.success(t('messagesPage.toast.batchRetrySuccess', { count: success }))
  else if (success) toast.warning(t('messagesPage.toast.batchPartial', { success, failed }))
  else toast.error(resolveApiErrorMessage(lastError, 'messagesPage.toast.batchRetryFailed'))
}

const deleteConfirm = ref<{ show: boolean; item: MessageItem | null }>({ show: false, item: null })
function openDeleteConfirm(item: MessageItem) { deleteConfirm.value = { show: true, item } }
function closeDeleteConfirm() { deleteConfirm.value = { show: false, item: null } }
async function confirmDelete() {
  const item = deleteConfirm.value.item
  if (!item) return
  try {
    const localID = localStorage.getItem('device_id') || ''
    await deleteMessage({ message_id: item.id, device_id: item.target_device_id || localID })
    store.removeMessage(messageKey(item))
    total.value = Math.max(0, total.value - 1)
    toast.success(t('messagesPage.toast.deleteSuccess'))
  } catch (error) {
    toast.error(resolveApiErrorMessage(error, 'messagesPage.toast.deleteFailed'))
  }
  closeDeleteConfirm()
}
async function deleteBatch() {
  if (!selectedRows.value.length) { toast.warning(t('messagesPage.toast.selectFirst')); return }
  let success = 0; let failed = 0; let lastError: unknown = null
  const localID = localStorage.getItem('device_id') || ''
  for (const row of selectedRows.value) {
    try { await deleteMessage({ message_id: row.id, device_id: row.target_device_id || localID }); store.removeMessage(messageKey(row)); success += 1 } catch (error) { lastError = error; failed += 1 }
  }
  total.value = Math.max(0, total.value - success)
  selectedRows.value = []
  if (success && !failed) toast.success(t('messagesPage.toast.batchDeleteSuccess', { count: success }))
  else if (success) toast.warning(t('messagesPage.toast.batchPartial', { success, failed }))
  else toast.error(resolveApiErrorMessage(lastError, 'messagesPage.toast.batchDeleteFailed'))
}
const clearAllConfirm = ref(false)
const clearing = ref(false)
function openClearAllConfirm() { clearAllConfirm.value = true }
async function confirmClearAll() {
  clearing.value = true
  try {
    const localID = localStorage.getItem('device_id') || ''
    await clearMessages({ device_id: localID })
    store.clearMessages(localID)
    total.value = 0
    toast.success(t('messagesPage.toast.clearAllSuccess'))
  } catch (error) {
    toast.error(resolveApiErrorMessage(error, 'messagesPage.toast.clearAllFailed'))
  }
  clearing.value = false
  clearAllConfirm.value = false
}

function responseMessageID(data: unknown, fallback: string) {
  if (!data || typeof data !== 'object') return fallback
  const outer = data as Record<string, unknown>; const nested = outer.data && typeof outer.data === 'object' ? outer.data as Record<string, unknown> : outer
  return typeof nested.id === 'string' ? nested.id : fallback
}

async function retryOne(item: MessageItem) {
  const senderID = localStorage.getItem('device_id') || ''; const target = item.target_device_id
  if (!senderID || !target) throw new Error('missing device info')
  if (item.type === 'TEXT') { const response = await sendTextMessage({ content: item.content || '', sender_device_id: senderID, targets: [target] }); return `${responseMessageID(response.data, item.id)}:${target}` }
  if (item.type === 'FILE' && item.file_id) { const response = await sendFileMessage({ file_id: item.file_id, sender_device_id: senderID, targets: [target] }); return `${responseMessageID(response.data, item.id)}:${target}` }
  throw new Error('unsupported message')
}

function flashHighlight(rowKey: string) { if (!rowKey) return; highlightRowKey.value = rowKey; window.setTimeout(() => { if (highlightRowKey.value === rowKey) highlightRowKey.value = '' }, 2200) }
</script>

<style scoped>
.messages-total { display: flex; justify-content: flex-end; margin-bottom: 14px; }
.filter-panel { margin-bottom: 20px; background: var(--brutal-muted); }
.filter-grid { display: grid; grid-template-columns: repeat(5, minmax(145px, 1fr)); gap: 10px; }
.filter-actions { display: flex; flex-wrap: wrap; align-items: center; gap: 10px; margin-top: 14px; }
.action-divider { width: 3px; height: 30px; background: var(--brutal-border-color); }
.message-table-wrap { overflow-x: auto; border: 2px solid var(--brutal-border-color); box-shadow: 3px 3px 0 var(--brutal-shadow-color); }
.message-table { width: 100%; min-width: 1160px; border-collapse: collapse; background: var(--brutal-bg); }
.message-table th, .message-table td { padding: 9px 9px; border-right: 2px solid var(--brutal-border-color); border-bottom: 2px solid var(--brutal-border-color); text-align: left; vertical-align: middle; }
.message-table th { background: var(--brutal-primary); font-size: 12px; font-weight: 500; letter-spacing: .05em; }
.message-table tr:last-child td { border-bottom: 0; }
.message-table th:last-child, .message-table td:last-child { border-right: 0; }
.message-table td:last-child { min-width: 200px; }
.content-cell { max-width: 260px; overflow-wrap: anywhere; }
.content-cell a, .download-link { display: inline-flex; align-items: center; gap: 5px; font-weight: 600; }
.row-actions { display: flex; flex-wrap: nowrap; align-items: center; gap: 7px; white-space: nowrap; }
.row-highlight { background: #fff0bd !important; animation: row-flash 650ms ease 2; }
.pagination { display: flex; justify-content: center; gap: 8px; margin-top: 20px; }
.message-cards { display: none; }
@keyframes row-flash { 50% { background: var(--brutal-accent); } }
.confirm-overlay { position: fixed; inset: 0; z-index: 1000; display: grid; place-items: center; background: rgba(0, 0, 0, 0.5); }
.confirm-dialog { width: 90%; max-width: 400px; padding: 20px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: var(--brutal-bg); box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.confirm-dialog strong { display: block; margin-bottom: 8px; font-size: 16px; }
.confirm-dialog p { margin: 0 0 20px; color: var(--brutal-muted-foreground); font-size: 14px; }
.confirm-actions { display: flex; gap: 10px; justify-content: flex-end; }

@media (max-width: 1100px) { .filter-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 780px) {
  .filter-grid { grid-template-columns: 1fr; }
  .action-divider { width: 100%; height: 3px; }
  .message-table-wrap { display: none; }
  .message-cards { display: grid; gap: 15px; }
  .message-card { display: grid; gap: 12px; }
  .message-card__head { display: flex; align-items: center; justify-content: space-between; gap: 10px; }
  .message-card__head label { display: flex; align-items: center; gap: 8px; }
  .message-time { margin: 0; color: var(--brutal-muted-foreground); font-size: 12px; }
  .message-content { margin: 0; padding: 12px; border: 2px solid var(--brutal-border-color); background: var(--brutal-muted); overflow-wrap: anywhere; }
  dl { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin: 0; }
  dl div { min-width: 0; }
  dt { color: var(--brutal-muted-foreground); font-size: 10px; font-weight: 600; }
  dd { margin: 3px 0 0; overflow: hidden; font-size: 12px; text-overflow: ellipsis; white-space: nowrap; }
  .filter-actions :deep(button) { min-height: 44px; }
  .row-actions { flex-wrap: wrap; }
  .row-actions :deep(button) { flex: 1 1 auto; min-height: 44px; }
}
</style>
