package domain

import (
	"fmt"
	"time"

	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

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

func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return Task{
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	}
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

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("invalid `CompletedAt` can't be null: %w", core_errors.ErrInvalidArgument)
		}

		if t.CompletedAt.Before(*t.CompletedAt) {
			return fmt.Errorf("invalid `CompletedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("invalid `CompletedAt` must be `nil  if `Completed` == `false`: %w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}
