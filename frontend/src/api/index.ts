import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

export function getDevices() {
  return api.get('/devices')
}

export function sendText(data: any) {
  return api.post('/message/text', data)
}

export function sendFileMessage(data: any) {
  return api.post('/message/file', data)
}

export function uploadFile(data: FormData) {
  return api.post('/file/upload', data, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function updateMessageStatus(data: any) {
  return api.post('/message/status', data)
}

export function login(data: any) {
  return api.post('/auth/login', data)
}

export function getConfig() {
  return api.get('/config')
}

export function updateConfig(data: any) {
  return api.post('/config', data)
}

export default api
