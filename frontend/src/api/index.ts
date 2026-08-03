import axios from 'axios'
import i18n from '../i18n'

export const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
})

function tByKey(key: string, params?: Record<string, unknown>) {
  if (i18n.global.te(key)) {
    return i18n.global.t(key, params)
  }
  return i18n.global.t('error.UNKNOWN')
}

function inferErrorCodeFromStatus(err: any): string {
  const status = err?.response?.status
  const url = String(err?.config?.url || '')
  const text = String(err?.response?.data || '').toLowerCase()

  if (status === 401) {
    if (url.includes('/auth/login') || text.includes('wrong password')) {
      return 'AUTH_PASSWORD_INVALID'
    }
    return 'AUTH_REQUIRED'
  }
  if (status === 404) {
    if (url.includes('/file/download') || url.includes('/download/')) {
      return 'FILE_NOT_FOUND'
    }
    return 'NOT_FOUND'
  }
  if (status === 409) {
    return 'CONFLICT'
  }
  if (status === 400) {
    return 'INVALID_REQUEST'
  }
  if (status >= 500) {
    return 'SERVER_ERROR'
  }
  return 'UNKNOWN'
}

export function resolveApiErrorMessage(err: unknown, fallbackKey = 'error.UNKNOWN') {
  if (!axios.isAxiosError(err)) {
    return i18n.global.t(fallbackKey)
  }

  const payload = err.response?.data as any
  const apiCode = payload?.error?.code || payload?.code
  const params = payload?.error?.params || undefined
  const code = typeof apiCode === 'string' && apiCode.trim() !== ''
    ? apiCode.trim().toUpperCase()
    : inferErrorCodeFromStatus(err)

  const errorKey = `error.${code}`
  if (i18n.global.te(errorKey)) {
    return tByKey(errorKey, params)
  }

  if (typeof payload?.error === 'string' && payload.error.trim() !== '') {
    return payload.error
  }
  if (typeof payload === 'string' && payload.trim() !== '') {
    return payload
  }

  return i18n.global.t(fallbackKey)
}

export interface LoginPayload {
  device_id?: string
  password?: string
}

export function authStatus() {
  return api.get('/auth/status')
}

export function login(data: LoginPayload) {
  return api.post('/auth/login', data)
}

export function logout() {
  return api.post('/auth/logout')
}

export function getDevices() {
  return api.get('/devices')
}

export function registerDevice(data: {
  id: string
  display_name: string
  role?: string
  platform?: string
  browser?: string
}) {
  return api.post('/device/register', data)
}

export function heartbeatDevice(id: string) {
  return api.post('/device/heartbeat', { id })
}

export function renameDevice(data: { id: string; display_name: string }) {
  return api.post('/device/rename', data)
}

export function removeDevice(data: { id: string }) {
  return api.post('/device/remove', data)
}

export function sendTextMessage(data: {
  content: string
  sender_device_id: string
  targets: string[]
}) {
  return api.post('/message/text', data)
}

export function sendFileMessage(data: {
  file_id: string
  sender_device_id: string
  targets: string[]
}) {
  return api.post('/message/file', data)
}

export function uploadFile(file: File) {
  const form = new FormData()
  form.append('file', file)
  return api.post('/file/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function downloadFileURL(fileId: string) {
  return `/api/file/download/${fileId}`
}

export function updateMessageStatus(data: {
  message_id: string
  device_id: string
  status: 'CREATED' | 'DELIVERED' | 'READ'
}) {
  return api.post('/message/status', data)
}

export interface MessageQuery {
  page?: number
  size?: number
  type?: '' | 'TEXT' | 'FILE'
  status?: '' | 'CREATED' | 'DELIVERED' | 'READ'
  device_id?: string
  sender_device_id?: string
  target_device_id?: string
}

export interface MessageListItem {
  message_id: string
  type: 'TEXT' | 'FILE'
  content?: string
  file_id?: string
  sender_device_id: string
  target_device_id?: string
  status: 'CREATED' | 'DELIVERED' | 'READ'
  created_at: string
}

export function getMessages(params: MessageQuery) {
  return api.get('/messages', { params })
}

export function getHistory(page = 1, size = 50) {
  return api.get('/history', {
    params: { page, size },
  })
}

export function getConfig() {
  return api.get('/config')
}

export function updateConfig(data: {
  port?: number
  auto_open?: boolean
  password?: string
}) {
  return api.post('/config', data)
}

export function qrCodeURL(urlText?: string) {
  const params = new URLSearchParams()
  if (urlText && urlText.trim() !== '') {
    params.set('url', urlText.trim())
  }
  params.set('t', String(Date.now()))
  return `/api/qr?${params.toString()}`
}

export default api
