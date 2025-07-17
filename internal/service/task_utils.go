package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	api "github.com/makhkets/7.17.25/internal/api/gen"
	"github.com/makhkets/7.17.25/internal/repository"
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
