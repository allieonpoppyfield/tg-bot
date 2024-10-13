.run-d: 
	$(info Деплой приложения в докер...)
	docker-compose build
	docker compose up -d

rund: .run-d


migrate-up: 
	goose -dir ./migrations postgres "user=beslan dbname=buddymap password=beslan sslmode=disable host=localhost" up

migrate-down: 
	goose -dir ./migrations postgres "user=beslan dbname=buddymap password=beslan sslmode=disable host=localhost" down

