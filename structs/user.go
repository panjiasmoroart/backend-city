package structs

// Struct ini digunakan untuk menampilkan data user sebagai response API
type UserResponse struct {
	Id          uint            `json:"id"`
	Name        string          `json:"name"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
	Permissions map[string]bool `json:"permissions,omitempty"` // Menampilkan permissions yang dimiliki user
	Roles       []RoleResponse  `json:"roles,omitempty"`
	Token       *string         `json:"token,omitempty"`
	CreatedAt   string          `json:"created_at,omitempty"`
	UpdatedAt   string          `json:"updated_at,omitempty"`
}

type UserSimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Struct ini digunakan untuk menerima data saat proses create user
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password" binding:"required"`
	RoleIDs  []uint `json:"role_ids" binding:"required"` // IDs of roles assigned to the user
}

// Struct ini digunakan untuk menerima data saat proses update user
type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password,omitempty"`
	RoleIDs  []uint `json:"role_ids"` // IDs of roles assigned to the user
}

// Struct ini digunakan saat user melakukan proses login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
