export type WSHandler = (message: any) => void

let socket: WebSocket | null = null
const handlers: WSHandler[] = []

export function connectWebSocket() {
  if (socket) return

  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  socket = new WebSocket(`${protocol}//${location.host}/ws`)

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data)
    handlers.forEach((handler) => handler(data))
  }

  socket.onclose = () => {
    socket = null
  }
}

export function onWebSocketMessage(handler: WSHandler) {
  handlers.push(handler)
}

export function sendWebSocket(data: any) {
  socket?.send(JSON.stringify(data))
}
