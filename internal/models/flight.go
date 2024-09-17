package models

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type FlightJSON struct {
	Id              uuid.UUID
	SourceAirportId uuid.UUID
	DestAirportId   uuid.UUID
	Passengers      []Ticket
	Seats           uint
}

type Flight struct {
	Id              uuid.UUID
	SourceAirportId uuid.UUID
	DestAirportId   uuid.UUID
	Passengers      []Ticket
	Seats           uint
	Queue           chan *Session // Canal de fila para reservas
	mu              sync.Mutex
}

func (f *Flight) AcceptReservation(reservation Reservation) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.Seats > 0 {
		f.Seats--
		f.Passengers = append(f.Passengers, reservation.Ticket)
		return nil
	}
	return errors.New("no seats available")

}
