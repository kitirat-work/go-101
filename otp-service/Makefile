DB_URL = mysql://myuser:mypassword@tcp(127.0.0.1:6033)/mydatabase

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./db/migrations $$name

migrate-up:
	migrate -database "$(DB_URL)" -path ./db/migrations -verbose up

migrate-down:
	@read -p "Enter number of migrations to roll back: " n; \
	migrate -database "$(DB_URL)" -path ./db/migrations -verbose down $$n

migrate-force:
	@read -p "Enter migration version to force: " version; \
	migrate -database "$(DB_URL)" -path ./db/migrations -verbose force $$version

create-db:
	docker compose up -d mysql