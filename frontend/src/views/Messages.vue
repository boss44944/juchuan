<template>
  <el-card v-for="item in messages" :key="item.id">
    <div v-if="item.type === 'TEXT'">{{ item.content }}</div>
    <div v-else-if="item.type === 'FILE'">
      <a :href="`/api/file/download/${item.file_id}`">Download</a>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useMessageStore } from '../stores/message'

const store = useMessageStore()
const messages = computed(() => store.messages)
</script>
