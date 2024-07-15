postgres:
	docker run --name WebSocket -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:latest

createdb:
	docker exec -it WebSocket createdb --username=root --owner=root webSocketDB

migrateup:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/webSocketDB?sslmode=disable" -verbose up


migratedown:
	migrate -path db/migrations -database "postgresql://root:postgres@localhost:5432/webSocketDB?sslmode=disable" -verbose down


#migrate create -ext sql -dir db/migrations -seq <migration-name>

