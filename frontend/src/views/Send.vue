<template>
  <el-card>
    <h2>{{ t('send.title') }}</h2>
    <el-select v-model="targets" multiple :placeholder="t('send.selectTargets')">
      <el-option
        v-for="device in devices"
        :key="device.id"
        :label="device.display_name"
        :value="device.id"
      />
    </el-select>

    <el-input v-model="content" type="textarea" :placeholder="t('send.inputText')" />
    <el-button @click="sendText">{{ t('send.sendText') }}</el-button>

    <el-upload :auto-upload="false" :on-change="selectFile">
      <el-button>{{ t('send.selectFile') }}</el-button>
    </el-upload>
    <el-button @click="sendFile">{{ t('send.sendFile') }}</el-button>
  </el-card>
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
