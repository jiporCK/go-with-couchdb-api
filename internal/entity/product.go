package entity

// Struct a user-defined type to store a collection of different fields into a single field. 
type Product struct {
	ID		string `json:"_id,omitempty"`
	Rev		string `json:"_rev,omitempty"`
	Name 	string `json:"name" validate:"required,min=3,max=100"`
	Price 	float64 `json:"price" validate:"required,gt=0"`   
}
