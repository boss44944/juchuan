<template>
  <el-card>
    <h2>发送</h2>
    <el-select v-model="targets" multiple placeholder="选择设备">
      <el-option
        v-for="device in devices"
        :key="device.id"
        :label="device.display_name"
        :value="device.id"
      />
    </el-select>

    <el-input v-model="content" type="textarea" placeholder="输入文字" />
    <el-button @click="sendText">发送文字</el-button>

    <el-upload :auto-upload="false" :on-change="selectFile">
      <el-button>选择文件</el-button>
    </el-upload>
    <el-button @click="sendFile">发送文件</el-button>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDeviceStore } from '../stores/device'
import { sendTextMessage, uploadFile, sendFileMessage } from '../api'

const store = useDeviceStore()
const devices = store.devices
const targets = ref<string[]>([])
const content = ref('')
const file = ref<File | null>(null)

function selectFile(upload:any){
  file.value = upload.raw
}

async function sendText(){
  await sendTextMessage({content: content.value, targets: targets.value})
}

async function sendFile(){
  if(!file.value) return
  const result:any = await uploadFile(file.value)
  await sendFileMessage({file_id: result.data.file_id, targets: targets.value})
}
</script>
