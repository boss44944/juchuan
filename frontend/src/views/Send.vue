<template>
  <section class="view-panel send-view">
    <header class="panel-header">
      <div>
        <h2 class="panel-title">{{ t('menu.send') }}</h2>
        <p class="panel-subtitle">LOW-LATENCY PAYLOAD DELIVERY</p>
      </div>
    </header>

    <el-select v-model="targets" multiple :placeholder="t('send.selectTargets')" class="target-select">
      <el-option
        v-for="device in devices"
        :key="device.id"
        :label="device.display_name"
        :value="device.id"
      />
    </el-select>

    <div class="send-grid">
      <el-card shadow="never" class="composer-card">
        <h3>{{ t('send.sendText') }}</h3>
        <el-input v-model="content" type="textarea" :rows="8" :placeholder="t('send.inputText')" />
        <el-button type="primary" class="composer-btn" @click="sendText">{{ t('send.sendText') }}</el-button>
      </el-card>

      <el-card shadow="never" class="composer-card">
        <h3>{{ t('send.sendFile') }}</h3>
        <el-upload :auto-upload="false" :on-change="selectFile" class="upload-box">
          <el-button>{{ t('send.selectFile') }}</el-button>
        </el-upload>
        <div class="file-name" v-if="file">{{ file.name }}</div>
        <div class="file-name" v-else>-</div>
        <el-button type="primary" class="composer-btn" @click="sendFile">{{ t('send.sendFile') }}</el-button>
      </el-card>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useDeviceStore } from '../stores/device'
import { useMessageStore } from '../stores/message'
import { sendTextMessage, uploadFile, sendFileMessage, resolveApiErrorMessage } from '../api'

const { t } = useI18n()
const store = useDeviceStore()
const messageStore = useMessageStore()
const devices = computed(() => store.devices)
const targets = ref<string[]>([])
const content = ref('')
const file = ref<File | null>(null)

onMounted(async () => {
  await store.load()
})

function selectFile(upload:any){
  file.value = upload.raw
}

async function sendText(){
  if (!content.value.trim()) {
    ElMessage.warning(t('send.toast.textRequired'))
    return
  }
  if (targets.value.length === 0) {
    ElMessage.warning(t('send.toast.targetRequired'))
    return
  }

  try {
    const senderID = localStorage.getItem('device_id') || ''
    const res: any = await sendTextMessage({
      content: content.value,
      sender_device_id: senderID,
      targets: targets.value,
    })
    const msg = res.data?.data || res.data
    if (msg?.id) {
      messageStore.addMessage({
        ...msg,
        status: 'CREATED',
      })
    }
    content.value = ''
    ElMessage.success(t('send.toast.textSent'))
  } catch (err) {
    ElMessage.error(resolveApiErrorMessage(err))
  }
}

async function sendFile(){
  if (!file.value) {
    ElMessage.warning(t('send.toast.fileRequired'))
    return
  }
  if (targets.value.length === 0) {
    ElMessage.warning(t('send.toast.targetRequired'))
    return
  }

  try {
    const result: any = await uploadFile(file.value)
    const fileId = result.data?.data?.file_id || result.data?.file_id
    const senderID = localStorage.getItem('device_id') || ''
    const res: any = await sendFileMessage({
      file_id: fileId,
      sender_device_id: senderID,
      targets: targets.value,
    })
    const msg = res.data?.data || res.data
    if (msg?.id) {
      messageStore.addMessage({
        ...msg,
        status: 'CREATED',
      })
    }
    file.value = null
    ElMessage.success(t('send.toast.fileSent'))
  } catch (err) {
    ElMessage.error(resolveApiErrorMessage(err))
  }
}
</script>

<style scoped>
.target-select {
  width: 100%;
  margin-bottom: 16px;
}

.send-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.composer-card h3 {
  margin: 0 0 12px;
  font-size: 16px;
  color: #dff8ff;
}

.upload-box {
  margin-bottom: 14px;
}

.file-name {
  min-height: 30px;
  display: flex;
  align-items: center;
  color: #a5d5dd;
  font-size: 13px;
  margin-bottom: 16px;
}

.composer-btn {
  width: 100%;
}

@media (max-width: 980px) {
  .send-grid {
    grid-template-columns: 1fr;
  }
}
</style>
