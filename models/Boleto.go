package models

type Boleto struct {
	TypeableLine string `json:"typeableLine"`
	BarCode      string `json:"barCode"`
}
