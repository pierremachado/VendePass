package server

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"vendepass/internal/dao"
	"vendepass/internal/models"
)

func HandleConn(conn net.Conn) {
	defer conn.Close()

	session := &models.Session{Connection: conn, LastTimeActive: time.Now()}
	dao.GetSessionDAO().Insert(session)

	defer dao.GetSessionDAO().Delete(session)

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("erro na leitura do buffer", err)
			return
		}

		session.LastTimeActive = time.Now()

		if n > 0 {
			var request models.Request

			err = json.Unmarshal(buffer[:n], &request)
			if err != nil {
				fmt.Println("Erro ao desserializar o JSON:", err)
				continue
			}

			handleRequest(request, conn, session)
		}
	}

}

func CleanupSessions(timeout time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, session := range dao.GetSessionDAO().FindAll() {
			if time.Since(session.LastTimeActive) > timeout {
				fmt.Printf("Encerrando sessão %s por inatividade\n", session.ID)
				session.Connection.Close()
				dao.GetSessionDAO().Delete(session)
			}
		}
	}
}

func handleRequest(request models.Request, conn net.Conn, session *models.Session) {
	switch {
	case request.Action == "login":
		login(request.Data, conn, session)
	case request.Action == "logout":
		logout(request.Data, conn)
	}
}
