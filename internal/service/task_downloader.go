package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	api "github.com/makhkets/7.17.25/internal/api/gen"
	"github.com/makhkets/7.17.25/pkg/utils"
)

// type Task struct {
// 	ID uuid.UUID `json:"id"`
// 	Status TaskStatus `json:"status"`
// 	Files []FileInfo `json:"files"`
// 	ArchiveUrl OptNilString `json:"archiveUrl"`
// 	CreatedAt time.Time `json:"createdAt"`
// 	UpdatedAt time.Time `json:"updatedAt"`

// URL url.URL `json:"url"`
// Filename string `json:"filename"`
// Status FileInfoStatus `json:"status"`
// Error OptNilString `json:"error"`

func (s *Service) downloadFile(ctx context.Context, fileURL string, uuid uuid.UUID, NotAllowedExtensions []string) (*api.FileInfo, error) {
	// извлекаем оригинальное имя файла из URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	path := parsedURL.Path
	if path == "" {
		path = "/"
	}

	// извлекаем имя из пути
	originalFilename := ""
	if path != "/" {
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			originalFilename = parts[len(parts)-1]
		}
	}

	if originalFilename == "" || !strings.Contains(originalFilename, ".") {
		originalFilename = uuid.String()
	}

	extension := ""
	if strings.Contains(originalFilename, ".") {
		parts := strings.Split(originalFilename, ".")
		if len(parts) > 1 {
			extension = "." + parts[len(parts)-1]
		}
	}

	// Проверяем, разрешено ли расширение файла
	isAllowed := true
	for _, ext := range NotAllowedExtensions {
		if strings.EqualFold(extension, ext) {
			isAllowed = false
			break
		}
	}
	if !isAllowed {
		slog.Warn("File extension not allowed, skipping download",
			slog.String("taskID", uuid.String()),
			slog.String("fileURL", fileURL),
			slog.String("filename", originalFilename),
			slog.String("extension", extension),
		)

		task, err := s.Repo.FindTaskByID(ctx, uuid)
		if err != nil {
			slog.Error("Failed to find task for file", slog.String("taskID", uuid.String()), slog.String("error", err.Error()))
			return nil, err
		}

		fileInfo := api.FileInfo{
			URL:      *parsedURL,
			Filename: originalFilename,
			Status:   api.FileInfoStatusFailed,
			Error:    api.OptNilString{Value: "file extension not allowed", Null: false},
		}

		task.Files = append(task.Files, fileInfo)
		if err := s.Repo.UpdateTaskByID(ctx, uuid, task); err != nil {
			slog.Error("Failed to update task with file", slog.String("taskID", uuid.String()), slog.String("error", err.Error()))
			return nil, err
		}
		return &fileInfo, nil
	}

	uniqueFilename := uuid.String() + "_" + originalFilename
	filePath := utils.FindDirectoryName("photo_storage") + "/" + uniqueFilename

	task, err := s.Repo.FindTaskByID(ctx, uuid)
	if err != nil {
		slog.Error("Failed to find task for file addition",
			slog.String("taskID", uuid.String()),
			slog.String("error", err.Error()))
		return nil, err
	}

	status := api.FileInfoStatusPending
	fileInfo := api.FileInfo{
		URL:      *parsedURL,
		Filename: uniqueFilename,
		Status:   api.FileInfoStatusPending,
		Error:    api.OptNilString{},
	}

	task.Files = append(task.Files, fileInfo)
	err = s.Repo.UpdateTaskByID(ctx, uuid, task)
	if err != nil {
		slog.Error("Failed to update task with new file",
			slog.String("taskID", uuid.String()),
			slog.String("error", err.Error()))
		return nil, err
	}

	// окончательно обновляем статус файла
	defer func() {

		fileinfo, err := s.FindAndDeleteFile(task, fileInfo.Filename)
		if err != nil {
			return
		}
		fileinfo.Status = status
		s.Repo.UpdateTaskByID(ctx, uuid, task)

	}()

	resp, err := http.Get(fileURL)
	if err != nil {
		status = api.FileInfoStatusFailed
		fileInfo.Error = api.OptNilString{
			Value: err.Error(),
			Null:  false,
		}
		return nil, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		status = api.FileInfoStatusFailed
		fileInfo.Error = api.OptNilString{
			Value: err.Error(),
			Null:  false,
		}

		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	bytesWritten, err := io.Copy(out, resp.Body)
	if err != nil {
		status = api.FileInfoStatusFailed
		fileInfo.Error = api.OptNilString{
			Value: err.Error(),
			Null:  false,
		}
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	status = api.FileInfoStatusDownloaded
	slog.Info("File downloaded successfully",
		slog.String("taskID", uuid.String()),
		slog.String("filePath", filePath),
		slog.Int64("bytesWritten", bytesWritten))

	return &fileInfo, nil
}

// количество файлов в задаче
func (s *Service) HowManyFiles(uuid uuid.UUID) int {
	dir := utils.FindDirectoryName("photo_storage")
	files, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("Failed to read directory",
			slog.String("taskID", uuid.String()),
			slog.String("directory", dir),
			slog.String("error", err.Error()))
		return 0
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if len(name) >= len(uuid.String()) && name[:len(uuid.String())] == uuid.String() {
				count++
			}
		}
	}

	return count
}

// находим файл в задаче по названию
func (s *Service) FindAndDeleteFile(task *api.Task, filename string) (*api.FileInfo, error) {
	for i, file := range task.Files {
		if file.Filename == filename {
			foundFile := file
			task.Files = append(task.Files[:i], task.Files[i+1:]...)

			return &foundFile, nil
		}
	}
	return nil, fmt.Errorf("file %s not found in task", filename)
}
