package services

import (
	"go-api/internal/dal"
	"go-api/internal/models"
	"time"
)

type Task struct {
	ID        string    `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	UserId    string    `json:"userId,omitempty" gorm:"user_id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// TODO implement relations
	//User   User   `json:"user" gorm:"foreignKey:UserID"`
}

type TaskService interface {
	CreateTask(task *models.Task) (*models.Task, error)
	GetTaskById(id string) (*Task, error)
	GetTasks() ([]*models.Task, error)
}

func NewTaskResponse(model *models.Task) *Task {
	return &Task{
		ID:        model.ID,
		Title:     model.Title,
		UserId:    model.UserId,
		CreatedAt: model.CreatedAt.Local(),
	}
}

type taskService struct {
	dao *dal.Query
}

func NewTaskService(dao *dal.Query) TaskService {
	return &taskService{dao: dao}
}

func (s *taskService) CreateTask(task *models.Task) (*models.Task, error) {
	if err := s.dao.Task.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetTaskById(id string) (*Task, error) {
	t := s.dao.Task
	task, err := t.Where(t.ID.Eq(id)).Preload(t.User).First()
	if err != nil {
		return nil, err
	}
	taskResp := NewTaskResponse(task)
	return taskResp, nil
}

func (s *taskService) GetTasks() ([]*models.Task, error) {
	t := s.dao.Task
	tasks, err := t.Preload(t.User).Find()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
