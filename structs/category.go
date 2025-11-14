package structs

type CategoryResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CategorySimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
