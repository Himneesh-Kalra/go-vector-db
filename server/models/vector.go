package models

type Vector struct {
	ID     string            `json:"id"`
	Values []float32         `json:"values"`
	Meta   map[string]string `json:"meta,omitempty"`
}
