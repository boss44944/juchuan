import { onMounted } from 'vue'
import { heartbeatDevice, registerDevice } from '@/api'
import { useDeviceStore } from '@/stores/device'
import { useMessageStore } from '@/stores/message'
import { connectWebSocket, disconnectWebSocket, onWebSocketMessage } from '@/websocket/client'
import type { AccessRole } from '@/utils/role'

let wsBound = false
let heartbeatTimer: number | null = null

export function useAppRuntime(role: AccessRole) {
  const deviceStore = useDeviceStore()
  const messageStore = useMessageStore()

  onMounted(async () => {
    if (role === 'server') {
      await deviceStore.load().catch(() => undefined)
    }

    const id = (localStorage.getItem('device_id') || '').trim()
    if (!id) return

    await registerDevice({
      id,
      display_name: id,
      role,
      platform: navigator.platform,
      browser: navigator.userAgent,
    }).catch(() => undefined)

    if (heartbeatTimer == null) {
      heartbeatTimer = window.setInterval(() => void heartbeatDevice(id), 30000)
    }
    if (!wsBound) {
      onWebSocketMessage((event) => {
        deviceStore.handleEvent(event)
        messageStore.handleEvent(event)
      })
      wsBound = true
    }
    connectWebSocket()
  })
}

export function stopAppRuntime() {
  if (heartbeatTimer != null) {
    window.clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
  disconnectWebSocket()
}
