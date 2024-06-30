// cmd/shared/model/jsonModel.go

package model

// Status defines the standard status structure
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SingleResponse defines the standard single response structure
type SingleResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// PagedResponse defines the standard paged response structure
type PagedResponse struct {
	Status Status        `json:"status"`
	Data   []interface{} `json:"data,omitempty"`
	Paging Paging        `json:"paging"`
}
