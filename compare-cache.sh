#!/bin/bash

echo "========================================="
echo "Сравнение производительности с Redis и без"
echo "========================================="

# Функция для тестирования
test_performance() {
    local name=$1
    echo ""
    echo "📊 $name"
    echo "-----------------------------------------"
    
    # 5 запросов для прогрева
    for i in {1..5}; do
        curl -s http://localhost:5050/api/v1/tasks > /dev/null
    done
    
    # 100 запросов и замер времени
    local total=0
    for i in {1..100}; do
        time=$(curl -s -o /dev/null -w "%{time_total}" http://localhost:5050/api/v1/tasks)
        total=$(echo "$total + $time" | bc)
    done
    
    local avg=$(echo "scale=3; $total / 100 * 1000" | bc)
    echo "Среднее время: ${avg} мкс (микросекунд)"
    
    # Быстрый тест с ab
    echo "Запуск нагрузочного теста (100 запросов, 10 конкурентных)..."
    ab -n 100 -c 10 http://localhost:5050/api/v1/tasks 2>/dev/null | grep -E "Requests per second|Time per request"
}

# Тест с Redis
test_performance "✅ С Redis (кеш включён)"

# Останавливаем Redis
echo ""
echo "🔄 Останавливаем Redis..."
docker-compose stop todoapp-redis
sleep 2

# Тест без Redis
test_performance "❌ Без Redis (прямые запросы в БД)"

# Запускаем Redis обратно
echo ""
echo "🔄 Запускаем Redis..."
docker-compose start todoapp-redis

echo ""
echo "✅ Сравнение завершено!"
