<template>
  <el-card>
    <h2>{{ t('menu.config') }}</h2>
    <el-form>
      <el-form-item label="服务端口">
        <el-input v-model="config.port" />
      </el-form-item>
      <el-form-item label="自动打开">
        <el-switch v-model="config.auto_open" />
      </el-form-item>
      <el-form-item label="访问密码">
        <el-input v-model="config.password" type="password" show-password />
      </el-form-item>
      <el-form-item label="语言">
        <el-select v-model="config.language">
          <el-option label="中文" value="zh-CN" />
          <el-option label="English" value="en-US" />
          <el-option label="日本語" value="ja-JP" />
        </el-select>
      </el-form-item>
      <el-button @click="save">保存</el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getConfig, updateConfig } from '../api'

const { t } = useI18n()

const config = reactive<any>({
  port: 8000,
  auto_open: false,
  password: '',
  language: 'zh-CN',
})

onMounted(async () => {
  const res = await getConfig()
  Object.assign(config, res.data)
})

async function save() {
  await updateConfig(config)
}
</script>
