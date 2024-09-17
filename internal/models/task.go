package models

type Task struct {
	BaseModel
	Title  string `json:"title"`
	UserId string `json:"user_id" gorm:"user_id"`
	User   User   `json:"user" gorm:"foreignKey:UserId;references:id"`
}
