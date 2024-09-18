package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type FlightJSON struct {
	Id              uuid.UUID
	SourceAirportId uuid.UUID
	DestAirportId   uuid.UUID
	Passengers      []*Ticket
	Seats           uint
}

type Flight struct {
	Id              uuid.UUID
	SourceAirportId uuid.UUID
	DestAirportId   uuid.UUID
	Passengers      []*Ticket
	Seats           uint
	Queue           chan *Session // Canal de fila para reservas
	Mu              sync.Mutex
}

func (f *Flight) AcceptReservation() (*Ticket, error) {
	f.Mu.Lock()
	defer f.Mu.Unlock()

	if f.Seats > 0 {
		f.Seats--
		ticket := new(Ticket)
		ticket.Id = uuid.New()
		ticket.FlightId = f.Id
		f.Passengers = append(f.Passengers, ticket)
		return ticket, nil
	}
	return nil, errors.New("no seats available")

}

func (f *Flight) ProcessReservations() {
	for session := range f.Queue {
		ticket, err := f.AcceptReservation()
		if err != nil {
			fmt.Printf("Sessão %s: erro ao reservar para o voo %s - %s\n", session.ID, f.Id, err)
		} else {
			ticket.ClientId = session.ClientID
			id := uuid.New()
			session.Reservations[id] = Reservation{
				Id:        id,
				CreatedAt: time.Now(),
				Ticket:    ticket,
			}
			fmt.Printf("Sessão %s: voo %s reservado com sucesso!\n", session.ID, f.Id)
		}
	}
}
