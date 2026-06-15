package tasks_events

import "time"

type TaskCreated struct {
	TaskID string `json:"task_id"`

	Title string `json:"title"`

	Description *string `json:"description" jsonschema:"omitempty" example:"null"`

	AuthorUserID string `json:"author_user_id"`

	AuthorName string `json:"author_name"`

	Completed bool `json:"completed"`

	CreatedAt time.Time `json:"created_at"`
}

// GetKey - возвращает ключ для партиционирования в Kafka
// Все сообщения с одинаковым ключом попадают в одну партицию
// Это гарантирует порядок обработки для конкретной задачи
func (e TaskCreated) GetKey() string {
	return e.TaskID
}

// GetType - возвращает тип события
// Используется потребителем для маршрутизации (например, task.created -> Telegram)
func (e TaskCreated) GetType() string {
	return "task.created"
}
