// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AddFileToTask implements addFileToTask operation.
//
// Добавить файл в задачу.
//
// POST /tasks/{taskId}/files
func (UnimplementedHandler) AddFileToTask(ctx context.Context, req *AddFileRequest, params AddFileToTaskParams) (r AddFileToTaskRes, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateTask implements createTask operation.
//
// Создать новую задачу.
//
// POST /tasks
func (UnimplementedHandler) CreateTask(ctx context.Context) (r CreateTaskRes, _ error) {
	return r, ht.ErrNotImplemented
}

// DownloadTaskArchive implements downloadTaskArchive operation.
//
// Скачать ZIP архив задачи.
//
// GET /tasks/{taskId}/download
func (UnimplementedHandler) DownloadTaskArchive(ctx context.Context, params DownloadTaskArchiveParams) (r *DownloadTaskArchiveOKHeaders, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTaskStatus implements getTaskStatus operation.
//
// Получить статус задачи.
//
// GET /tasks/{taskId}
func (UnimplementedHandler) GetTaskStatus(ctx context.Context, params GetTaskStatusParams) (r GetTaskStatusRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetTasks implements getTasks operation.
//
// Получить список всех задач.
//
// GET /tasks
func (UnimplementedHandler) GetTasks(ctx context.Context) (r []Task, _ error) {
	return r, ht.ErrNotImplemented
}
