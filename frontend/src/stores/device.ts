import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DeviceIdentity } from '../types/device'

export const useDeviceStore = defineStore('device', () => {
  const devices = ref<DeviceIdentity[]>([])

  const current = ref<DeviceIdentity>({
    id: '',
    display_name: ''
  })

  function setCurrent(device: DeviceIdentity) {
    current.value = device
  }

  function setDevices(list: DeviceIdentity[]) {
    devices.value = list
  }

  function addDevice(device: DeviceIdentity) {
    const index = devices.value.findIndex(d => d.id === device.id)
    if (index >= 0) {
      devices.value[index] = device
    } else {
      devices.value.push(device)
    }
  }

  function removeDevice(id: string) {
    devices.value = devices.value.filter(d => d.id !== id)
  }

  return {
    devices,
    current,
    setCurrent,
    setDevices,
    addDevice,
    removeDevice
  }
})
