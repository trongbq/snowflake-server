package snowflake


import (
    "os"
    "net/http"
    "log"
    "strings"
    "encoding/json"
)


const AuthorizationHeaderPrefix = "Bearer"

var APIKey string

func init() {
    APIKey = os.Getenv("API_KEY")
    if len(APIKey) == 0 {
        panic("API Key is missing!")
    }
}

type Server struct {
    idWorker *IDWorker
}

func NewServer(mID int) (*Server, error) {
    idWorker, err := NewIDWorker(mID)
    if err != nil {
        return nil, err
    }

    s := Server { idWorker }

    return &s, nil
}

func (s *Server) Start() {
    mux := http.NewServeMux()
    mux.Handle("/ping", http.HandlerFunc(ping))
    mux.Handle("/nextid", authMiddleware(http.HandlerFunc(s.nextID)))
    mux.Handle("/stats", authMiddleware(http.HandlerFunc(s.stats)))

    log.Println("Listening on :8000...")
    err := http.ListenAndServe(":8000", mux)
    if err != nil {
        log.Fatal(err)
    }
}

func (s *Server) nextID(w http.ResponseWriter, r *http.Request) {
    id := s.idWorker.NextID()

    idResp := struct {
        ID int64 `json:"id"`
    } { id }

    data, err := json.Marshal(idResp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
    data, err := s.idWorker.Stats()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Pong"))
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := strings.TrimSpace(r.Header.Get("Authorization"))
        if len(auth) == 0 {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }
        if !strings.HasPrefix(auth, AuthorizationHeaderPrefix) {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }
        tokens := strings.Split(auth, " ")
        if len(tokens) != 2 || tokens[1] != APIKey {
            http.Error(w, "API Key is incorrect", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
