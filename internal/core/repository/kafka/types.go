package core_kafka

import "context"


type Event interface {
    // GetKey - возвращает ключ для партиционирования
    // Сообщения с одинаковым ключом попадают в одну партицию
    GetKey() string
    
    // GetType - возвращает тип события (например, "task.created")
    // Используется для маршрутизации у потребителя
    GetType() string
}


type Producer interface {
    // Publish - отправляет событие в брокер
    // key - ключ для партиционирования
    // event - само событие (должно реализовать интерфейс Event)
    Publish(ctx context.Context, key string, event Event) error
    
    Close() error
}