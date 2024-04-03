run:
	go run  ./cmd --config=./configs/local.yaml

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root ozinshe

dropdb:
	docker exec -it postgres12 dropdb ozinshe

migrateup:
	migrate -database 'postgresql://root:123@localhost:5432/ozinshe?sslmode=disable' -path db/migration up

migratedown:
	migrate -database 'postgresql://root:123@localhost:5432/ozinshe?sslmode=disable' -path db/migration down

sqlc:
	sqlc generate

sqlcw:
	docker run --rm -v ${CURDIR}:/src -w /src kjconroy/sqlc generate

psql:
	docker exec -it postgres12  psql -U root -d ozinshe

root:
	docker exec -it postgres12  psql -U root