package services

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go-api/internal/dal"
	"go-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User interface {
	GetID() string
	GetName() string
	GetEmail() string
}

func NewUserResponse(model *models.User, includeTasks bool) User {
	return CreateUserDTO(model, includeTasks)
}

type UserBase struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type UserDTO struct {
	UserBase
	Tasks []interface{} `json:"tasks,omitempty"`
}

func CreateUserDTO(model *models.User, expand bool) User {
	tasks := make([]interface{}, 0, len(model.Tasks))
	for _, task := range model.Tasks {
		if expand {
			tasks = append(tasks, Task{
				ID:     task.ID,
				Title:  task.Title,
				UserId: task.UserId,
			})
		} else {
			tasks = append(tasks, task.ID)
		}

	}
	return &UserDTO{
		UserBase: UserBase{
			ID:    model.ID,
			Name:  model.Name,
			Email: model.Email,
		},
		Tasks: tasks,
	}
}

func (u UserBase) GetID() string    { return u.ID }
func (u UserBase) GetName() string  { return u.Name }
func (u UserBase) GetEmail() string { return u.Email }

type UserService interface {
	CreateUser(user models.User) (User, error)
	GetUserById(id string, expand bool) (User, error)
	RegisterUser(name string, email string, password string) (*models.User, error)
	Login(email string, hashedPassword string) (string, error)
}

type service struct {
	dao *dal.Query
}

func NewUserService(dao *dal.Query) UserService {
	return &service{dao: dao}
}

func (s *service) CreateUser(user models.User) (User, error) {
	if err := s.dao.User.Create(&user); err != nil {
		return nil, err
	}
	resp := NewUserResponse(&user, false)
	return resp, nil
}

func (s *service) GetUserById(id string, expand bool) (User, error) {
	u := s.dao.User
	user, err := u.Where(u.ID.Eq(id)).Preload(u.Tasks).First()
	if err != nil {
		return nil, err
	}

	userResp := NewUserResponse(user, expand)
	return userResp, nil
}

func (s *service) Login(email string, password string) (string, error) {
	u := s.dao.User
	user, err := u.Where(u.Email.Eq(email)).First()
	if err != nil {
		return "nil", err
	}
	if err := verifyPassword(user.Password, password); err != nil {
		return "nil", err
	}
	return s.generateJWT(user)
}

func (s *service) generateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID, // Subject (user ID)
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
		"iat":   time.Now().Unix(),                     // Issued At
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Expires in 24 hours
		"jti":   uuid.New().String(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("CRYPTO_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *service) RegisterUser(name string, email string, password string) (*models.User, error) {
	// Hash the user's password
	hashedPassword, err := hashPassword(password)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	// Create the user in the database
	if err := s.dao.User.DO.Create(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
