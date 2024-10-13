### Инициализация migrate
``migrate create -ext sql -dir ./schema -seq init``

### Применение миграций
``migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436?sslmode=disable' up``