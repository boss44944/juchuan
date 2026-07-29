import axios from 'axios'

const api = axios.create({
  baseURL: '/api'
})

export interface LoginRequest {
  device_id: string
  password?: string
}

export function login(data: LoginRequest) {
  return api.post('/auth/login', data)
}
