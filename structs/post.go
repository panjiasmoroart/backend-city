package structs

type PostResponse struct {
	Id         uint   `json:"id"`
	Image      string `json:"image"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Content    string `json:"content"`
	CategoryID uint   `json:"category_id"`
	UserID     uint   `json:"user_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PostWithRelationResponse struct {
	Id        uint                   `json:"id"`
	Image     string                 `json:"image"`
	Title     string                 `json:"title"`
	Slug      string                 `json:"slug"`
	Content   string                 `json:"content,omitempty"`
	Category  CategorySimpleResponse `json:"category,omitempty"` // Hanya nama kategori
	User      UserSimpleResponse     `json:"user,omitempty"`     // Hanya nama user
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

type PostCreateRequest struct {
	// Image      string `form:"image" binding:"required"`
	Title      string `form:"title" binding:"required"`
	Content    string `form:"content" binding:"required"`
	CategoryID uint   `form:"category_id" binding:"required"`
}

type PostUpdateRequest struct {
	// Image      string `form:"image"`
	Title      string `form:"title" binding:"required"`
	Content    string `form:"content" binding:"required"`
	CategoryID uint   `form:"category_id" binding:"required"`
}
