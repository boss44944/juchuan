package main

import (
 "encoding/json"
 "os"
 "path/filepath"
)

type Config struct {
 Port int `json:"port"`
 AutoOpen bool `json:"auto_open"`
 Password string `json:"password"`
 MaxUploadMB int64 `json:"max_upload_mb"`
}

func defaultConfig() Config { return Config{Port:8000,AutoOpen:true,MaxUploadMB:512} }

func configFile() string {
 d,_:=os.UserConfigDir()
 return filepath.Join(d,"Juchuan","config.json")
}

func loadConfig() Config {
 c:=defaultConfig()
 b,e:=os.ReadFile(configFile())
 if e==nil { _=json.Unmarshal(b,&c) }
 return c
}

func saveConfig(c Config) error {
 p:=configFile()
 _=os.MkdirAll(filepath.Dir(p),0755)
 b,_:=json.MarshalIndent(c,"","  ")
 return os.WriteFile(p,b,0644)
}
