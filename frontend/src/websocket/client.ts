export type WSHandler = (message: any) => void

let socket: WebSocket | null = null
const handlers: WSHandler[] = []
let reconnectTimer: number | null = null

function getDeviceId() {
  const value = localStorage.getItem('device_id')
  return value ? value.trim() : ''
}

export function connectWebSocket() {
  if (socket) return
  const deviceId = getDeviceId()
  if (!deviceId) return

  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  socket = new WebSocket(`${protocol}//${location.host}/ws?device=${encodeURIComponent(deviceId)}`)

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data)
    handlers.forEach((handler) => handler(data))
  }

  socket.onclose = () => {
    socket = null
    if (reconnectTimer == null) {
      reconnectTimer = window.setTimeout(() => {
        reconnectTimer = null
        connectWebSocket()
      }, 1500)
    }
  }
}

export function onWebSocketMessage(handler: WSHandler) {
  handlers.push(handler)
}

export function sendWebSocket(data: any) {
  socket?.send(JSON.stringify(data))
}
