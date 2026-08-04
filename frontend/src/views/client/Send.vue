<template>
  <section class="client-send" aria-labelledby="client-send-title">
    <header class="client-page-head"><p>JUCHUAN / SEND</p><h1 id="client-send-title">{{ t('client.send.title') }}</h1><span>{{ t('client.send.target') }}</span></header>

    <div class="send-tabs" role="tablist" :aria-label="t('client.send.chooseType')">
      <button type="button" role="tab" :aria-selected="mode === 'text'" :class="{ active: mode === 'text' }" @click="mode = 'text'"><MessageSquareText :size="20" />{{ t('client.inbox.text') }}</button>
      <button type="button" role="tab" :aria-selected="mode === 'file'" :class="{ active: mode === 'file' }" @click="mode = 'file'"><FileUp :size="20" />{{ t('client.inbox.file') }}</button>
    </div>

    <Card v-if="mode === 'text'" padding="lg" class="send-panel">
      <div class="send-panel__title"><span><MessageSquareText :size="23" /></span><div><h2>{{ t('send.sendText') }}</h2><p>{{ t('client.send.textHint') }}</p></div></div>
      <Textarea v-model="content" size="lg" class="client-textarea" :placeholder="t('send.inputText')" :aria-label="t('send.inputText')" />
      <div class="character-count">{{ content.length }}</div>
      <Button variant="primary" size="lg" class="send-action" :loading="sending" :disabled="sending" @click="sendText"><SendIcon :size="20" />{{ t('client.send.sendToComputer') }}</Button>
    </Card>

    <Card v-else padding="lg" class="send-panel">
      <div class="send-panel__title"><span class="file-mark"><FileUp :size="23" /></span><div><h2>{{ t('send.sendFile') }}</h2><p>{{ t('client.send.fileHint') }}</p></div></div>
      <input ref="fileInput" class="visually-hidden" type="file" @change="selectFile" />
      <button type="button" class="file-picker" @click="fileInput?.click()">
        <FileIcon :size="36" /><span><strong>{{ file ? file.name : t('send.selectFile') }}</strong><small>{{ file ? formatFileSize(file.size) : t('client.send.fileEmpty') }}</small></span>
      </button>
      <Button variant="secondary" size="lg" class="send-action" :loading="sending" :disabled="sending" @click="sendFile"><Upload :size="20" />{{ t('client.send.sendToComputer') }}</Button>
    </Card>

    <p class="privacy-note"><ShieldCheck :size="17" />{{ t('client.send.localOnly') }}</p>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { File as FileIcon, FileUp, MessageSquareText, Send as SendIcon, ShieldCheck, Upload } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import { useToast } from '@/composables/useToast'
import { useMessageStore, type MessageItem } from '@/stores/message'
import { resolveApiErrorMessage, sendFileMessage, sendTextMessage, uploadFile } from '@/api'

const { t } = useI18n()
const toast = useToast()
const messageStore = useMessageStore()
const mode = ref<'text' | 'file'>('text')
const content = ref('')
const file = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const sending = ref(false)

function asRecord(value: unknown): Record<string, unknown> | null { return typeof value === 'object' && value !== null ? value as Record<string, unknown> : null }
function extractMessage(responseData: unknown): MessageItem | null {
  const outer = asRecord(responseData); const nested = asRecord(outer?.data); const message = nested ?? outer
  if (!message || typeof message.id !== 'string') return null
  return { ...(message as Partial<MessageItem>), id: message.id, type: typeof message.type === 'string' ? message.type : '', status: 'CREATED' }
}
function selectFile(event: Event) { const input = event.target as HTMLInputElement; file.value = input.files?.[0] ?? null }
function formatFileSize(bytes: number) { if (bytes < 1024) return `${bytes} B`; if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`; return `${(bytes / (1024 * 1024)).toFixed(1)} MB` }
async function sendText() {
  if (!content.value.trim()) { toast.warning(t('send.toast.textRequired')); return }
  sending.value = true
  try {
    const response = await sendTextMessage({ content: content.value, sender_device_id: localStorage.getItem('device_id') || '', targets: ['server'] })
    const message = extractMessage(response.data); if (message) messageStore.addMessage(message)
    content.value = ''; toast.success(t('send.toast.textSent'))
  } catch (error) { toast.error(resolveApiErrorMessage(error)) }
  finally { sending.value = false }
}
async function sendFile() {
  if (!file.value) { toast.warning(t('send.toast.fileRequired')); return }
  sending.value = true
  try {
    const uploadResponse = await uploadFile(file.value)
    const uploadOuter = asRecord(uploadResponse.data); const uploadData = asRecord(uploadOuter?.data) ?? uploadOuter
    const fileID = typeof uploadData?.file_id === 'string' ? uploadData.file_id : typeof uploadData?.id === 'string' ? uploadData.id : ''
    if (!fileID) throw new Error('missing file id')
    const response = await sendFileMessage({ file_id: fileID, sender_device_id: localStorage.getItem('device_id') || '', targets: ['server'] })
    const message = extractMessage(response.data); if (message) messageStore.addMessage(message)
    file.value = null; if (fileInput.value) fileInput.value.value = ''; toast.success(t('send.toast.fileSent'))
  } catch (error) { toast.error(resolveApiErrorMessage(error)) }
  finally { sending.value = false }
}
</script>

<style scoped>
.client-send { display: grid; gap: 17px; }
.client-page-head p { margin: 0 0 4px; color: #9a4b1e; font-size: 9px; font-weight: 900; letter-spacing: .14em; }
.client-page-head h1 { margin: 0; font-size: 30px; line-height: 1; }
.client-page-head span { display: inline-block; margin-top: 9px; padding: 5px 8px; border: 2px solid var(--brutal-border-color); border-radius: 99px; background: var(--brutal-accent); font-size: 11px; font-weight: 900; }
.send-tabs { display: grid; grid-template-columns: 1fr 1fr; gap: 6px; padding: 5px; border: 2px solid var(--brutal-border-color); border-radius: 7px; background: var(--brutal-muted); }
.send-tabs button { display: flex; align-items: center; justify-content: center; gap: 7px; min-height: 48px; border: 2px solid transparent; border-radius: 5px; background: transparent; color: var(--brutal-fg); font-weight: 900; cursor: pointer; }
.send-tabs button.active { border-color: var(--brutal-border-color); background: var(--brutal-bg); box-shadow: 2px 2px 0 var(--brutal-shadow-color); }
.send-tabs button:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 1px; }
.send-panel { display: grid; gap: 17px; box-shadow: 4px 4px 0 var(--brutal-shadow-color); }
.send-panel__title { display: flex; align-items: center; gap: 12px; }
.send-panel__title > span { display: grid; flex: 0 0 44px; width: 44px; height: 44px; place-items: center; border: 3px solid var(--brutal-border-color); background: var(--brutal-primary); box-shadow: 3px 3px 0 var(--brutal-shadow-color); }
.send-panel__title > .file-mark { background: var(--brutal-secondary); color: #fff; }
.send-panel h2 { margin: 0; font-size: 20px; }
.send-panel p { margin: 3px 0 0; color: var(--brutal-muted-foreground); font-size: 12px; line-height: 1.4; }
:deep(.client-textarea) { min-height: 210px; resize: vertical; }
.character-count { margin-top: -11px; color: var(--brutal-muted-foreground); font-size: 11px; font-weight: 900; text-align: right; }
.file-picker { display: grid; min-height: 220px; width: 100%; place-items: center; align-content: center; gap: 13px; padding: 20px; border: 3px dashed var(--brutal-border-color); border-radius: 7px; background: var(--brutal-muted); color: var(--brutal-fg); text-align: center; cursor: pointer; }
.file-picker span { display: grid; gap: 4px; max-width: 100%; }
.file-picker strong { overflow-wrap: anywhere; }
.file-picker small { color: var(--brutal-muted-foreground); }
.file-picker:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 2px; }
.send-action { width: 100%; min-height: 52px; }
.privacy-note { display: flex; align-items: center; justify-content: center; gap: 7px; margin: 0; color: var(--brutal-muted-foreground); font-size: 11px; font-weight: 800; text-align: center; }
</style>
