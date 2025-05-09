# Event Management System

## Quick Run (Docker)

```bash
docker compose --env-file .env -f docker-compose.yml -p eventmanagement up -d --build
```
**Note** : Make sure to update port from `.env` if your local mysql is already running on port `3306`

It will run the postman collection from the `postman_collection` folder. 

## API Endpoints

The application provides several API endpoints for managing events.

| Method | Endpoint      | Description                           | Payload                                                                                                                                                                            |
|--------|---------------|---------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GET    | `/`           | Checks if the server is running.      | None                                                                                                                                                                               |
| POST   | `/events`     | Creates a new event.                  | ```json { "title": "Event Title", "description": "Event Description", "start_time": "2025-05-10T10:00:00Z", "end_time": "2025-05-10T12:00:00Z", "created_by": "Creator Name" } ``` |
| GET    | `/events`     | Retrieves all events.                 | None                                                                                                                                                                               |
| GET    | `/event/{id}` | Retrieves a specific event by its ID. | None                                                                                                                                                                               |
| PUT    | `/event/{id}` | Updates an existing event by its ID.  | ```json { "title": "Updated Event Title" } ```                                                                                                                                     |
| DELETE | `/event/{id}` | Deletes an existing event by its ID.  | None                                                                                                                                                                               |

## Postman 

 - [Postman Collection](./postman_collection/GoLangProgram.postman_collection.json)
 - [Environment File](./postman_collection/GoLangProgram.postman_environment.json)

## Gitpod
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/tuhin47/EventManagementGo)

