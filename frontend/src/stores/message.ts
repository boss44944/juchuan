import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface MessageItem {
  id: string
  type: string
  content?: string
  file_id?: string
}

export const useMessageStore = defineStore('message', () => {
  const messages = ref<MessageItem[]>([])

  function addMessage(message: MessageItem) {
    messages.value.unshift(message)
  }

  function handleEvent(event: any) {
    if (event.type === 'MESSAGE_RECEIVED') {
      addMessage(event.data)
    }
  }

  return {
    messages,
    addMessage,
    handleEvent
  }
})
