package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	api "github.com/makhkets/7.17.25/internal/api/gen"
)

type Repository interface {
	CreateTask(ctx context.Context, task *api.Task) error

	UpdateTaskByID(ctx context.Context, taskID uuid.UUID, task *api.Task) error
	FindTaskByID(ctx context.Context, taskID uuid.UUID) (*api.Task, error)
}

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) CreateTask(ctx context.Context, task *api.Task) error {
	if err := s.Repo.CreateTask(ctx, task); err != nil {
		return err
	}
	return nil
}

// добавление таски
func (s *Service) AddFileToTask(ctx context.Context, taskID uuid.UUID, fileURL string, maxFiles int, allowedExtensions []string) {
	task, err := s.findTask(ctx, taskID)
	if err != nil {
		slog.Error("Failed to find task",
			slog.String("taskID", taskID.String()),
			slog.String("error", err.Error()))
		return
	}

	currentFilesCount := s.HowManyFiles(taskID)

	// проверяем не превышен ли лимит файлов
	if currentFilesCount >= maxFiles {
		slog.Warn("File limit exceeded, skipping download",
			slog.String("taskID", taskID.String()),
			slog.Int("currentFiles", currentFilesCount),
			slog.Int("maxFiles", maxFiles))
		return
	}

	_, err = s.downloadFile(ctx, fileURL, taskID, allowedExtensions)
	if err != nil {
		slog.Error("Failed to download file",
			slog.String("taskID", taskID.String()),
			slog.String("fileURL", fileURL),
			slog.String("error", err.Error()))
		fmt.Println("error", err)
		return
	}

	// проверяем количество файлов после загрузки
	newFilesCount := s.HowManyFiles(taskID)

	// если количество файлов равно максимальному, то завершаем задачу
	if newFilesCount >= maxFiles {
		task.Status = api.TaskStatusCompleted
		task.UpdatedAt = time.Now()
		err := s.Repo.UpdateTaskByID(ctx, taskID, task)
		if err != nil {
			slog.Error("Failed to update task status",
				slog.String("taskID", taskID.String()),
				slog.String("error", err.Error()))
		} else {
			slog.Info("Task status updated to completed",
				slog.String("taskID", taskID.String()))
		}
	}
}
