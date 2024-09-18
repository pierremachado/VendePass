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
	http.HandleFunc("/route", handleGetRoute)
	http.HandleFunc("/flights", handleGetFlights)
	http.HandleFunc("/reservation", handleReservation)
	http.HandleFunc("/cart", handleGetCart)
	http.HandleFunc("/ticket", handleTicket)
	http.HandleFunc("/tickets", handleGetTickets)
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

func handleGetTickets(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{
		Action: "tickets",
		Auth:   token,
	})
}

func handleTicket(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	switch r.Method {
	case http.MethodPost:
		handleBuyTicket(w, r)
	case http.MethodDelete:
		handleCancelTicket(w, r)
	default:
		http.Error(w, "only POST or DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleBuyTicket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	var buyTicket models.BuyTicket

	err := json.NewDecoder(r.Body).Decode(&buyTicket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "buy",
		Auth:   token,
		Data:   buyTicket,
	})

}

func handleCancelTicket(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	var ticketId models.CancelBuyRequest

	err := json.NewDecoder(r.Body).Decode(&ticketId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "cancel-buy",
		Auth:   token,
		Data:   ticketId,
	})

}

func handleGetCart(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{
		Action: "cart",
		Auth:   token,
	})

}

func handleReservation(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	switch r.Method {
	case http.MethodPost:
		handleMakeReservations(w, r)
	case http.MethodDelete:
		handleCancelReservation(w, r)
	default:
		http.Error(w, "only POST or DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleMakeReservations(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	var flightIds models.FlightsRequest

	err := json.NewDecoder(r.Body).Decode(&flightIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	writeAndReturnResponse(w, models.Request{
		Action: "reservation",
		Auth:   token,
		Data:   flightIds,
	})

}

func handleCancelReservation(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	var reservationId models.CancelReservationRequest

	err := json.NewDecoder(r.Body).Decode(&reservationId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "cancel-reservation",
		Auth:   token,
		Data:   reservationId,
	})

}

func handleGetFlights(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	var flightIds models.FlightsRequest

	err := json.NewDecoder(r.Body).Decode(&flightIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "flights",
		Auth:   token,
		Data:   flightIds,
	})

}

func handleGetRoute(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()

	src := queryParams.Get("src")
	dest := queryParams.Get("dest")

	token := r.Header.Get("Authorization")
	writeAndReturnResponse(w, models.Request{
		Action: "route",
		Auth:   token,
		Data: models.RouteRequest{
			Source: src,
			Dest:   dest,
		},
	})

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
