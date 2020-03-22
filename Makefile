.UP:db-start migrate-up
db-start:
	@echo "-----running db-create database test_authorization-----"	
	docker run --rm --name test_postgres -v $(shell pwd)/migrations:/migrations -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=test_authorization -d -p 5432:5432 postgres
	sleep 3
migrate-up:
	#make migrate-reset
	@echo "-----running migrate-up in database test_authorization-----"
	docker exec -it test_postgres psql -d test_authorization -p 5432 -h localhost -U postgres -f /migrations/migrate_up.sql
.DOWN:migrate-down db-stop
migrate-down:
	@echo "-----running migrate-down in database test_authorization-----"
	docker exec -it test_postgres psql -d test_authorization -p 5432 -h localhost -U postgres -f /migrations/migrate_down.sql
db-stop:
	@echo "-----running db-stop database test_authorization-----"	
	docker stop test_postgres
.TEST:test
test:
	go test -v -race -timeout 30s ./...	