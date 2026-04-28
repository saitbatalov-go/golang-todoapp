package domain

import (
	"fmt"
	"time"

	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
)

type Task struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	completedAt *time.Time,
	createdAt time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
	}
}

func (t *TaskPatch) Validate() error {

	if t.Title.Set && t.Title.Value == nil {
		return fmt.Errorf(
			"PATCH invalid `Title` can't be value null: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		nil,
		time.Now(),
		authorUserID,
	)
}

func (t *Task) Validate() error {

	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf("invalid `Title` len: %d:%w", titleLen, core_errors.ErrInvalidArgument)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf("invalid `Description` len: %d:%w", descriptionLen, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	if patch.Title.Set {
		t.Title = *patch.Title.Value
	}
	if patch.Description.Set {
		t.Description = patch.Description.Value // напрямую копируем указатель
	}

	if err := t.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	return nil
}
