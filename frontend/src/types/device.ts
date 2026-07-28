export interface DeviceIdentity {
  id: string
  display_name: string
  platform?: string
  browser?: string
  status?: 'online' | 'offline'
  last_seen?: string
}
