package dto

type CompanyValueResponse struct {
	ID          uint     `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Examples    []string `json:"examples,omitempty"`
}

type CompanyValueListResponse struct {
	Items []CompanyValueResponse `json:"items"`
	Total int                    `json:"total"`
}
