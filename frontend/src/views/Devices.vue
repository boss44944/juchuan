<template>
  <el-table :data="devices">
    <el-table-column prop="display_name" label="Device" />
    <el-table-column prop="platform" label="Platform" />
    <el-table-column prop="status" label="Status" />
  </el-table>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useDeviceStore } from '../stores/device'

const store = useDeviceStore()
const devices = computed(() => store.devices)

onMounted(() => {
  // device loading will be connected with api layer
})
</script>
