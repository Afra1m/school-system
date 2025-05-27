# School Management System

Система управления школой для завучей и учителей. Позволяет управлять списками учеников и учителей, вести учет оценок, анализировать успеваемость.

## Запуск проекта

### Требования
- Go 1.16+
- Node.js 14+
- PostgreSQL 12+

### Шаги запуска

1. **Настройка базы данных**
```bash
psql -U your_username -d your_database -f backend/database/migrations.sql
```

2. **Запуск бэкенда**
```bash
cd backend
go mod download
go run main.go
```

3. **Запуск фронтенда**
```bash
cd frontend
npm install
npm start
```

### Тестовые аккаунты
- Завуч: deputy3 / password12345
- Учитель: teacher1 / password12345