package structs

// DashboardResponse struct untuk response data dashboard
type DashboardResponse struct {
	CategoriesCount int64 `json:"categories_count"`
	PostsCount      int64 `json:"posts_count"`
	ProductsCount   int64 `json:"products_count"`
	AparatursCount  int64 `json:"aparaturs_count"`
}
