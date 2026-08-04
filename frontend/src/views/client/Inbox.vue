<template>
  <section class="inbox-view" aria-labelledby="client-inbox-title">
    <header class="client-page-head">
      <div><p>JUCHUAN / INBOX</p><h1 id="client-inbox-title">{{ t('client.inbox.title') }}</h1></div>
      <Button variant="outline" size="sm" :loading="loading" @click="loadInbox(1)"><RefreshCw :size="17" />{{ t('client.inbox.refresh') }}</Button>
    </header>

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

    <div v-else class="message-list" aria-live="polite">
      <Card v-for="item in filteredMessages" :key="messageKey(item)" padding="default" :class="item.status !== 'READ' ? 'mobile-message mobile-message--unread' : 'mobile-message'">
        <div class="mobile-message__top">
          <span class="message-type" :class="`message-type--${item.type.toLowerCase()}`">
            <MessageSquareText v-if="item.type === 'TEXT'" :size="20" /><FileIcon v-else :size="20" />
          </span>
          <div class="mobile-message__meta">
            <strong>{{ item.type === 'TEXT' ? t('client.inbox.textMessage') : t('client.inbox.fileMessage') }}</strong>
            <small>{{ formatTime(item.created_at) }}</small>
          </div>
          <Badge :variant="item.status === 'READ' ? 'success' : 'accent'" size="sm" dot>{{ statusLabel(item.status) }}</Badge>
        </div>

        <p v-if="item.type === 'TEXT'" class="message-copy">{{ item.content || '—' }}</p>
        <div v-else class="file-copy"><FileIcon :size="28" /><span><strong>{{ t('client.inbox.receivedFile') }}</strong><small>{{ shortFileID(item.file_id) }}</small></span></div>

        <div class="mobile-message__from"><Monitor :size="15" /><span>{{ t('client.inbox.fromComputer') }}</span><code>{{ item.sender_device_id || 'server' }}</code></div>

        <Button v-if="item.type === 'TEXT'" variant="primary" size="default" class="message-action" @click="copyText(item)"><Copy :size="19" />{{ t('client.inbox.copy') }}</Button>
        <a v-else :href="fileURL(item)" class="download-action" download @click="markRead(item)"><Download :size="19" />{{ t('client.inbox.download') }}</a>
      </Card>
    </div>

    <Button v-if="hasMore" variant="outline" size="lg" class="load-more" :loading="loadingMore" @click="loadMore">{{ t('client.inbox.loadMore') }}</Button>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { CircleAlert, Copy, Download, File as FileIcon, Inbox as InboxIcon, MessageSquareText, Monitor, RefreshCw } from '@lucide/vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { useToast } from '@/composables/useToast'
import { useMessageStore, type MessageItem } from '@/stores/message'
import { downloadFileURL, getMessages, resolveApiErrorMessage, updateMessageStatus, type MessageListItem } from '@/api'

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
const messages = computed(() => store.messages.filter((item) => item.target_device_id === localID))
const filteredMessages = computed(() => typeFilter.value ? messages.value.filter((item) => item.type === typeFilter.value) : messages.value)
const unreadCount = computed(() => messages.value.filter((item) => item.status !== 'READ').length)
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
    const response = await getMessages({ page: nextPage, size: pageSize, target_device_id: localID })
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
function shortFileID(value?: string) { return value ? `${value.slice(0, 8)}…` : '—' }
function formatTime(value?: string) { if (!value) return '—'; const date = new Date(value); return Number.isNaN(date.getTime()) ? value : date.toLocaleString() }
function statusLabel(status?: string) { return status === 'READ' ? t('messagesPage.status.read') : t('messagesPage.status.delivered') }
async function markRead(item: MessageItem) {
  if (!localID || item.status === 'READ') return
  try { await updateMessageStatus({ message_id: item.id, device_id: localID, status: 'READ' }); store.updateStatus(item.id, 'READ', localID) }
  catch (error) { toast.error(resolveApiErrorMessage(error, 'messagesPage.toast.markReadFailed')) }
}
async function copyText(item: MessageItem) {
  try { await navigator.clipboard.writeText(item.content || ''); await markRead(item); toast.success(t('client.inbox.copied')) }
  catch { toast.error(t('client.inbox.copyFailed')) }
}
</script>

<style scoped>
.inbox-view { display: grid; gap: 16px; }
.client-page-head { display: flex; align-items: flex-end; justify-content: space-between; gap: 12px; }
.client-page-head p { margin: 0 0 4px; color: #9a4b1e; font-size: 9px; font-weight: 900; letter-spacing: .14em; }
.client-page-head h1 { margin: 0; font-size: 30px; line-height: 1; }
.inbox-summary { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.inbox-summary span { display: flex; align-items: baseline; gap: 6px; padding: 11px 13px; border: 2px solid var(--brutal-border-color); border-radius: 6px; background: #fff; font-size: 12px; font-weight: 900; }
.inbox-summary strong { font-size: 23px; }
.inbox-summary__unread { background: var(--brutal-accent) !important; }
.inbox-tabs { display: grid; grid-template-columns: repeat(3, 1fr); gap: 6px; padding: 5px; border: 2px solid var(--brutal-border-color); border-radius: 7px; background: var(--brutal-muted); }
.inbox-tabs button { display: flex; align-items: center; justify-content: center; gap: 5px; min-height: 44px; border: 2px solid transparent; border-radius: 5px; background: transparent; color: var(--brutal-fg); font-weight: 900; cursor: pointer; }
.inbox-tabs button.active { border-color: var(--brutal-border-color); background: var(--brutal-bg); box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
.inbox-tabs button:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 1px; }
.client-state { display: grid; min-height: 240px; place-items: center; align-content: center; gap: 10px; padding: 28px; border: 3px dashed var(--brutal-border-color); border-radius: 7px; background: var(--brutal-muted); text-align: center; }
.client-state p { margin: 0; color: var(--brutal-muted-foreground); }
.client-state--error { color: var(--brutal-destructive); }
.state-pulse { width: 42px; height: 42px; border: 4px solid var(--brutal-border-color); border-right-color: var(--brutal-primary); border-radius: 50%; animation: spin .8s linear infinite; }
.message-list { display: grid; gap: 14px; }
.mobile-message { position: relative; display: grid; gap: 13px; box-shadow: 3px 3px 0 var(--brutal-shadow-color); }
.mobile-message--unread::before { position: absolute; top: -3px; bottom: -3px; left: -3px; width: 7px; border: 2px solid var(--brutal-border-color); background: var(--brutal-secondary); content: ''; }
.mobile-message__top { display: flex; align-items: center; gap: 10px; }
.message-type { display: grid; flex: 0 0 40px; width: 40px; height: 40px; place-items: center; border: 2px solid var(--brutal-border-color); background: var(--brutal-primary); box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
.message-type--file { background: var(--brutal-secondary); color: #fff; }
.mobile-message__meta { display: grid; flex: 1; min-width: 0; }
.mobile-message__meta small { margin-top: 2px; color: var(--brutal-muted-foreground); font-size: 11px; }
.message-copy { max-height: 180px; margin: 0; padding: 12px; border: 2px solid var(--brutal-border-color); background: var(--brutal-muted); line-height: 1.55; overflow: auto; overflow-wrap: anywhere; white-space: pre-wrap; }
.file-copy { display: flex; align-items: center; gap: 11px; padding: 13px; border: 2px solid var(--brutal-border-color); background: var(--brutal-muted); }
.file-copy span { display: grid; min-width: 0; }
.file-copy small { margin-top: 3px; color: var(--brutal-muted-foreground); }
.mobile-message__from { display: flex; align-items: center; gap: 5px; min-width: 0; color: var(--brutal-muted-foreground); font-size: 11px; }
.mobile-message__from code { margin-left: auto; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.message-action, .download-action { width: 100%; min-height: 48px; }
.download-action { display: flex; align-items: center; justify-content: center; gap: 8px; border: 3px solid var(--brutal-border-color); border-radius: 6px; background: var(--brutal-secondary); box-shadow: 4px 4px 0 var(--brutal-shadow-color); color: #fff; font-weight: 900; text-decoration: none; }
.download-action:active { box-shadow: none; transform: translate(3px, 3px); }
.download-action:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 2px; }
.load-more { width: 100%; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (prefers-reduced-motion: reduce) { .state-pulse { animation: none; } }
</style>
