package entity

type Product struct {
	ID		string `json:"_id,omitempty"`   // Optional: will be omitted from JSON if empty
	Rev		string `json:"_rev,omitempty"`  // Optional: will be omitted from JSON if empty
	Name 	string `json:"name" validate:"required,min=3,max=100"` // Required: must be between 3 and 100 characters
	Price 	float64 `json:"price" validate:"required,gt=0"`       // Required: must be greater than 0
}

