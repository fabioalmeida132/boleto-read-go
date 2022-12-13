package models

type Boleto struct {
	TypeableLine string   `json:"typeableLine"`
	BarCode      string   `json:"barCode"`
	FindTypes    []string `json:"findTypes"`
	HasPassword  bool     `json:"hasPassword"`
	Message      string   `json:"message,omitempty"`
}
