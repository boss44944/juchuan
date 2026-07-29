import { defineStore } from 'pinia'
import { ref } from 'vue'
import { updateMessageStatus } from '../api'

export interface MessageItem {
  row_key?: string
  id: string
  type: string
  content?: string
  file_id?: string
  history_id?: string
  sender_device_id?: string
  target_device_id?: string
  created_at?: string
  status?: 'CREATED' | 'DELIVERED' | 'READ'
}

export const useMessageStore = defineStore('message', () => {
  const messages = ref<MessageItem[]>([])

  function addMessage(message: MessageItem) {
    const key = message.row_key || `${message.id}:${message.target_device_id || ''}`
    const index = messages.value.findIndex((item) => (item.row_key || `${item.id}:${item.target_device_id || ''}`) === key)
    if (index >= 0) {
      messages.value[index] = { ...messages.value[index], ...message, row_key: key }
      return
    }
    messages.value.unshift({ ...message, row_key: key })
  }

  function setMessages(list: MessageItem[]) {
    messages.value = [...list]
  }

  function updateStatus(messageId: string, status: 'CREATED' | 'DELIVERED' | 'READ', deviceId?: string) {
    for (const msg of messages.value) {
      if (msg.id !== messageId) {
        continue
      }
      if (deviceId && msg.target_device_id && msg.target_device_id !== deviceId) {
        continue
      }
      msg.status = status
    }
  }

  function handleEvent(event: any) {
    if (event.type === 'MESSAGE_RECEIVED') {
      const data = event.data || {}
      const localDeviceID = localStorage.getItem('device_id') || ''
      const targetID = localDeviceID || data.target_device_id || ''
      addMessage({
        id: data.id,
        type: data.type,
        content: data.content,
        file_id: data.file_id,
        sender_device_id: data.sender_device_id,
        target_device_id: targetID,
        created_at: data.created_at,
        row_key: `${data.id}:${targetID}`,
        ...data,
        status: 'DELIVERED',
      })

      if (localDeviceID && data.id && data.sender_device_id !== localDeviceID) {
        void updateMessageStatus({
          message_id: data.id,
          device_id: localDeviceID,
          status: 'READ',
        })
      }
    }

    if (event.type === 'MESSAGE_STATUS_UPDATED' && event.data) {
      updateStatus(event.data.message_id, event.data.status, event.data.device_id)
    }
  }

  return {
    messages,
    addMessage,
    setMessages,
    updateStatus,
    handleEvent
  }
})
