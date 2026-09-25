package main

import (
 "encoding/json"
 "log"
 "net/http"
 "os"
 "github.com/updateguard/updateguard/internal/domain"
 "github.com/updateguard/updateguard/internal/store"
)
func main(){
 db:=store.NewMemory()
 mux:=http.NewServeMux()
 mux.HandleFunc("/health",func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(http.StatusOK);_,_=w.Write([]byte("ok\n"))})
 mux.HandleFunc("/ready",func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(http.StatusOK);_,_=w.Write([]byte("ready\n"))})
 mux.HandleFunc("/api/v1/dashboard",func(w http.ResponseWriter,_ *http.Request){write(w,map[string]int{"services":0,"up_to_date":0,"updates_available":0,"failed_updates":0})})
 mux.HandleFunc("/api/v1/audit",func(w http.ResponseWriter,r *http.Request){write(w,db.Audit(r.Context()))})
 mux.HandleFunc("/api/v1/update-states",func(w http.ResponseWriter,_ *http.Request){write(w,[]domain.UpdateState{domain.Pending,domain.Precheck,domain.BackingUp,domain.Pulling,domain.Deploying,domain.Verifying,domain.Stabilizing,domain.Success,domain.Failed,domain.RollingBack,domain.RolledBack,domain.RollbackFailed,domain.Cancelled})})
 addr:=os.Getenv("UPDATEGUARD_LISTEN_ADDR");if addr==""{addr=":8080"};log.Printf("UpdateGuard control plane listening on %s",addr);log.Fatal(http.ListenAndServe(addr,securityHeaders(mux)))
}
func write(w http.ResponseWriter,v any){w.Header().Set("Content-Type","application/json");_ = json.NewEncoder(w).Encode(v)}
func securityHeaders(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("X-Content-Type-Options","nosniff");w.Header().Set("Cache-Control","no-store");next.ServeHTTP(w,r)})}
