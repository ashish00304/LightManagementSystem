package models

type Light struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "ON" or "OFF"
}

type TurnOnRequest struct {
	Duration int `json:"duration,omitempty"`
}
