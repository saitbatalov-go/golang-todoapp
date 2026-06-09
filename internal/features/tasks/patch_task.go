package tasks_service

import (
	"context"
	"fmt"

	ports_in "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/in"
	ports_out_repository "github.com/saitbatalov-go/golang-todoapp/internal/features/tasks/ports/out/repository"
)


func (s *TasksService) PatchTask(
	ctx context.Context,
	params ports_in.PatchTaskParams,
) (ports_in.PatchTaskResult, error) {
	repoGetTaskResult, err := s.tasksRepository.GetTask(
		ctx,
		ports_out_repository.NewGetTaskParams(params.ID),
	)
	if err != nil {
		return ports_in.PatchTaskResult{}, fmt.Errorf("get task from repository: %w", err)
	}

	task := repoGetTaskResult.Task
	if err := task.ApplyPatch(params.Patch); err != nil {
		return ports_in.PatchTaskResult{}, fmt.Errorf("apply task patch: %w", err)
	}

	updateTaskResult, err := s.tasksRepository.UpdateTask(
		ctx,
		ports_out_repository.NewUpdateTaskParams(task),
	)
	if err != nil {
		return ports_in.PatchTaskResult{}, fmt.Errorf("update task in repository: %w", err)
	}

	return ports_in.NewPatchTaskResult(
		updateTaskResult.Task,
	), nil
}
