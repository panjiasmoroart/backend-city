package structs

type PageResponse struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PagetWithRelationResponse struct {
	Id        uint               `json:"id"`
	Title     string             `json:"title"`
	Slug      string             `json:"slug"`
	Content   string             `json:"content,omitempty"`
	User      UserSimpleResponse `json:"user,omitempty"` // Hanya nama user
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}

type PageCreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type PageUpdateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
