# CRUD-home-appliance-store
Create CRUD for home appliance store with swagger documentation

# Install
## create .env file in project root directiory
## run dokcer-compose up

# Enviromant variable
## example
```
env=local
POSTGRES_HOST=postgres-db
POSTGRES_PORT=5432
POSTGRES_USER=root
POSTGRES_PASSWORD=root
POSTGRES_DB=service-crud-db
timeout=2s
http_server:
    address:crud-service
    port:8080
```