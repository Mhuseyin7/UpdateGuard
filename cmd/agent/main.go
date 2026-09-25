// The agent is the only component mounted with the Docker socket. In production it must
// authenticate mTLS/enrollment credentials and implement only the routes in agent.Client.
package main
import("log";"net/http";"os")
func main(){
 mux:=http.NewServeMux()
 mux.HandleFunc("/health",func(w http.ResponseWriter,_ *http.Request){_,_=w.Write([]byte("ok\n"))})
 // Intentionally no proxy/general Docker API endpoint. Add handlers only after authorization checks.
 addr:=os.Getenv("UPDATEGUARD_AGENT_LISTEN_ADDR");if addr==""{addr=":8081"};log.Printf("UpdateGuard restricted agent listening on %s",addr);log.Fatal(http.ListenAndServe(addr,mux))
}
