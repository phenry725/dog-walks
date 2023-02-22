run:
	docker compose up -d
purge:
	docker compose down --rmi all --volumes