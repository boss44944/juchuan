import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
})

export interface LoginRequest {
  device_id: string
  password?: string
}

export function authStatus() {
  return api.get('/auth/status')
}

export function login(data: LoginRequest) {
  return api.post('/auth/login', data)
}

export function logout() {
  return api.post('/auth/logout')
}
