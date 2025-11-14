package structs

type SliderResponse struct {
	Id          uint   `json:"id"`
	Image       string `json:"image"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SliderCreateRequest struct {
	Description string `form:"description" binding:"required"`
}
