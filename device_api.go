package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type DeviceRequest struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	Platform    string `json:"platform"`
	Browser     string `json:"browser"`
}

func (s *Server) DeviceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { WriteError(w,http.StatusBadRequest,"INVALID_REQUEST",nil); return }
	if req.ID == "" || req.DisplayName == "" { WriteError(w,http.StatusBadRequest,"DEVICE_INFO_REQUIRED",nil); return }
	if s.devices.NameExists(req.DisplayName) { WriteError(w,http.StatusConflict,"DEVICE_NAME_EXISTS",map[string]interface{}{"name":req.DisplayName}); return }
	if req.Role == "" { req.Role="client" }
	device:=&Device{ID:req.ID,DisplayName:req.DisplayName,Role:req.Role,Platform:req.Platform,Browser:req.Browser,LastSeen:time.Now().Unix(),Status:DeviceStatusOnline}
	s.devices.Add(device)
	if s.db != nil { _ = SaveDevice(s.db, device) }
	if s.events != nil { s.events.Publish(DeviceEvent{Type:DeviceOnlineEvent,Data:*device}) }
	WriteJSON(w,http.StatusOK,APIResponse{Success:true,Data:map[string]string{"device_id":req.ID}})
}

func (s *Server) DevicesHandler(w http.ResponseWriter,r *http.Request){ WriteJSON(w,http.StatusOK,APIResponse{Success:true,Data:s.devices.List()}) }

func (s *Server) DeviceHeartbeatHandler(w http.ResponseWriter,r *http.Request){
	var req struct{ ID string `json:"id"` }
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil { WriteError(w,http.StatusBadRequest,"INVALID_REQUEST",nil); return }
	device,ok:=s.devices.Get(req.ID); if !ok { WriteError(w,http.StatusNotFound,"DEVICE_NOT_FOUND",nil); return }
	wasOffline:=device.Status==DeviceStatusOffline
	device.LastSeen=time.Now().Unix(); device.Status=DeviceStatusOnline
	if s.db!=nil { _=SaveDevice(s.db,device) }
	if wasOffline && s.events!=nil { s.events.Publish(DeviceEvent{Type:DeviceOnlineEvent,Data:*device}) }
	WriteJSON(w,http.StatusOK,APIResponse{Success:true})
}

func (s *Server) DeviceRenameHandler(w http.ResponseWriter,r *http.Request){
	var req struct{ ID string `json:"id"`; DisplayName string `json:"display_name"` }
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil { WriteError(w,http.StatusBadRequest,"INVALID_REQUEST",nil); return }
	device,ok:=s.devices.Get(req.ID); if !ok { WriteError(w,http.StatusNotFound,"DEVICE_NOT_FOUND",nil); return }
	if req.DisplayName=="" { WriteError(w,http.StatusBadRequest,"DEVICE_NAME_REQUIRED",nil); return }
	if device.DisplayName!=req.DisplayName && s.devices.NameExists(req.DisplayName){ WriteError(w,http.StatusConflict,"DEVICE_NAME_EXISTS",map[string]interface{}{"name":req.DisplayName}); return }
	device.DisplayName=req.DisplayName
	if s.db!=nil { _=UpdateDeviceName(s.db,req.ID,req.DisplayName) }
	WriteJSON(w,http.StatusOK,APIResponse{Success:true,Data:device})
}

func (s *Server) DeviceRemoveHandler(w http.ResponseWriter,r *http.Request){
	var req struct{ ID string `json:"id"` }
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil { WriteError(w,http.StatusBadRequest,"INVALID_REQUEST",nil); return }
	if !s.devices.Remove(req.ID){ WriteError(w,http.StatusNotFound,"DEVICE_NOT_FOUND",nil); return }
	if s.db!=nil { _=DeleteDevice(s.db,req.ID) }
	if s.events!=nil { s.events.Publish(DeviceEvent{Type:"DEVICE_REMOVED",Data:Device{ID:req.ID}}) }
	WriteJSON(w,http.StatusOK,APIResponse{Success:true})
}
