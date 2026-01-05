# Assignment 2 - Concurrent Task Service

**Student:** Nurzhan Izimbetov  
**Group:** SE-2436

## Project Structure

```
Assignment2/
├── go.mod                      # Go module file
├── main.go                     # Application entry point
├── internal/
│   ├── api/
│   │   └── handler.go         # HTTP request handlers
│   ├── queue/
│   │   └── queue.go           # Task queue implementation
│   ├── worker/
│   │   └── worker.go          # Worker pool and monitoring
│   ├── store/
│   │   └── repository.go      # Generic in-memory repository
│   └── model/
│       └── task.go            # Data models
└── README.md                   # This file
```

## How to Run

```bash
go run .
```

The server will start on port 8080.

## API Endpoints

### 1. POST /tasks
Create a new task.

**Request:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"payload":"test task"}'
```

**Response (201):**
```json
{
  "id": "1",
  "payload": "test task",
  "status": "PENDING"
}
```

### 2. GET /tasks
Get all tasks.

**Request:**
```bash
curl http://localhost:8080/tasks
```

**Response (200):**
```json
[
  {
    "id": "1",
    "payload": "test task",
    "status": "DONE"
  }
]
```

### 3. GET /tasks/{id}
Get a specific task by ID.

**Request:**
```bash
curl http://localhost:8080/tasks/1
```

**Response (200):**
```json
{
  "id": "1",
  "payload": "test task",
  "status": "DONE"
}
```

**Response (404):** Task not found

### 4. GET /stats
Get server statistics.

**Request:**
```bash
curl http://localhost:8080/stats
```

**Response (200):**
```json
{
  "submitted": 10,
  "completed": 7,
  "in_progress": 1
}
```

## Testing the Service

### Test Scenario 1: Basic Task Flow
```bash
# create task
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"payload":"Task 1"}'

# get all
curl http://localhost:8080/tasks

# wait 3 seconds
sleep 3

# Get all tasks (status should be DONE)
curl http://localhost:8080/tasks
```

### Test Scenario 2: Multiple Tasks
```bash
# submit 5 tasks
for i in {1..5}; do
  curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d "{\"payload\":\"Task $i\"}"
done

# check stats
curl http://localhost:8080/stats
```

### Test Scenario 3: Graceful Shutdown
```bash
# submit tasks
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"payload":"Task"}'

# press ctrl+C in terminal running the server
```
