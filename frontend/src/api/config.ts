import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
})

export function getConfig(){
  return api.get('/config')
}

export function saveConfig(data:any){
  return api.post('/config', data)
}

export const updateConfig = saveConfig
