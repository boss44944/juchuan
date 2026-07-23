package main

import "net/http"

func DownloadHandler(w http.ResponseWriter,r *http.Request){
 http.ServeFile(w,r,r.URL.Path)
}
