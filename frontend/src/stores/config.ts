import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useConfigStore = defineStore('config', () => {
  const config = ref<Record<string, any>>({})

  function setConfig(value:Record<string,any>) {
    config.value = value
  }

  return { config, setConfig }
})
