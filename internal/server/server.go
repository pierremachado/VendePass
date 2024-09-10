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

	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("erro na leitura do buffer", err)
		return
	}

	if n > 0 {
		var request models.Request

		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			fmt.Println("Erro ao desserializar o JSON:", err)
		}

		handleRequest(request, conn)
	}

}

func CleanupSessions(timeout time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, session := range dao.GetSessionDAO().FindAll() {
			if time.Since(session.LastTimeActive) > timeout {
				fmt.Printf("Encerrando sessão %s por inatividade\n", session.ID)
				dao.GetSessionDAO().Delete(session)
			}
		}
	}
}

func handleRequest(request models.Request, conn net.Conn) {
	switch {
	case request.Action == "login":
		login(request.Data, conn)
	case request.Action == "logout":
		logout(request.Data, conn)
	}
}

func WriteNewResponse(response models.Response, conn net.Conn) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}
