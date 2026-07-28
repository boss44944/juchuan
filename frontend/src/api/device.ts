import request from './request'

export interface RegisterDeviceRequest {
  id: string
  display_name: string
  platform?: string
  browser?: string
}

export function registerDevice(data: RegisterDeviceRequest) {
  return request.post('/api/device/register', data)
}

export function checkDeviceName(display_name: string) {
  return request.get('/api/device/check-name', {
    params: { display_name }
  })
}

export function getDevices() {
  return request.get('/api/devices')
}
