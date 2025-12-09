# Kube-Serverless API Documentation

Base URL: `http://<api-server-address>/api/v1`

## Authentication

Currently, the API is unauthenticated. For production use, implement:
- API keys
- OAuth 2.0
- Kubernetes RBAC integration

## Endpoints

### List Functions

```http
GET /functions
```

**Response**:
```json
[
  {
    "name": "hello-world",
    "runtime": "nodejs18",
    "handler": "index.handler",
    "minReplicas": 0,
    "maxReplicas": 10,
    "status": {
      "state": "running",
      "replicas": 2,
      "endpoint": "hello-world.kube-serverless.svc.cluster.local",
      "lastDeployment": "2024-01-15T10:30:00Z"
    }
  }
]
```

### Create Function

```http
POST /functions
```

**Request Body**:
```json
{
  "name": "my-function",
  "runtime": "nodejs18",
  "handler": "index.handler",
  "code": "module.exports.handler = async (event) => { return { statusCode: 200, body: 'Hello!' }; };",
  "environment": {
    "VAR1": "value1"
  },
  "minReplicas": 0,
  "maxReplicas": 10,
  "triggers": [
    {
      "type": "http",
      "config": {
        "path": "/my-function"
      }
    }
  ]
}
```

**Response**: `201 Created`
```json
{
  "name": "my-function",
  "runtime": "nodejs18",
  "handler": "index.handler",
  "status": {
    "state": "deploying"
  }
}
```

### Get Function

```http
GET /functions/{name}
```

**Response**: `200 OK`
```json
{
  "name": "my-function",
  "runtime": "nodejs18",
  "handler": "index.handler",
  "code": "...",
  "minReplicas": 0,
  "maxReplicas": 10,
  "status": {
    "state": "running",
    "replicas": 1,
    "endpoint": "my-function.kube-serverless.svc.cluster.local"
  }
}
```

### Update Function

```http
PUT /functions/{name}
```

**Request Body**: Same as Create Function

**Response**: `200 OK`

### Delete Function

```http
DELETE /functions/{name}
```

**Response**: `204 No Content`

### Invoke Function

```http
POST /functions/{name}/invoke
```

**Request Body**:
```json
{
  "key": "value",
  "data": [1, 2, 3]
}
```

**Response Headers**:
- `X-Function-Duration`: Execution time in seconds
- `X-Cold-Start`: `true` or `false`

**Response**: `200 OK`
```json
{
  "statusCode": 200,
  "body": {
    "message": "Function executed successfully",
    "result": "..."
  }
}
```

### Get Function Metrics

```http
GET /functions/{name}/metrics
```

**Response**: `200 OK`
```json
{
  "invocations": 1523,
  "coldStarts": 12,
  "avgDuration": 0.234,
  "errorRate": 0.02,
  "costEstimate": 0.15
}
```

## Health Endpoints

### Health Check

```http
GET /health
```

**Response**: `200 OK`
```json
{
  "status": "healthy"
}
```

### Readiness Check

```http
GET /ready
```

**Response**: `200 OK`
```json
{
  "status": "ready"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 404 Not Found
```json
{
  "error": "Function not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error: ..."
}
```

## Function Specification

### Runtimes

- `nodejs18` - Node.js 18 (Alpine)
- `python39` - Python 3.9 (Alpine)
- `go119` - Go 1.19 (Alpine)

### Handler Format

- **Node.js**: `filename.exportname` (e.g., `index.handler`)
- **Python**: `filename.functionname` (e.g., `handler.handler`)
- **Go**: `FunctionName` (e.g., `Handler`)

### Trigger Types

#### HTTP Trigger
```json
{
  "type": "http",
  "config": {
    "path": "/custom-path"
  }
}
```

#### Cron Trigger
```json
{
  "type": "cron",
  "config": {
    "schedule": "*/5 * * * *"
  }
}
```

#### Queue Trigger
```json
{
  "type": "queue",
  "config": {
    "queue": "function-queue",
    "exchange": "function-triggers"
  }
}
```

## Event Object

Functions receive an event object with the following structure:

```javascript
{
  "body": {},           // Request body (parsed JSON)
  "headers": {},        // Request headers
  "method": "POST",     // HTTP method
  "path": "/invoke",    // Request path
  "query": {}           // Query parameters
}
```

## Response Format

Functions should return:

```javascript
{
  "statusCode": 200,
  "body": "Response data (string or object)"
}
```

Or simply return data directly:

```javascript
{
  "message": "Hello World"
}
```

## Rate Limiting

Not currently implemented. Consider implementing:
- Per-function rate limits
- Per-user rate limits
- Global rate limits

## Versioning

API version is included in the path: `/api/v1/...`

Future versions will be `/api/v2/...`, etc.

## SDK Examples

### JavaScript/Node.js

```javascript
const axios = require('axios');

const client = axios.create({
  baseURL: 'http://api-server/api/v1'
});

// Deploy function
await client.post('/functions', {
  name: 'my-function',
  runtime: 'nodejs18',
  handler: 'index.handler',
  code: 'module.exports.handler = ...'
});

// Invoke function
const result = await client.post('/functions/my-function/invoke', {
  data: 'test'
});
```

### Python

```python
import requests

BASE_URL = 'http://api-server/api/v1'

# Deploy function
response = requests.post(f'{BASE_URL}/functions', json={
    'name': 'my-function',
    'runtime': 'python39',
    'handler': 'handler',
    'code': 'def handler(event): ...'
})

# Invoke function
result = requests.post(f'{BASE_URL}/functions/my-function/invoke',
    json={'data': 'test'})
```

### cURL

```bash
# Deploy function
curl -X POST http://api-server/api/v1/functions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-function",
    "runtime": "nodejs18",
    "handler": "index.handler",
    "code": "module.exports.handler = async (e) => ({statusCode: 200, body: \"OK\"});"
  }'

# Invoke function
curl -X POST http://api-server/api/v1/functions/my-function/invoke \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
```
