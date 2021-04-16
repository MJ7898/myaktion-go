package model

type Account struct {
	Iban string `json:"iban"`
	Name string `json:"name"`
	NameOfBank string `json:"nameOfBank"`
}
