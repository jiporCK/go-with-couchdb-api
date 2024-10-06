package entity

type Product struct {
	ID		string `json:"_id,omitempty"`
	Rev		string `json:"_rev,omitempty"`
	Name	string `json:"name"`
	Price 	float64 `json:"price"`
}
