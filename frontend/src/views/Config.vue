<template>
  <el-card>
    <h2>配置</h2>
    <el-form>
      <el-form-item label="服务端口">
        <el-input v-model="config.port" />
      </el-form-item>
      <el-form-item label="自动打开">
        <el-switch v-model="config.auto_open" />
      </el-form-item>
      <el-button @click="save">保存</el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { getConfig, updateConfig } from '../api'

const config = reactive<any>({
  port: 8000,
  auto_open: false,
})

onMounted(async () => {
  const res = await getConfig()
  Object.assign(config, res.data)
})

async function save() {
  await updateConfig(config)
}
</script>
