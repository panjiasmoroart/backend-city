package structs

type ProductResponse struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Owner     string `json:"owner"`
	Price     int    `json:"price"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ProductWithRelationResponse struct {
	Id        uint               `json:"id"`
	Title     string             `json:"title"`
	Slug      string             `json:"slug"`
	Content   string             `json:"content,omitempty"`
	Image     string             `json:"image,omitempty"`
	Owner     string             `json:"owner"`
	Price     int                `json:"price,omitempty"`
	Address   string             `json:"address,omitempty"`
	Phone     string             `json:"phone,omitempty"`
	User      UserSimpleResponse `json:"user,omitempty"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}

type ProductCreateRequest struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	Owner   string `form:"owner" binding:"required"`
	Price   int    `form:"price" binding:"required"`
	Address string `form:"address" binding:"required"`
	Phone   string `form:"phone" binding:"required"`
}

type ProductUpdateRequest struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	Owner   string `form:"owner" binding:"required"`
	Price   int    `form:"price" binding:"required"`
	Address string `form:"address" binding:"required"`
	Phone   string `form:"phone" binding:"required"`
}
