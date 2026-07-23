package main

import (
 "encoding/json"
 "net/http"
 "strconv"
)

func (s *Server) HistoryHandler(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","application/json")

 page,_:=strconv.Atoi(r.URL.Query().Get("page"))
 size,_:=strconv.Atoi(r.URL.Query().Get("size"))
 if page<=0 {page=1}
 if size<=0 {size=20}

 list,err:=ListHistory(s.db,page,size)
 if err!=nil {
  http.Error(w,err.Error(),500)
  return
 }

 json.NewEncoder(w).Encode(list)
}
