package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

type visit struct {
    Visitor string
    Time    time.Time
}

func randomString(n int) string {
    letters := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func checkError(err error) {
    if err == nil {
        return
    }
    fmt.Printf("error: %s\n", err.Error())
    os.Exit(1)
}

var database *sql.DB

func initDatabase() {
    var err error = nil
    database, err = sql.Open("sqlite3", "./tracker.db")
    checkError(err)
    statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS visits (tracker TEXT, visitor TEXT, time DATETIME DEFAULT CURRENT_TIMESTAMP);")
    checkError(err)
    _, err = statement.Exec()
    checkError(err)
}

func getVisitsCount(tracker string) int {
    count := 0
    rows, err := database.Query("SELECT COUNT(*) FROM visits WHERE tracker = ?;", tracker)
    checkError(err)
    rows.Next()
    err = rows.Scan(&count)
    checkError(err)
    err = rows.Close()
    checkError(err)
    return count
}

func addVisit(tracker string, visitor string) {
    statement, err := database.Prepare("INSERT INTO visits (tracker, visitor) VALUES (?, ?)")
    checkError(err)
    _, err = statement.Exec(tracker, visitor)
    checkError(err)
}

func getVisits(tracker string) []visit {
    count := getVisitsCount(tracker)
    if count == 0 {
        fmt.Print("Baaad\n")
        os.Exit(1)
    }

    visits := make([]visit, count)

    rows, err := database.Query("SELECT * FROM visits WHERE tracker = ?;", tracker)
    defer rows.Close()
    checkError(err)
    i := 0
    for rows.Next() {
        var ignore string
        rows.Scan(&ignore, &visits[i].Visitor, &visits[i].Time)
        i++
    }
    return visits
}

func createTracker() string {
    tracker := randomString(50)
    for getVisitsCount(tracker) != 0 {
        tracker = randomString(50)
    }
    addVisit(tracker, "-1.-1.-1.-1")
    return tracker
}

func httpVisits(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tracker := vars["tracker"]
    if getVisitsCount(tracker) == 0 {
        w.WriteHeader(404)
    } else {
        w.WriteHeader(200)
        visits := getVisits(tracker)
        bytes, _ := json.Marshal(visits)
        w.Write(bytes)
    }
}

func httpCreate(w http.ResponseWriter, r *http.Request) {
    tracker := createTracker()
    w.WriteHeader(200)
    bytes, _ := json.Marshal(tracker)
    w.Write(bytes)
}

func httpTrack(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tracker := vars["tracker"]
    if getVisitsCount(tracker) == 0 {
        w.WriteHeader(404)
    } else {
        ip := r.RemoteAddr
        forwarded := r.Header.Get("X-FORWARDED-FOR")
        if forwarded != "" {
            ip = forwarded
        }
        addVisit(tracker, ip)
        w.WriteHeader(200)
    }
}


func main() {
    // seed
    rand.Seed(time.Now().UnixNano())

    //database
    initDatabase()
    defer database.Close()

    //router
    router := mux.NewRouter()
    router.HandleFunc("/visits/{tracker:[0-9a-zA-Z]+}", httpVisits)
    router.HandleFunc("/track/{tracker:[0-9a-zA-Z]+}", httpTrack)
    router.HandleFunc("/create", httpCreate)
    http.Handle("/", router)
    http.ListenAndServe(":3000", nil)
}

