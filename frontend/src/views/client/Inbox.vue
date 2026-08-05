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

    <section class="inbox-list-panel" :aria-busy="loading">
      <div v-if="loading" class="client-state" role="status">
        <span class="state-pulse" /><span>{{ t('client.inbox.loading') }}</span>
      </div>
      <div v-else-if="loadError" class="client-state client-state--error" role="alert">
        <CircleAlert :size="24" /><strong>{{ loadError }}</strong><Button variant="outline" size="sm" @click="loadInbox(page)">{{ t('client.inbox.retry') }}</Button>
      </div>
      <div v-else-if="messages.length === 0" class="client-state">
        <InboxIcon :size="42" /><strong>{{ t('client.inbox.emptyTitle') }}</strong><p>{{ t('client.inbox.emptyHint') }}</p>
      </div>

      <div v-else class="chat-list" aria-live="polite">
        <div
          v-for="item in filteredMessages"
          :key="messageKey(item)"
          class="swipe-wrapper"
          @touchstart.passive="onTouchStart(item, $event)"
          @touchmove="onTouchMove"
          @touchend="onTouchEnd(item)"
          @touchcancel="resetSwipe"
        >
          <Button variant="danger" size="sm" class="swipe-action swipe-action--delete" :aria-label="t('client.inbox.delete')" tabindex="-1" aria-hidden="true">
            <Trash2 :size="18" />{{ t('client.inbox.delete') }}
          </Button>
          <Button v-if="item.type === 'TEXT'" variant="outline" size="sm" class="swipe-action swipe-action--copy" :aria-label="t('client.inbox.copy')" tabindex="-1" aria-hidden="true">
            <Copy :size="18" />{{ t('client.inbox.copy') }}
          </Button>
          <div
            class="swipe-content"
            :class="{ 'swipe-content--dragging': isSwiping(messageKey(item)) }"
            :style="{ transform: `translateX(${swipeOffset(messageKey(item), item.type === 'TEXT')})` }"
          >
            <div class="chat-row" :class="item.sender_device_id === localID ? 'chat-row--sent' : 'chat-row--received'">
              <ChatBubble
                :message="chatMessage(item)"
                :show-avatar="false"
                :class="chatBubbleClass(item)"
              >
                <template v-if="item.type === 'FILE'">
                  <a :href="fileURL(item)" class="chat-file" download @click="markRead(item)"><Download :size="16" />{{ t('client.inbox.download') }}</a>
                </template>
                <template v-else>{{ item.content }}</template>
              </ChatBubble>
            </div>
          </div>
        </div>
      </div>

      <Pagination
        v-if="total > pageSize"
        :model-value="page"
        :total="total"
        :page-size="pageSize"
        layout="prev, pager, next"
        size="sm"
        :sibling-count="0"
        :show-first-last="false"
        :disabled="loading"
        class="inbox-pagination"
        @update:model-value="loadInbox"
      />
    </section>

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
import { ChatBubble, type ChatMessage } from 'brutx-ui-vue/chat-bubble'
import { Pagination } from 'brutx-ui-vue/pagination'
import { useToast } from '@/composables/useToast'
import { useMessageStore, type MessageItem } from '@/stores/message'
import { clearMessages, deleteMessage, downloadFileURL, getMessages, resolveApiErrorMessage, updateMessageStatus, type MessageListItem } from '@/api'

type InboxType = '' | 'TEXT' | 'FILE'
const { locale, t } = useI18n()
const toast = useToast()
const store = useMessageStore()
const localID = (localStorage.getItem('device_id') || '').trim()
const typeFilter = ref<InboxType>('')
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const loadError = ref('')
const swipeState = ref({ id: '', startX: 0, startY: 0, currentX: 0, dragging: false })
const messages = computed(() =>
    store.messages.filter((item) => item.sender_device_id === localID || item.target_device_id === localID)
)
const filteredMessages = computed(() => {
    const list = typeFilter.value ? messages.value.filter((item) => item.type === typeFilter.value) : [...messages.value]
    return list.sort((a, b) => messageTime(b.created_at) - messageTime(a.created_at))
})
const unreadCount = computed(() => store.messages.filter((item) => item.target_device_id === localID && item.status !== 'READ').length)
const tabs = computed(() => [
  { value: '' as InboxType, label: t('client.inbox.all'), icon: InboxIcon },
  { value: 'TEXT' as InboxType, label: t('client.inbox.text'), icon: MessageSquareText },
  { value: 'FILE' as InboxType, label: t('client.inbox.file'), icon: FileIcon },
])
onMounted(() => void loadInbox(1))

function toMessage(item: MessageListItem): MessageItem {
  return { row_key: `${item.message_id}:${item.target_device_id || ''}`, id: item.message_id, type: item.type, content: item.content, file_id: item.file_id, sender_device_id: item.sender_device_id, target_device_id: item.target_device_id, status: item.status, created_at: item.created_at }
}
async function loadInbox(nextPage: number) {
  loading.value = true
  loadError.value = ''
  try {
    const response = await getMessages({ page: nextPage, size: pageSize, device_id: localID, type: typeFilter.value })
    const data = response.data?.data || { items: [], total: 0 }
    const incoming = (Array.isArray(data.items) ? data.items : []).map((item: MessageListItem) => toMessage(item))
    store.setMessages(incoming)
    total.value = Number(data.total || 0)
    page.value = nextPage
  } catch (error) { loadError.value = resolveApiErrorMessage(error, 'messagesPage.toast.loadFailed') }
  finally { loading.value = false }
}
function changeType(value: InboxType) {
  if (typeFilter.value === value) return
  typeFilter.value = value
  void loadInbox(1)
}
function messageKey(item: MessageItem) { return item.row_key || `${item.id}:${item.target_device_id || ''}` }
function fileURL(item: MessageItem) { return item.file_id ? downloadFileURL(item.file_id) : '#' }
function formatTime(value?: string) {
  if (!value) return ''

  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) {
    return new Intl.DateTimeFormat(locale.value, {
      hour: '2-digit',
      minute: '2-digit',
      hourCycle: 'h23',
    }).format(date)
  }

  // Go time.String() may include a zone and monotonic suffix, which Date cannot parse.
  const match = value.match(/^\d{4}-\d{2}-\d{2}[T ](\d{2}):(\d{2})/)
  return match ? `${match[1]}:${match[2]}` : ''
}
function messageTime(value?: string) {
  if (!value) return 0
  const parsed = new Date(value).getTime()
  if (!Number.isNaN(parsed)) return parsed

  const normalized = value.match(/^\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?/i)?.[0]?.replace(' ', 'T')
  const fallback = normalized ? new Date(normalized).getTime() : Number.NaN
  return Number.isNaN(fallback) ? 0 : fallback
}
function chatBubbleClass(item: MessageItem) {
  return item.sender_device_id === localID
    ? 'inbox-chat-bubble inbox-chat-bubble--sent'
    : 'inbox-chat-bubble inbox-chat-bubble--received'
}
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
function onTouchStart(item: MessageItem, event: TouchEvent) {
  const touch = event.touches[0]
  if (!touch) return
  swipeState.value = {
    id: messageKey(item),
    startX: touch.clientX,
    startY: touch.clientY,
    currentX: touch.clientX,
    dragging: false,
  }
}
function onTouchMove(event: TouchEvent) {
  if (!swipeState.value.id) return
  const touch = event.touches[0]
  if (!touch) return
  const dx = touch.clientX - swipeState.value.startX
  const dy = touch.clientY - swipeState.value.startY
  if (!swipeState.value.dragging && Math.abs(dx) > Math.abs(dy) && Math.abs(dx) > 10) {
    swipeState.value.dragging = true
  }
  if (swipeState.value.dragging) {
    event.preventDefault()
    swipeState.value.currentX = touch.clientX
  }
}
function onTouchEnd(item: MessageItem) {
  const dx = swipeState.value.currentX - swipeState.value.startX
  if (swipeState.value.dragging && dx > 80) openDeleteConfirm(item)
  else if (swipeState.value.dragging && dx < -80 && item.type === 'TEXT') void copyText(item)
  resetSwipe()
}
function resetSwipe() {
  swipeState.value = { id: '', startX: 0, startY: 0, currentX: 0, dragging: false }
}
function isSwiping(key: string) {
  return swipeState.value.id === key && swipeState.value.dragging
}
function swipeOffset(key: string, canCopy: boolean) {
  if (!isSwiping(key)) return '0px'
  const dx = swipeState.value.currentX - swipeState.value.startX
  return `${Math.max(canCopy ? -112 : 0, Math.min(112, dx))}px`
}
const deleteConfirm = ref<{ show: boolean; item: MessageItem | null }>({ show: false, item: null })
function openDeleteConfirm(item: MessageItem) { deleteConfirm.value = { show: true, item } }
function closeDeleteConfirm() { deleteConfirm.value = { show: false, item: null } }
async function confirmDelete() {
  const item = deleteConfirm.value.item
  if (!item || !localID) return
  try {
    await deleteMessage({ message_id: item.id, device_id: localID })
    const remainingTotal = Math.max(0, total.value - 1)
    const lastPage = Math.max(1, Math.ceil(remainingTotal / pageSize))
    await loadInbox(Math.min(page.value, lastPage))
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
.inbox-list-panel { display: grid; gap: 16px; padding: 16px; border: 3px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: #fff; box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.chat-list { display: flex; flex-direction: column; gap: 14px; }
.swipe-wrapper { position: relative; overflow: hidden; border-radius: var(--brutal-radius); touch-action: pan-y; }
.swipe-content { position: relative; z-index: 1; background: #fff; transition: transform 150ms ease-out; will-change: transform; }
.swipe-content--dragging { transition: none; }
.swipe-action { position: absolute; z-index: 0; top: 0; bottom: 0; width: 104px; height: 100%; border-radius: 0; box-shadow: none; pointer-events: none; }
.swipe-action--delete { left: 0; }
.swipe-action--copy { right: 0; background: var(--brutal-accent); }
.chat-row { display: flex; flex-direction: column; gap: 4px; max-width: 100%; }
.chat-row--sent { align-items: flex-end; }
.chat-row--received { align-items: flex-start; }
:deep(.inbox-chat-bubble) {
  min-width: 72px;
  padding: 12px 16px;
  overflow-wrap: anywhere;
}
:deep(.inbox-chat-bubble--sent) { background: var(--brutal-primary); }
:deep(.inbox-chat-bubble--received) { background: var(--brutal-muted); }
.chat-file { display: inline-flex; align-items: center; gap: 6px; color: inherit; font-weight: 600; text-decoration: underline; }
.inbox-pagination { width: 100%; flex-wrap: wrap; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (prefers-reduced-motion: reduce) { .state-pulse { animation: none; } .swipe-content { transition: none; } }
.confirm-overlay { position: fixed; inset: 0; z-index: 1000; display: grid; place-items: center; background: rgba(0, 0, 0, 0.5); }
.confirm-dialog { width: 90%; max-width: 340px; padding: 20px; border: 2px solid var(--brutal-border-color); border-radius: var(--brutal-radius); background: var(--brutal-bg); box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.confirm-dialog strong { display: block; margin-bottom: 8px; font-size: 16px; }
.confirm-dialog p { margin: 0 0 20px; color: var(--brutal-muted-foreground); font-size: 14px; }
.confirm-actions { display: flex; gap: 10px; justify-content: flex-end; }
</style>
