package service

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
	api "github.com/makhkets/7.17.25/internal/api/gen"
	"github.com/makhkets/7.17.25/internal/repository"
	"github.com/makhkets/7.17.25/pkg/utils"
)

func (s *Service) findTask(ctx context.Context, taskID uuid.UUID) (*api.Task, error) {
	task, err := s.Repo.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			slog.Error("Task not found", slog.String("taskID", taskID.String()))
			return nil, fmt.Errorf("task not found")
		}
		slog.Error("Repository error",
			slog.String("taskID", taskID.String()),
			slog.String("error", err.Error()))
		return nil, err
	}

	return task, nil
}

func (s *Service) createArchive(files []string, archiveName string) error {
	file, err := os.Create(utils.FindDirectoryName("photo_storage") + "/" + archiveName)
	if err != nil {
		return err
	}
	defer file.Close()

	zip := zip.NewWriter(file)
	defer zip.Close()

	for _, file := range files {
		err := s.addFileToZip(zip, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) addFileToZip(zipWriter *zip.Writer, filename string) error {
	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Получаем инфо о файле
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Создаём заголовок
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Имя файла внутри архива
	header.Name = info.Name()
	header.Method = zip.Deflate // метод сжатия

	// Создаём writer для этого файла в архиве
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Копируем содержимое файла в архив
	_, err = io.Copy(writer, file)
	return err
}
