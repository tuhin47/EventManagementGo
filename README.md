# Event Management System API

This README provides example `curl` requests for interacting with the Event Management System API. You can use these requests in Postman or directly in your terminal.

## Base URL

The server runs on `http://localhost:8080`.

## Endpoints

### 1. Welcome Endpoint

**Description:** Displays a welcome message.

**Request:**
```bash
curl -X GET http://localhost:8080/
```

### 2. Create Event Endpoint

**Description:** Creates a new event.

**Request:**
```bash
curl -X POST http://localhost:8080/create \
     -H "Content-Type: application/json" \
     -d '{
           "name": "Event Name", 
           "description": "Event Description", 
           "location": "Event Location", 
           "start_time": "2025-05-10T10:00:00Z", 
           "end_time": "2025-05-10T12:00:00Z", 
           "organizer": "Organizer Name", 
           "capacity": 100
         }'
```

### 3. Get All Events Endpoint

**Description:** Retrieves all events.

**Request:**
```bash
curl -X GET http://localhost:8080/events \
     -H "Content-Type: application/json"
```

## Notes
- Replace `Event Name`, `2025-05-10`, and `Event Location` with actual event details in the `/create` endpoint.
- Ensure the server is running before making requests.