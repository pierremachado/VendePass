package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

type Auth struct {
	Token uuid.UUID `json:"token"`
}

func main() {
	// var clients []models.Client
	// readJson(&clients)

	var response *models.Response
	var tokens []uuid.UUID
	var res Auth
	var bytes []byte
	clients := dao.GetClientDAO().FindAll()
	for i := 0; i < 100; i++ {
		client := clients[i%len(clients)]
		response = writeAndReturnResponse(models.Request{
			Action: "login",
			Data: models.LoginCredentials{
				Username: client.Username,
				Password: client.Password,
			},
		})
		var error models.Response

		bytes, _ = json.Marshal(response)
		json.Unmarshal(bytes, &error)
		fmt.Println(error.Error)
		bytes, _ = json.Marshal(response.Data)

		json.Unmarshal(bytes, &res)

		tokens = append(tokens, res.Token)
	}

	id, _ := uuid.Parse("650e8400-e29b-41d4-a716-446655440011")

	var wg sync.WaitGroup
	start := make(chan struct{}) // Canal para sinalizar o início das reservas

	for _, token := range tokens {
		wg.Add(1)

		go func(t uuid.UUID) {
			defer wg.Done() // Marca essa goroutine como concluída ao terminar

			// Espera até que a sinalização para iniciar as reservas seja recebida
			<-start

			response := writeAndReturnResponse(models.Request{
				Action: "reservation",
				Auth:   t.String(),
				Data: models.FlightsRequest{
					FlightIds: []uuid.UUID{id},
				},
			})
			// Exibe o erro se houver
			if response.Error != "" {
				fmt.Println(response.Error)

			} else {
				fmt.Println("Reserva realizada com sucesso: ", token)
				fmt.Println(response.Data)
			}
		}(token)
	}

	// Inicia todas as goroutines ao mesmo tempo
	close(start) // Sinaliza que as reservas podem começar

	// Espera todas as goroutines finalizarem
	wg.Wait()
	fmt.Println("Todas as reservas foram processadas.")

}

func readJson(clients *[]models.Client) {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	jsonPath := filepath.Join(baseDir, "internal", "stubs", "clients.json")

	f, err := os.OpenFile(jsonPath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)

	err = json.Unmarshal(b, &clients)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}

}

func writeAndReturnResponse(req models.Request) *models.Response {
	conn, err := net.Dial("tcp", "localhost"+":"+"8888")
	if err != nil {
		return nil
	}
	defer conn.Close()

	buffer, _ := json.Marshal(req)

	_, writeErr := conn.Write(buffer)
	if writeErr != nil {
		return nil
	}

	receive := make([]byte, 2048)
	n, readErr := conn.Read(receive)
	if readErr != nil {
		return nil
	}

	var responseData models.Response
	err = json.Unmarshal(receive[:n], &responseData)
	if err != nil {
		return nil
	}

	return &responseData
}
