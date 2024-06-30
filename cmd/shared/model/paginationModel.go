// cmd/shared/model/paginationModel.go

package model

// Paging defines the pagination structure
type Paging struct {
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPerPage"`
	TotalRows   int `json:"totalRows"`
	TotalPages  int `json:"totalPages"`
}
