package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

func HandleConn(conn net.Conn, sessions map[uuid.UUID]*models.Session) {
	defer conn.Close()
	var request models.Request

	sessionID := uuid.New()
	session := &models.Session{ID: sessionID, Connection: conn, LastTimeActive: time.Now()}
	sessions[sessionID] = session

	defer delete(sessions, sessionID)

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)

		if err == io.EOF {
			fmt.Println("Fim da conexão com ", sessionID)
			return
		} else if err != nil {
			fmt.Println("erro na leitura do buffer", err)
		}

		session.LastTimeActive = time.Now()

		json.Unmarshal(buffer[:n], &request)
		handleRequest(request, conn)
	}

}

func handleRequest(request models.Request, conn net.Conn) {
	switch {
	case request.Action == "login":
		login(request.Data, conn)
	}
}
