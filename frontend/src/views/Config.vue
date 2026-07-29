<template>
  <el-card>
    <h2>{{ t('menu.config') }}</h2>
    <el-form>
      <el-form-item :label="t('configPage.labels.port')">
        <el-input v-model="config.port" />
      </el-form-item>
      <el-form-item :label="t('configPage.labels.autoOpen')">
        <el-switch v-model="config.auto_open" />
      </el-form-item>
      <el-form-item :label="t('configPage.labels.password')">
        <el-input v-model="config.password" type="password" show-password />
      </el-form-item>
      <el-form-item :label="t('configPage.labels.language')">
        <el-select v-model="config.language">
          <el-option :label="t('configPage.languages.zhCN')" value="zh-CN" />
          <el-option :label="t('configPage.languages.enUS')" value="en-US" />
          <el-option :label="t('configPage.languages.jaJP')" value="ja-JP" />
        </el-select>
      </el-form-item>
      <el-button @click="save">{{ t('configPage.save') }}</el-button>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getConfig, resolveApiErrorMessage, updateConfig } from '../api'

const { t } = useI18n()

const config = reactive<any>({
  port: 8000,
  auto_open: false,
  password: '',
  language: 'zh-CN',
})

onMounted(async () => {
  try {
    const res = await getConfig()
    Object.assign(config, res.data)
  } catch (err) {
    ElMessage.error(resolveApiErrorMessage(err))
  }
})

async function save() {
  try {
    await updateConfig({
      port: Number(config.port),
      auto_open: !!config.auto_open,
      password: config.password,
    })
    ElMessage.success(t('configPage.saved'))
  } catch (err) {
    ElMessage.error(resolveApiErrorMessage(err))
  }
}
</script>
