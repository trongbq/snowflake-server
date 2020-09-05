package main

import (
  "net/http"
  "log"
  "strings"
  "encoding/json"
  "snowflake-server/snowflake"
)


const AuthorizationHeaderPrefix = "Bearer"

var idWorker *snowflake.IDWorker

func init() {
    var err error
    idWorker, err = snowflake.NewIDWorker(1)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/ping", http.HandlerFunc(ping))
    mux.Handle("/nextid", authMiddleware(http.HandlerFunc(nextID)))
    mux.Handle("/stats", authMiddleware(http.HandlerFunc(stats)))

    log.Println("Listening on :8000...")
    err := http.ListenAndServe(":8000", mux)
    if err != nil {
        log.Fatal(err)
    }
}

func ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Pong"))
}

func nextID(w http.ResponseWriter, r *http.Request) {
    id := idWorker.NextID()

    idResp := struct {
        ID int64 `json:"id"`
    } { id }

    data, err := json.Marshal(idResp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func stats(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Not Implemented"))
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, AuthorizationHeaderPrefix) {
        }
        next.ServeHTTP(w, r)
    })
}
