.run: 
	$(info Запуск приложения...)
	go run src/cmd/main.go

run: .run 

compose:
	docker-compose -p buddymap_proj up --build

compose-re:
	docker-compose -p buddymap_proj up --d --force-recreate
