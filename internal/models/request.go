package models

type Request struct {
	Action string `json:"Action"`
	Data   interface{}
}
