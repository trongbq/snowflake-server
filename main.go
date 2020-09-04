package main

import (
  "net/http"
  "log"
)


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
    w.Write([]byte("Not Implemented"))
}

func stats(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Not Implemented"))
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    })
}
