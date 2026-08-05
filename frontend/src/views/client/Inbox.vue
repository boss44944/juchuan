<template>
  <section class="inbox-view">
    <div class="client-page-head client-page-head--bare">
      <div class="client-page-head__actions">
        <Button variant="outline" size="sm" :loading="loading" @click="loadInbox(1)"><RefreshCw :size="17" />{{ t('client.inbox.refresh') }}</Button>
        <Button v-if="messages.length > 0" variant="danger" size="sm" @click="openClearAllConfirm"><Trash2 :size="17" />{{ t('client.inbox.clearAll') }}</Button>
      </div>
    </div>

    <div class="inbox-summary">
      <span><strong>{{ total }}</strong>{{ t('client.inbox.total') }}</span>
      <span class="inbox-summary__unread"><strong>{{ unreadCount }}</strong>{{ t('client.inbox.unread') }}</span>
    </div>

    <div class="inbox-tabs" role="tablist" :aria-label="t('client.inbox.filter')">
      <button v-for="tab in tabs" :key="tab.value" type="button" role="tab" :aria-selected="typeFilter === tab.value" :class="{ active: typeFilter === tab.value }" @click="changeType(tab.value)">
        <component :is="tab.icon" :size="18" />{{ tab.label }}
      </button>
    </div>

    <div v-if="loading && messages.length === 0" class="client-state" role="status">
      <span class="state-pulse" /><span>{{ t('client.inbox.loading') }}</span>
    </div>
    <div v-else-if="loadError" class="client-state client-state--error" role="alert">
      <CircleAlert :size="24" /><strong>{{ loadError }}</strong><Button variant="outline" size="sm" @click="loadInbox(1)">{{ t('client.inbox.retry') }}</Button>
    </div>
    <div v-else-if="messages.length === 0" class="client-state">
      <InboxIcon :size="42" /><strong>{{ t('client.inbox.emptyTitle') }}</strong><p>{{ t('client.inbox.emptyHint') }}</p>
    </div>

    <div v-else class="chat-list" aria-live="polite">
      <div v-for="item in filteredMessages" :key="messageKey(item)" class="chat-row" :class="item.sender_device_id === localID ? 'chat-row--sent' : 'chat-row--received'">
        <ChatBubble :message="chatMessage(item)" :show-avatar="false">
          <template v-if="item.type === 'FILE'">
            <a :href="fileURL(item)" class="chat-file" download @click="markRead(item)"><Download :size="16" />{{ t('client.inbox.download') }}</a>
          </template>
          <template v-else>{{ item.content }}</template>
        </ChatBubble>
        <div class="chat-tools">
          <Button v-if="item.type === 'TEXT'" variant="outline" size="sm" @click="copyText(item)"><Copy :size="15" />{{ t('client.inbox.copy') }}</Button>
          <button type="button" class="chat-delete" :aria-label="t('client.inbox.delete')" @click="openDeleteConfirm(item)"><Trash2 :size="15" /></button>
        </div>
      </div>
    </div>

    <Button v-if="hasMore" variant="outline" size="lg" class="load-more" :loading="loadingMore" @click="loadMore">{{ t('client.inbox.loadMore') }}</Button>

    <Teleport to="body">
      <div v-if="deleteConfirm.show" class="confirm-overlay" @click.self="closeDeleteConfirm">
        <div class="confirm-dialog">
          <strong>{{ t('client.inbox.deleteTitle') }}</strong>
          <p>{{ t('client.inbox.deleteConfirm') }}</p>
          <div class="confirm-actions">
            <Button variant="outline" size="sm" @click="closeDeleteConfirm">{{ t('client.inbox.cancel') }}</Button>
            <Button variant="danger" size="sm" @click="confirmDelete">{{ t('client.inbox.confirmDelete') }}</Button>
          </div>
        </div>
      </div>
      <div v-if="clearAllConfirm" class="confirm-overlay" @click.self="clearAllConfirm = false">
        <div class="confirm-dialog">
          <strong>{{ t('client.inbox.clearAllTitle') }}</strong>
          <p>{{ t('client.inbox.clearAllConfirm') }}</p>
          <div class="confirm-actions">
            <Button variant="outline" size="sm" @click="clearAllConfirm = false">{{ t('client.inbox.cancel') }}</Button>
            <Button variant="danger" size="sm" :loading="clearing" @click="confirmClearAll">{{ t('client.inbox.confirmClearAll') }}</Button>
          </div>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { CircleAlert, Copy, Download, File as FileIcon, Inbox as InboxIcon, MessageSquareText, RefreshCw, Trash2 } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { ChatBubble, type ChatMessage } from '@/components/ui/chat-bubble'
import { useToast } from '@/composables/useToast'
import { useMessageStore, type MessageItem } from '@/stores/message'
import { clearMessages, deleteMessage, downloadFileURL, getMessages, resolveApiErrorMessage, updateMessageStatus, type MessageListItem } from '@/api'

type InboxType = '' | 'TEXT' | 'FILE'
const { t } = useI18n()
const toast = useToast()
const store = useMessageStore()
const localID = (localStorage.getItem('device_id') || '').trim()
const typeFilter = ref<InboxType>('')
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const loadingMore = ref(false)
const loadError = ref('')
const messages = computed(() =>
    store.messages.filter((item) => item.sender_device_id === localID || item.target_device_id === localID)
)
const filteredMessages = computed(() => {
    const list = typeFilter.value ? messages.value.filter((item) => item.type === typeFilter.value) : [...messages.value]
    return list.sort((a, b) => String(a.created_at).localeCompare(String(b.created_at)))
})
const unreadCount = computed(() => store.messages.filter((item) => item.target_device_id === localID && item.status !== 'READ').length)
const hasMore = computed(() => messages.value.length < total.value)
const tabs = computed(() => [
  { value: '' as InboxType, label: t('client.inbox.all'), icon: InboxIcon },
  { value: 'TEXT' as InboxType, label: t('client.inbox.text'), icon: MessageSquareText },
  { value: 'FILE' as InboxType, label: t('client.inbox.file'), icon: FileIcon },
])
onMounted(() => void loadInbox(1))

function toMessage(item: MessageListItem): MessageItem {
  return { row_key: `${item.message_id}:${item.target_device_id || ''}`, id: item.message_id, type: item.type, content: item.content, file_id: item.file_id, sender_device_id: item.sender_device_id, target_device_id: item.target_device_id, status: item.status, created_at: item.created_at }
}
async function loadInbox(nextPage: number, append = false) {
  append ? loadingMore.value = true : loading.value = true
  loadError.value = ''
  try {
    const response = await getMessages({ page: nextPage, size: pageSize, device_id: localID })
    const data = response.data?.data || { items: [], total: 0 }
    const incoming = (Array.isArray(data.items) ? data.items : []).map((item: MessageListItem) => toMessage(item))
    if (append) {
      const merged = [...store.messages]
      for (const item of incoming) {
        const index = merged.findIndex((current) => messageKey(current) === messageKey(item))
        if (index >= 0) merged[index] = item
        else merged.push(item)
      }
      store.setMessages(merged)
    } else {
      store.setMessages(incoming)
    }
    total.value = Number(data.total || 0)
    page.value = nextPage
  } catch (error) { loadError.value = resolveApiErrorMessage(error, 'messagesPage.toast.loadFailed') }
  finally { loading.value = false; loadingMore.value = false }
}
function changeType(value: InboxType) { typeFilter.value = value }
function loadMore() { void loadInbox(page.value + 1, true) }
function messageKey(item: MessageItem) { return item.row_key || `${item.id}:${item.target_device_id || ''}` }
function fileURL(item: MessageItem) { return item.file_id ? downloadFileURL(item.file_id) : '#' }
function formatTime(value?: string) { if (!value) return '—'; const date = new Date(value); return Number.isNaN(date.getTime()) ? value : date.toLocaleString() }
function mapStatus(status?: string): ChatMessage['status'] {
  switch (status) {
    case 'READ': return 'read'
    case 'DELIVERED': return 'delivered'
    case 'CREATED': return 'sending'
    default: return undefined
  }
}
function chatMessage(item: MessageItem): ChatMessage {
  const sent = item.sender_device_id === localID
  return {
    id: messageKey(item),
    content: item.content || '',
    variant: sent ? 'sent' : 'received',
    name: sent ? t('client.inbox.me') : (item.sender_device_id || t('client.inbox.fromComputer')),
    timestamp: formatTime(item.created_at),
    status: sent ? mapStatus(item.status) : undefined,
  }
}
async function markRead(item: MessageItem) {
  if (!localID || item.status === 'READ') return
  try { await updateMessageStatus({ message_id: item.id, device_id: localID, status: 'READ' }); store.updateStatus(item.id, 'READ', localID) }
  catch (error) { toast.error(resolveApiErrorMessage(error, 'messagesPage.toast.markReadFailed')) }
}
function fallbackCopy(text: string): boolean {
  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.style.cssText = 'position:fixed;left:-9999px;top:-9999px;opacity:0'
  document.body.appendChild(textarea)
  textarea.select()
  const ok = document.execCommand('copy')
  document.body.removeChild(textarea)
  return ok
}
async function copyText(item: MessageItem) {
  const text = item.content || ''
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else if (!fallbackCopy(text)) {
      throw new Error('fallback failed')
    }
    await markRead(item)
    toast.success(t('client.inbox.copied'))
  } catch {
    if (fallbackCopy(text)) {
      await markRead(item)
      toast.success(t('client.inbox.copied'))
    } else {
      toast.error(t('client.inbox.copyFailed'))
    }
  }
}
const deleteConfirm = ref<{ show: boolean; item: MessageItem | null }>({ show: false, item: null })
function openDeleteConfirm(item: MessageItem) { deleteConfirm.value = { show: true, item } }
function closeDeleteConfirm() { deleteConfirm.value = { show: false, item: null } }
async function confirmDelete() {
  const item = deleteConfirm.value.item
  if (!item || !localID) return
  try {
    await deleteMessage({ message_id: item.id, device_id: localID })
    store.removeMessage(messageKey(item))
    toast.success(t('client.inbox.deleted'))
  } catch (error) {
    toast.error(resolveApiErrorMessage(error, 'client.inbox.deleteFailed'))
  }
  closeDeleteConfirm()
}
const clearAllConfirm = ref(false)
const clearing = ref(false)
function openClearAllConfirm() { clearAllConfirm.value = true }
async function confirmClearAll() {
  if (!localID) return
  clearing.value = true
  try {
    await clearMessages({ device_id: localID })
    store.clearMessages(localID)
    total.value = 0
    toast.success(t('client.inbox.clearedAll'))
  } catch (error) {
    toast.error(resolveApiErrorMessage(error, 'client.inbox.clearAllFailed'))
  }
  clearing.value = false
  clearAllConfirm.value = false
}
</script>

<style scoped>
.inbox-view { display: grid; gap: 16px; }
.client-page-head { display: flex; align-items: flex-end; justify-content: space-between; gap: 12px; }
.client-page-head--bare { justify-content: flex-end; }
.client-page-head p { margin: 0 0 4px; color: #9a4b1e; font-size: 9px; font-weight: 600; letter-spacing: .14em; }
.client-page-head h1 { margin: 0; font-size: 22px; line-height: 1.2; }
.client-page-head__actions { display: flex; gap: 8px; }
.inbox-summary { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.inbox-summary span { display: flex; align-items: baseline; gap: 6px; padding: 9px 12px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: #fff; font-size: 12px; font-weight: 600; }
.inbox-summary strong { font-size: 18px; }
.inbox-summary__unread { background: var(--brutal-accent) !important; }
.inbox-tabs { display: grid; grid-template-columns: repeat(3, 1fr); gap: 6px; padding: 4px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: var(--brutal-muted); }
.inbox-tabs button { display: flex; align-items: center; justify-content: center; gap: 5px; min-height: 40px; border: 2px solid transparent; border-radius: var(--brutal-radius); background: transparent; color: var(--brutal-fg); font-weight: 600; cursor: pointer; }
.inbox-tabs button.active { border-color: var(--brutal-border-color); background: var(--brutal-bg); box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
.inbox-tabs button:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 1px; }
.client-state { display: grid; min-height: 200px; place-items: center; align-content: center; gap: 10px; padding: 20px; border: 2px dashed var(--brutal-border-color); border-radius: var(--brutal-radius); background: var(--brutal-muted); text-align: center; }
.client-state p { margin: 0; color: var(--brutal-muted-foreground); }
.client-state--error { color: var(--brutal-destructive); }
.state-pulse { width: 38px; height: 38px; border: 4px solid var(--brutal-border-color); border-right-color: var(--brutal-primary); border-radius: 50%; animation: spin .8s linear infinite; }
.chat-list { display: flex; flex-direction: column; gap: 14px; }
.chat-row { display: flex; flex-direction: column; gap: 4px; max-width: 100%; }
.chat-row--sent { align-items: flex-end; }
.chat-row--received { align-items: flex-start; }
.chat-file { display: inline-flex; align-items: center; gap: 6px; color: inherit; font-weight: 600; text-decoration: underline; }
.chat-tools { display: flex; align-items: center; gap: 6px; padding: 0 2px; }
.chat-row--sent .chat-tools { justify-content: flex-end; }
.chat-delete { display: inline-flex; align-items: center; justify-content: center; width: 28px; height: 28px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: transparent; color: var(--brutal-destructive); cursor: pointer; }
.chat-delete:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 1px; }
.load-more { width: 100%; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (prefers-reduced-motion: reduce) { .state-pulse { animation: none; } }
.confirm-overlay { position: fixed; inset: 0; z-index: 1000; display: grid; place-items: center; background: rgba(0, 0, 0, 0.5); }
.confirm-dialog { width: 90%; max-width: 340px; padding: 20px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: var(--brutal-bg); box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.confirm-dialog strong { display: block; margin-bottom: 8px; font-size: 16px; }
.confirm-dialog p { margin: 0 0 20px; color: var(--brutal-muted-foreground); font-size: 14px; }
.confirm-actions { display: flex; gap: 10px; justify-content: flex-end; }
</style>
