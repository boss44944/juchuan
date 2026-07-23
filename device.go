package main

import "sync"

type Device struct {
 ID string
 Name string
 Token string
 LastSeen int64
}

type DeviceManager struct {
 mu sync.Mutex
 devices map[string]*Device
}

func NewDeviceManager()*DeviceManager{
 return &DeviceManager{devices:make(map[string]*Device)}
}

func (m *DeviceManager) Add(d *Device){
 m.mu.Lock()
 defer m.mu.Unlock()
 m.devices[d.ID]=d
}

func (m *DeviceManager) List() []*Device{
 m.mu.Lock()
 defer m.mu.Unlock()
 out:=make([]*Device,0,len(m.devices))
 for _,d:=range m.devices { out=append(out,d) }
 return out
}
