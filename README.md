## Requirements
- Docker
- Docker Compose

## How to Run
1. Copy environment file
   cp .env.example .env.docker

2. Start all services
   docker compose up

3. Run database migration
   docker compose run --rm server ./migrate up

4. Run mock MQTT publisher
   docker compose --profile mock up mqtt-publisher --build