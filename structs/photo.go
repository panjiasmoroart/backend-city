package structs

type PhotoResponse struct {
	Id          uint   `json:"id"`
	Image       string `json:"image"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PhotoCreateRequest struct {
	Caption     string `form:"caption" binding:"required"`
	Description string `form:"description" binding:"required"`
}
