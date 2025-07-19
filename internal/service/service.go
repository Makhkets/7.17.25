package service

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	api "github.com/makhkets/7.17.25/internal/api/gen"
	"github.com/makhkets/7.17.25/pkg/utils"
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
func (s *Service) AddFileToTask(ctx context.Context, taskID uuid.UUID, fileURL string, maxFiles int, allowedExtensions []string, address string) {
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

	file, err := s.downloadFile(ctx, fileURL, taskID, allowedExtensions)
	if err != nil {
		slog.Error("Failed to download file",
			slog.String("taskID", taskID.String()),
			slog.String("fileURL", fileURL),
			slog.String("error", err.Error()))
		fmt.Println("error", err)
		return
	}

	file.Status = api.FileInfoStatusDownloaded
	task.Files = append(task.Files, *file)

	// проверяем количество файлов после загрузки
	newFilesCount := s.HowManyFiles(taskID)

	// если количество файлов равно максимальному, то завершаем задачу
	if newFilesCount >= maxFiles {
		defer func() {
			err := s.Repo.UpdateTaskByID(ctx, taskID, task)
			if err != nil {
				slog.Error("Failed to update task status",
					slog.String("taskID", taskID.String()),
					slog.String("error", err.Error()))
			}
		}()

		slog.Info("Task status updated to completed",
			slog.String("taskID", taskID.String()))

		// uniqueFilename := uuid.String() + "_" + originalFilename
		// filePath := utils.FindDirectoryName("photo_storage") + "/" + uniqueFilename

		files, err := filepath.Glob(utils.FindDirectoryName("photo_storage") + "/" + taskID.String() + "_*")
		if err != nil {
			slog.Error("Failed to get files",
				slog.String("error", err.Error()))
		}
		slog.Info("Files", slog.Any("files", files))

		archiveName := taskID.String() + ".zip"
		err = s.createArchive(files, archiveName)
		if err != nil {
			slog.Error("Failed to create archive",
				slog.String("error", err.Error()))
		}
		slog.Debug("Archive created", slog.String("archiveName", archiveName))

		// localhost:8080/download/taskID

		task.ArchiveUrl = api.OptNilString{
			// /tasks/{taskId}/download:
			Value: fmt.Sprintf("http://%s/download/%s", address, taskID),
			Null:  false,
		}
		task.Status = api.TaskStatusCompleted
		task.UpdatedAt = time.Now()

		// пользователь запрашивает статус, получает ссылку на архив, после чего удаляется архив и картинки
		// если пользователь оставил картинки и не завершил задачу, то сервис сам все почистит
	}
}

func (s *Service) GetTaskByID(ctx context.Context, taskID uuid.UUID, address string, maxFiles int) (*api.Task, error) {
	task, err := s.findTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// проверяем количество файлов после загрузки
	newFilesCount := s.HowManyFiles(taskID)
	if newFilesCount >= maxFiles {
		task.ArchiveUrl.SetTo(fmt.Sprintf("http://%s/download/%s", address, taskID))
	}

	return task, nil
}
