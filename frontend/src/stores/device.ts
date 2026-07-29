import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DeviceIdentity } from '../types/device'
import { getDevices } from '../api'

export const useDeviceStore = defineStore('device', () => {
  const devices = ref<DeviceIdentity[]>([])

  const current = ref<DeviceIdentity>({
    id: '',
    display_name: ''
  })

  async function load() {
    const res = await getDevices()
    const payload = Array.isArray(res.data) ? res.data : (res.data.data || [])
    setDevices(payload)
  }

  function setCurrent(device: DeviceIdentity) {
    current.value = device
  }

  function setDevices(list: DeviceIdentity[]) {
    devices.value = [...list].sort(sortByOnline)
  }

  function addDevice(device: DeviceIdentity) {
    const index = devices.value.findIndex(d => d.id === device.id)
    if (index >= 0) {
      devices.value[index] = device
    } else {
      devices.value.push(device)
    }
    devices.value.sort(sortByOnline)
  }

  function removeDevice(id: string) {
    devices.value = devices.value.filter(d => d.id !== id)
  }

  function handleEvent(event: any) {
    if (event.type === 'DEVICE_ONLINE') {
      addDevice(event.data)
    }

    if (event.type === 'DEVICE_OFFLINE') {
      const device = devices.value.find(d => d.id === event.data.id)
      if (device) {
        device.status = 'offline'
        devices.value.sort(sortByOnline)
      }
    }
  }

  function sortByOnline(a: DeviceIdentity, b: DeviceIdentity) {
    const as = a.status === 'online' ? 1 : 0
    const bs = b.status === 'online' ? 1 : 0
    if (as !== bs) {
      return bs - as
    }
    return a.display_name.localeCompare(b.display_name)
  }

  return {
    devices,
    current,
    load,
    setCurrent,
    setDevices,
    addDevice,
    removeDevice,
    handleEvent
  }
})
