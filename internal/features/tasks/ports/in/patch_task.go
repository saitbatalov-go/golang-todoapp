package ports_in

import (
	"github.com/google/uuid"
	"github.com/saitbatalov-go/golang-todoapp/internal/core/domain"
)

type PatchTaskParams struct {
	ID    uuid.UUID
	Patch domain.TaskPatch
}

func NewPatchTaskParams(
	id uuid.UUID,
	patch domain.TaskPatch,
) PatchTaskParams {
	return PatchTaskParams{
		ID:    id,
		Patch: patch,
	}
}

type PatchTaskResult struct {
	Task domain.Task
}

func NewPatchTaskResult(
	task domain.Task,
) PatchTaskResult {
	return PatchTaskResult{
		Task: task,
	}
}
