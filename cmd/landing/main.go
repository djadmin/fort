package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = "7431"
	csvPath     = "waitlist.csv"
	staticDir   = "landing"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	mux.HandleFunc("/api/waitlist", handleWaitlist)

	log.Printf("fort landing → http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleWaitlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"email required"}`)
		return
	}

	if err := appendEmail(body.Email); err != nil {
		log.Printf("waitlist write: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"server error"}`)
		return
	}

	log.Printf("waitlist: %s", body.Email)
	fmt.Fprintln(w, `{"ok":true}`)
}

func appendEmail(email string) error {
	f, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	cw := csv.NewWriter(f)
	if err := cw.Write([]string{email, time.Now().UTC().Format(time.RFC3339)}); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}
