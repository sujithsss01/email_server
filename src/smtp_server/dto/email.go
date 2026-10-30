package dto

type Email struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
