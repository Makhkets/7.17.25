package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	api "github.com/makhkets/7.17.25/internal/api/gen"
	"github.com/makhkets/7.17.25/internal/config"
	"github.com/makhkets/7.17.25/pkg/utils"
)

type Service interface {
	CreateTask(ctx context.Context, task *api.Task) error
	AddFileToTask(ctx context.Context, taskID uuid.UUID, fileURL string, maxFiles int)

	// FindTaskByID(ctx context.Context, taskID string) (*api.Task, error)
}

type ServerAPI struct {
	Service Service
	Config  *config.Config
}

func NewServerAPI(service Service, config *config.Config) *ServerAPI {
	return &ServerAPI{
		Service: service,
		Config:  config,
	}
}

// Создать новую задачу.
// POST /tasks
func (s *ServerAPI) CreateTask(ctx context.Context) (api.CreateTaskRes, error) {
	task := &api.Task{
		ID:        uuid.New(),
		Status:    api.TaskStatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.Service.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	slog.Debug("task created", slog.Any("task", task))
	return task, nil
}

// AddFileToTask implements addFileToTask operation.
// Добавить файл в задачу.
// POST /tasks/{taskId}/files
func (s *ServerAPI) AddFileToTask(ctx context.Context, req *api.AddFileRequest, params api.AddFileToTaskParams) (api.AddFileToTaskRes, error) {
	slog.Info("AddFileToTask handler called",
		slog.String("taskID", params.TaskId.String()),
		slog.String("fileURL", req.URL.String()),
		slog.Int("maxFiles", s.Config.Filter.MaxFiles))

	if !utils.IsValidURL(req.URL.String()) {
		slog.Error("Invalid URL provided", slog.String("url", req.URL.String()))
		return nil, fmt.Errorf("invalid url")
	}

	// добавить логику сервис и репозиторий
	s.Service.AddFileToTask(ctx, params.TaskId, req.URL.String(), s.Config.Filter.MaxFiles)

	return &api.Task{}, nil
}

// GetTaskStatus implements getTaskStatus operation.
// Получить статус задачи.
// GET /tasks/{taskId}
func (s *ServerAPI) GetTaskStatus(ctx context.Context, params api.GetTaskStatusParams) (api.GetTaskStatusRes, error) {
	panic("not implemented")

}

// GetTasks implements getTasks operation.
// Получить список всех задач.
// GET /tasks
func (s *ServerAPI) GetTasks(ctx context.Context) ([]api.Task, error) {
	return nil, nil
}
