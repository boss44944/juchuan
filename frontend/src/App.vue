<template>
  <div class="app-root theme-juchuan">
    <ToastHost />
    <router-view />
  </div>
</template>

<script setup lang="ts">
import ToastHost from '@/components/app/ToastHost.vue'
</script>

<style scoped>
.app-root { min-height: 100vh; }
</style>
