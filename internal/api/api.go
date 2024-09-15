package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"vendepass/internal/models"
)

const (
	port      = ":8081"
	CONN_PORT = "8080"
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
)

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/user", handleGetUser)
	log.Fatal(http.ListenAndServe(port, nil))
}

func allowCrossOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{Action: "get-user", Auth: token})

}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{Action: "logout", Auth: token})

}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var logCred models.LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&logCred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "login",
		Data:   logCred,
	})
}

func writeAndReturnResponse(w http.ResponseWriter, req models.Request) {
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	buffer, _ := json.Marshal(req)
	_, writeErr := conn.Write(buffer)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}

	receive := make([]byte, 2048)
	n, readErr := conn.Read(receive)
	if readErr != nil {
		http.Error(w, readErr.Error(), http.StatusInternalServerError)
		return
	}

	var responseData models.Response
	err = json.Unmarshal(receive[:n], &responseData)
	if err != nil {
		http.Error(w, "Failed to decode response from server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}
