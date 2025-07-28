# CRUD-home-appliance-store
Create CRUD for home appliance store with swagger documentation

# 🚀 Run
- create .env file in project root directiory
- run dokcer-compose -f project/docker-compose.yaml up -d --build

### Enviromant variable
#### Example:
```
# common variable
env=local

# postgres varibale
POSTGRES_HOST=postgres-db
POSTGRES_PORT=5432
POSTGRES_USER=root
POSTGRES_PASSWORD=root
POSTGRES_DB=service-crud-db
timeout=5s

# crud variable
crud_service_id=go-service-0001
crud_service_name=crud-service
crud_service_address=crud-service
crud_service_port=8080

# consul variable
consul_service_address=consul-service
consul_service_port=8500
consul_service_retry_delay=2s
consul_service_max_attempts=5
consul_service_survey_interval=10s
consul_service_survey_timeout=5s
```

# 🧪 Endpoints

| Method | URL                             | Auth | Description                     |
|--------|---------------------------------|------|---------------------------------|
| POST   | `/api/v1/products`              | 🔓   | create product                  |
| GET    | `/api/v1/products `             | 🔓   | get all products                |
| GET    | `/api/v1/products/:id`          | 🔓   | get product by id               |
| PATCH  | `/api/v1/products/:id?decrease=`| 🔓   | update product available stock  |
| DELETE | `/api/v1/products/:id`          | 🔓   | delete product by id            |
|--------|---------------------------------|------|---------------------------------|
| POST   | `/api/v1/images`                | 🔓   | create images                   |
| GET    | `/api/v1/images `               | 🔓   | get all images                  |
| GET    | `/api/v1/images/:id`            | 🔓   | get images by id                |
| PATCH  | `/api/v1/images/:id`            | 🔓   | update images available stock   |
| DELETE | `/api/v1/images/:id`            | 🔓   | delete images by id             |
|--------|---------------------------------|------|---------------------------------|
| POST   | `/api/v1/clients`               | 🔓   | create product                  |
| GET    | `/api/v1/clients `              | 🔓   | get all clients                 |
| GET    | `/api/v1/clients/:id`           | 🔓   | get client by id                |
| PATCH  | `/api/v1/clients/:id?decrease=` | 🔓   | update client available stock   |
| DELETE | `/api/v1/clients/:id`           | 🔓   | delete client by id             |
|--------|---------------------------------|------|---------------------------------|
| POST   | `/api/v1/suppliers`             | 🔓   | create suppplier                |
| GET    | `/api/v1/suppliers`             | 🔓   | get all suppliers               |
| GET    | `/api/v1/suppliers/:id`         | 🔓   | get supplier by id              |
| PATCH  | `/api/v1/suppliers/:id?decrease=`| 🔓   | update supplier available stock|
| DELETE | `/api/v1/suppliers/:id`         | 🔓   | delete supplier by id           |

## Tech stack
  
- Go — language
- Gin — HTTP-engine
- PGX — PostgreSQL-driver
- PostgreSQL — data base
- Docker / Docker Compose — build and run
- slog — logger
- consul - service registration
- swaggo - swagger generator