package ports_in

import "github.com/google/uuid"


type DeleteTaskParams struct {
	ID uuid.UUID
}

func NewDeleteTaskParams(
	id uuid.UUID,
) DeleteTaskParams {
	return DeleteTaskParams{
		ID: id,
	}
}

type DeleteTaskResult struct{}

func NewDeleteTaskResult() DeleteTaskResult {
	return DeleteTaskResult{}
}
