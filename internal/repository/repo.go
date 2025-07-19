package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	api "github.com/makhkets/7.17.25/internal/api/gen"
)

var tasks = map[uuid.UUID]*api.Task{}

var ErrNotFound = errors.New("task not found")

type Repo struct {
	// db.Sql...
}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) CreateTask(ctx context.Context, task *api.Task) error {
	slog.Info("Repo.CreateTask called",
		slog.String("taskID", task.ID.String()),
		slog.String("status", string(task.Status)))

	if task == nil {
		slog.Error("Task is nil")
		return errors.New("task is nil")
	}

	tasks[task.ID] = task
	slog.Info("Task created successfully",
		slog.String("taskID", task.ID.String()),
		slog.Int("totalTasks", len(tasks)))
	return nil
}

func (r *Repo) FindTaskByID(ctx context.Context, taskID uuid.UUID) (*api.Task, error) {
	slog.Info("Repo.FindTaskByID called", slog.String("taskID", taskID.String()))
	task, ok := tasks[taskID]
	if !ok {
		slog.Error("Task not found", slog.String("taskID", taskID.String()))
		return nil, ErrNotFound
	}

	slog.Info("Task found successfully",
		slog.String("taskID", taskID.String()),
		slog.String("status", string(task.Status)),
		slog.Int("filesCount", len(task.Files)))
	return task, nil
}

func (r *Repo) UpdateTaskByID(ctx context.Context, taskID uuid.UUID, task *api.Task) error {
	slog.Info("Repo.UpdateTaskByID called",
		slog.String("taskID", taskID.String()),
		slog.String("status", string(task.Status)),
		slog.Int("filesCount", len(task.Files)))

	if _, ok := tasks[taskID]; !ok {
		slog.Error("Task not found for update", slog.String("taskID", taskID.String()))
		return ErrNotFound
	}
	tasks[taskID] = task
	slog.Info("Task updated successfully",
		slog.String("taskID", taskID.String()),
		slog.String("status", string(task.Status)))
	return nil
}
