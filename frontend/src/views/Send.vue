<template>
  <section class="send-view theme-juchuan" aria-labelledby="send-page-title">
    <header class="send-hero">
      <div>
        <p class="send-eyebrow">JUCHUAN / LOCAL TRANSFER</p>
        <h2 id="send-page-title">{{ t('menu.send') }}</h2>
        <p class="send-subtitle">{{ t('send.subtitle') }}</p>
      </div>
      <Badge variant="accent" size="sm" class="selection-count">
        {{ t('send.selectedCount', { selected: targets.length, total: devices.length }) }}
      </Badge>
    </header>

    <section class="target-panel" aria-labelledby="target-panel-title">
      <div class="section-heading">
        <div>
          <h3 id="target-panel-title">{{ t('send.selectTargets') }}</h3>
          <p>{{ t('send.targetHint') }}</p>
        </div>
        <span class="step-marker" aria-hidden="true">01</span>
      </div>

      <div v-if="loadingDevices" class="target-state" role="status" aria-live="polite">
        <span class="state-block" />
        <span class="state-block state-block--wide" />
      </div>
      <div v-else-if="loadFailed" class="target-state target-state--error" role="alert">
        {{ t('error.UNKNOWN') }}
      </div>
      <div v-else-if="devices.length === 0" class="target-state">
        {{ t('send.noDevices') }}
      </div>
      <div v-else class="target-grid">
        <label
          v-for="device in devices"
          :key="device.id"
          class="target-option"
          :class="{ 'target-option--selected': targets.includes(device.id) }"
        >
          <input v-model="targets" type="checkbox" :value="device.id" />
          <span class="target-check" aria-hidden="true"><Check :size="16" :stroke-width="3" /></span>
          <span class="target-copy">
            <strong>{{ device.display_name }}</strong>
            <small>{{ device.platform || device.id }}</small>
          </span>
          <Badge
            :variant="device.status === 'online' ? 'success' : 'outline'"
            size="sm"
            dot
          >
            <Wifi v-if="device.status === 'online'" :size="13" aria-hidden="true" />
            <WifiOff v-else :size="13" aria-hidden="true" />
            {{ t(`devices.status.${device.status === 'online' ? 'online' : 'offline'}`) }}
          </Badge>
        </label>
      </div>
    </section>

    <div class="send-grid">
      <Card variant="default" padding="lg" class="composer-card">
        <div class="composer-heading">
          <span class="composer-icon" aria-hidden="true"><MessageSquareText :size="24" /></span>
          <div>
            <CardTitle>{{ t('send.sendText') }}</CardTitle>
            <CardDescription>{{ t('send.textHint') }}</CardDescription>
          </div>
          <span class="step-marker" aria-hidden="true">02</span>
        </div>

        <Textarea
          v-model="content"
          size="lg"
          class="message-textarea"
          :placeholder="t('send.inputText')"
          :aria-label="t('send.inputText')"
        />
        <div class="composer-meta" aria-live="polite">{{ content.length }}</div>

        <Button
          variant="primary"
          size="lg"
          class="composer-button"
          :loading="sendingText"
          :disabled="sendingText || sendingFile"
          @click="sendText"
        >
          <SendIcon :size="20" aria-hidden="true" />
          {{ t('send.sendText') }}
        </Button>
      </Card>

      <Card variant="default" padding="lg" class="composer-card">
        <div class="composer-heading">
          <span class="composer-icon composer-icon--orange" aria-hidden="true"><FileUp :size="24" /></span>
          <div>
            <CardTitle>{{ t('send.sendFile') }}</CardTitle>
            <CardDescription>{{ t('send.fileHint') }}</CardDescription>
          </div>
          <span class="step-marker" aria-hidden="true">03</span>
        </div>

        <input ref="fileInput" class="visually-hidden" type="file" @change="selectFile" />
        <button type="button" class="file-picker" @click="openFilePicker">
          <FileIcon :size="30" aria-hidden="true" />
          <span>
            <strong>{{ file ? file.name : t('send.selectFile') }}</strong>
            <small>{{ file ? formatFileSize(file.size) : t('send.fileEmpty') }}</small>
          </span>
        </button>

        <Button
          variant="secondary"
          size="lg"
          class="composer-button"
          :loading="sendingFile"
          :disabled="sendingText || sendingFile"
          @click="sendFile"
        >
          <Upload :size="20" aria-hidden="true" />
          {{ t('send.sendFile') }}
        </Button>
      </Card>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Check, File as FileIcon, FileUp, MessageSquareText, Send as SendIcon, Upload, Wifi, WifiOff } from '@lucide/vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardDescription, CardTitle } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'
import { useToast } from '@/composables/useToast'
import { useDeviceStore } from '../stores/device'
import { useMessageStore, type MessageItem } from '../stores/message'
import { resolveApiErrorMessage, sendFileMessage, sendTextMessage, uploadFile } from '../api'

const { t } = useI18n()
const toast = useToast()
const store = useDeviceStore()
const messageStore = useMessageStore()
const devices = computed(() => store.devices)
const targets = ref<string[]>([])
const content = ref('')
const file = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const loadingDevices = ref(true)
const loadFailed = ref(false)
const sendingText = ref(false)
const sendingFile = ref(false)

onMounted(async () => {
  try {
    await store.load()
  } catch (err) {
    loadFailed.value = true
    toast.error(resolveApiErrorMessage(err))
  } finally {
    loadingDevices.value = false
  }
})

function openFilePicker() {
  fileInput.value?.click()
}

function selectFile(event: Event) {
  const input = event.target as HTMLInputElement
  file.value = input.files?.[0] ?? null
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function asRecord(value: unknown): Record<string, unknown> | null {
  return typeof value === 'object' && value !== null ? value as Record<string, unknown> : null
}

function extractMessage(responseData: unknown): MessageItem | null {
  const outer = asRecord(responseData)
  const nested = asRecord(outer?.data)
  const message = nested ?? outer
  if (!message || typeof message.id !== 'string') return null

  return {
    ...(message as Partial<MessageItem>),
    id: message.id,
    type: typeof message.type === 'string' ? message.type : '',
    status: 'CREATED',
  }
}

async function sendText() {
  if (!content.value.trim()) {
    toast.warning(t('send.toast.textRequired'))
    return
  }
  if (targets.value.length === 0) {
    toast.warning(t('send.toast.targetRequired'))
    return
  }

  sendingText.value = true
  try {
    const senderID = localStorage.getItem('device_id') || ''
    const response = await sendTextMessage({
      content: content.value,
      sender_device_id: senderID,
      targets: targets.value,
    })
    const message = extractMessage(response.data)
    if (message) messageStore.addMessage(message)
    content.value = ''
    toast.success(t('send.toast.textSent'))
  } catch (err) {
    toast.error(resolveApiErrorMessage(err))
  } finally {
    sendingText.value = false
  }
}

async function sendFile() {
  if (!file.value) {
    toast.warning(t('send.toast.fileRequired'))
    return
  }
  if (targets.value.length === 0) {
    toast.warning(t('send.toast.targetRequired'))
    return
  }

  sendingFile.value = true
  try {
    const uploadResponse = await uploadFile(file.value)
    const uploadOuter = asRecord(uploadResponse.data)
    const uploadData = asRecord(uploadOuter?.data) ?? uploadOuter
    const fileId = typeof uploadData?.file_id === 'string' ? uploadData.file_id : ''
    const senderID = localStorage.getItem('device_id') || ''
    const response = await sendFileMessage({
      file_id: fileId,
      sender_device_id: senderID,
      targets: targets.value,
    })
    const message = extractMessage(response.data)
    if (message) messageStore.addMessage(message)
    file.value = null
    if (fileInput.value) fileInput.value.value = ''
    toast.success(t('send.toast.fileSent'))
  } catch (err) {
    toast.error(resolveApiErrorMessage(err))
  } finally {
    sendingFile.value = false
  }
}
</script>

<style scoped>
.theme-juchuan {
  --brutal-border-width: 3px;
  --brutal-border-color: #4b2b15;
  --brutal-shadow-color: #4b2b15;
  --brutal-shadow-offset-x: 5px;
  --brutal-shadow-offset-y: 5px;
  --brutal-radius: 7px;
  --brutal-bg: #fff8e7;
  --brutal-fg: #2b1a0f;
  --brutal-primary: #f3b63f;
  --brutal-primary-foreground: #2b1a0f;
  --brutal-secondary: #e8722a;
  --brutal-secondary-foreground: #fff8e7;
  --brutal-accent: #ffd86e;
  --brutal-accent-foreground: #2b1a0f;
  --brutal-muted: #f5ead2;
  --brutal-muted-foreground: #6c5139;
  --brutal-success: #86a95b;
  --brutal-success-foreground: #1d2a12;
  --brutal-destructive: #c9483b;
  --brutal-destructive-foreground: #fff8e7;
  --brutal-info: #397b9c;
  --brutal-info-foreground: #fff8e7;
  --brutal-ring: #2b1a0f;
  --brutal-placeholder: #8b7258;
  --brutal-pressed-offset: 2px;
}

.send-view {
  color: var(--brutal-fg);
  background: #fff8e7;
  border: 3px solid var(--brutal-border-color);
  border-radius: 9px;
  box-shadow: 8px 8px 0 var(--brutal-shadow-color);
  padding: 24px;
}

.send-hero,
.section-heading,
.composer-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.send-hero {
  padding: 4px 4px 22px;
  border-bottom: 3px solid var(--brutal-border-color);
}

.send-eyebrow {
  margin: 0 0 6px;
  color: #9a4b1e;
  font-size: 11px;
  font-weight: 900;
  letter-spacing: 0.16em;
}

.send-hero h2 {
  margin: 0;
  font-size: clamp(28px, 4vw, 44px);
  line-height: 1;
  font-weight: 900;
  letter-spacing: -0.035em;
}

.send-subtitle,
.section-heading p {
  margin: 8px 0 0;
  color: var(--brutal-muted-foreground);
  line-height: 1.55;
}

.selection-count {
  flex: 0 0 auto;
  white-space: nowrap;
}

.target-panel {
  margin: 24px 0;
  padding: 20px;
  border: 3px solid var(--brutal-border-color);
  border-radius: var(--brutal-radius);
  background: var(--brutal-muted);
  box-shadow: 5px 5px 0 var(--brutal-secondary);
}

.section-heading h3,
.composer-heading h3 {
  margin: 0;
  color: var(--brutal-fg);
  font-size: 18px;
  font-weight: 900;
}

.step-marker {
  color: #9a4b1e;
  font-size: 13px;
  font-weight: 900;
  letter-spacing: 0.08em;
}

.target-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
  gap: 12px;
  margin-top: 18px;
}

.target-option {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  padding: 12px;
  cursor: pointer;
  background: var(--brutal-bg);
  border: 3px solid var(--brutal-border-color);
  border-radius: var(--brutal-radius);
  box-shadow: 3px 3px 0 var(--brutal-shadow-color);
  transition: transform 150ms ease, box-shadow 150ms ease, background 150ms ease;
}

.target-option:hover {
  transform: translate(-1px, -1px);
  box-shadow: 5px 5px 0 var(--brutal-shadow-color);
}

.target-option:has(input:focus-visible) {
  outline: 3px solid var(--brutal-ring);
  outline-offset: 3px;
}

.target-option--selected {
  background: var(--brutal-accent);
}

.target-option input {
  position: absolute;
  width: 1px;
  height: 1px;
  opacity: 0;
}

.target-check {
  display: grid;
  flex: 0 0 24px;
  width: 24px;
  height: 24px;
  place-items: center;
  color: transparent;
  background: #fff;
  border: 3px solid var(--brutal-border-color);
}

.target-option--selected .target-check {
  color: var(--brutal-fg);
  background: var(--brutal-primary);
}

.target-copy {
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
}

.target-copy strong,
.target-copy small {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.target-copy small {
  margin-top: 3px;
  color: var(--brutal-muted-foreground);
  font-size: 11px;
}

.target-state {
  display: flex;
  min-height: 72px;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-top: 18px;
  padding: 16px;
  color: var(--brutal-muted-foreground);
  font-weight: 700;
  border: 3px dashed var(--brutal-border-color);
  background: rgba(255, 248, 231, 0.65);
}

.target-state--error {
  color: var(--brutal-destructive);
}

.state-block {
  width: 56px;
  height: 22px;
  background: #dfcda9;
  animation: pulse 1s ease-in-out infinite alternate;
}

.state-block--wide {
  width: 140px;
}

.send-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.composer-card {
  display: flex;
  min-width: 0;
  min-height: 390px;
  flex-direction: column;
}

.composer-heading {
  align-items: center;
  margin-bottom: 22px;
}

.composer-heading > div {
  flex: 1;
  min-width: 0;
}

.composer-icon {
  display: grid;
  flex: 0 0 46px;
  width: 46px;
  height: 46px;
  place-items: center;
  color: var(--brutal-fg);
  border: 3px solid var(--brutal-border-color);
  background: var(--brutal-primary);
  box-shadow: 3px 3px 0 var(--brutal-shadow-color);
}

.composer-icon--orange {
  color: var(--brutal-secondary-foreground);
  background: var(--brutal-secondary);
}

:deep(.message-textarea) {
  min-height: 180px;
  resize: vertical;
}

.composer-meta {
  min-height: 22px;
  margin: 7px 2px 16px;
  color: var(--brutal-muted-foreground);
  font-size: 12px;
  font-weight: 800;
  text-align: right;
}

.composer-button {
  width: 100%;
  margin-top: auto;
}

.file-picker {
  display: flex;
  width: 100%;
  min-height: 180px;
  align-items: center;
  justify-content: center;
  gap: 14px;
  margin: 0 0 22px;
  padding: 20px;
  cursor: pointer;
  color: var(--brutal-fg);
  text-align: left;
  background: var(--brutal-muted);
  border: 3px dashed var(--brutal-border-color);
  border-radius: var(--brutal-radius);
  transition: transform 150ms ease, box-shadow 150ms ease, background 150ms ease;
}

.file-picker:hover {
  transform: translate(-1px, -1px);
  background: var(--brutal-accent);
  box-shadow: 5px 5px 0 var(--brutal-shadow-color);
}

.file-picker:active {
  transform: translate(2px, 2px);
  box-shadow: none;
}

.file-picker:focus-visible {
  outline: 3px solid var(--brutal-ring);
  outline-offset: 3px;
}

.file-picker span {
  display: flex;
  min-width: 0;
  flex-direction: column;
}

.file-picker strong,
.file-picker small {
  overflow-wrap: anywhere;
}

.file-picker small {
  margin-top: 5px;
  color: var(--brutal-muted-foreground);
}

.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

@keyframes pulse {
  from { opacity: 0.45; }
  to { opacity: 0.9; }
}

@media (max-width: 980px) {
  .send-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 620px) {
  .send-view {
    padding: 14px;
    border-width: 2px;
    box-shadow: 5px 5px 0 var(--brutal-shadow-color);
  }

  .send-hero {
    align-items: flex-start;
    flex-direction: column;
    padding-bottom: 18px;
  }

  .selection-count {
    align-self: flex-start;
  }

  .target-panel {
    margin: 18px 0;
    padding: 14px;
  }

  .target-grid {
    grid-template-columns: 1fr;
  }

  .composer-card {
    min-height: 360px;
    padding: 18px;
  }

  .composer-heading {
    gap: 10px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .target-option,
  .file-picker,
  .state-block {
    transition: none;
    animation: none;
  }
}
</style>
