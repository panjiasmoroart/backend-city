package structs

type AparaturResponse struct {
	Id          uint   `json:"id"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type AparaturCreateRequest struct {
	Name        string `form:"name" binding:"required"`
	Position    string `form:"position" binding:"required"`
	Description string `form:"description" binding:"required"`
}

type AparaturUpdateRequest struct {
	Name        string `form:"name"`
	Position    string `form:"position"`
	Description string `form:"description"`
}
