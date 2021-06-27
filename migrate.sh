export MYSQL_CONNECTION_STRING="hasty:secret@tcp(127.0.0.1:3306)/hasty-db"
go run ./cmd/schemamigrate --direction=up --directory="./internal/storage/migrations"

# migrate create -ext sql -dir internal/storage/migrations -format "20060102150405" service
