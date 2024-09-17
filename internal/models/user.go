package models

type User struct {
	BaseModel
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Tasks    []Task `json:"tasks,omitempty" gorm:"foreignKey:UserId"`
	Role     Role   `json:"role,omitempty" gorm:"default:user"`
	Password string `json:"-"`
}

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)
