import { api } from './index'

export interface RegisterDeviceRequest {
  id: string
  display_name: string
  platform?: string
  browser?: string
}

export function registerDevice(data: RegisterDeviceRequest) {
  return api.post('/device/register', data)
}

export function checkDeviceName(display_name: string) {
  return api.get('/devices', {
    params: { display_name }
  })
}

export function getDevices() {
  return api.get('/devices')
}
