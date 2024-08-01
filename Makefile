test:
	docker compose -f docker-compose.test.yaml up

gen-mocks:
	mockery .
