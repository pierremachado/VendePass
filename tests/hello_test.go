package tests

import (
	"testing"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

const ret string = "Helo"

func TestAssertHello(t *testing.T) {
	if ret != "Hello" {
		t.Fatalf("should have been test")
	}
}

// enter on this folder and type on terminal "go test"
func TestInsertTicket(t *testing.T) {
	ticketRepo := dao.GetTicketDAO()

	clientId := uuid.New()
	flightId := uuid.New()

	ticket := models.Ticket{ClientId: clientId, FlightId: flightId}

	ticketRepo.Insert(&ticket)

	if ticket.Id == uuid.Nil {
		t.Fatalf("Id should be added")
	}

	if daoticket, _ := ticketRepo.FindById(ticket.Id); daoticket.Id != ticket.Id {
		t.Fatalf("ticket and ticket added should be equal")
	}

}
