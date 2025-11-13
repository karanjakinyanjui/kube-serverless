# Getting Started with Kube-Serverless

This guide will help you deploy and manage serverless functions on the Kube-Serverless platform.

## Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl configured to access your cluster
- Docker (for building custom runtimes)

## Installation

### 1. Install the Platform

```bash
cd k8s
./install.sh
```

This will:
- Create the `kube-serverless` namespace
- Install Custom Resource Definitions (CRDs)
- Deploy the API server
- Set up monitoring with Prometheus
- Configure RBAC

### 2. Verify Installation

```bash
kubectl get pods -n kube-serverless
```

You should see the API server and Prometheus running.

### 3. Get API Server URL

```bash
kubectl get svc kube-serverless-api -n kube-serverless
```

Note the EXTERNAL-IP or use port-forwarding:

```bash
kubectl port-forward -n kube-serverless svc/kube-serverless-api 8080:80
```

## Using the CLI

### Install CLI

```bash
cd cli
make build
sudo make install
```

### Deploy a Function

#### From a YAML file:

```bash
ksls deploy -f examples/nodejs-hello.yaml
```

#### Using flags:

```bash
ksls deploy my-function \
  --runtime nodejs18 \
  --handler index.handler \
  --code function.js \
  --min-replicas 0 \
  --max-replicas 10
```

### List Functions

```bash
ksls list
```

### Get Function Details

```bash
ksls get my-function
```

### Invoke a Function

```bash
ksls invoke my-function --payload '{"name": "World"}'
```

### View Metrics

```bash
ksls metrics my-function
```

### Delete a Function

```bash
ksls delete my-function
```

## Using the Dashboard

### Deploy the UI

```bash
cd ui
npm install
npm start
```

Or build and deploy to Kubernetes:

```bash
docker build -t kube-serverless-ui:latest .
kubectl apply -f k8s/ui-deployment.yaml
```

Access the dashboard at `http://localhost:3000`

## Writing Functions

### Node.js Example

```javascript
module.exports.handler = async (event) => {
  return {
    statusCode: 200,
    body: JSON.stringify({
      message: 'Hello from Node.js!',
      event: event
    })
  };
};
```

### Python Example

```python
def handler(event):
    return {
        'statusCode': 200,
        'body': {
            'message': 'Hello from Python!',
            'event': event
        }
    }
```

### Go Example

```go
package main

func Handler(event map[string]interface{}) (interface{}, error) {
    return map[string]interface{}{
        "statusCode": 200,
        "body": map[string]interface{}{
            "message": "Hello from Go!",
            "event":   event,
        },
    }, nil
}
```

## Event Triggers

### HTTP Triggers

All functions are automatically exposed via HTTP through the API server.

### Cron Triggers

Add to your function YAML:

```yaml
triggers:
  - type: cron
    config:
      schedule: "*/5 * * * *"  # Every 5 minutes
```

### Message Queue Triggers

Deploy the message queue infrastructure:

```bash
kubectl apply -f k8s/triggers/message-queue.yaml
```

Publish messages to trigger functions:

```bash
# Using RabbitMQ management API
curl -X POST http://rabbitmq:15672/api/exchanges/kube-serverless/function-triggers/publish \
  -H "Content-Type: application/json" \
  -d '{"routing_key": "function.my-function", "payload": "{\"data\": \"test\"}"}'
```

## Auto-Scaling

Functions automatically scale based on:
- CPU utilization (default: 80%)
- Custom metrics
- Scale to zero after inactivity (configurable)

Configure scaling in your function YAML:

```yaml
minReplicas: 0  # Scale to zero
maxReplicas: 50 # Maximum instances
```

## Monitoring

### Prometheus Metrics

Access Prometheus:

```bash
kubectl port-forward -n kube-serverless svc/prometheus 9090:9090
```

Visit `http://localhost:9090`

### Available Metrics

- `function_invocations_total` - Total function invocations
- `function_duration_seconds` - Function execution duration
- `function_cold_starts_total` - Cold start count
- `function_deployments_total` - Deployment count

### View Logs

```bash
kubectl logs -l function=my-function -n kube-serverless -f
```

## Best Practices

1. **Keep functions small and focused** - Single responsibility principle
2. **Handle cold starts** - First invocation may be slower
3. **Use environment variables** - For configuration
4. **Monitor metrics** - Track performance and costs
5. **Set appropriate scaling limits** - Balance cost and performance
6. **Test locally** - Before deploying to production
7. **Use version control** - For function code and configurations

## Troubleshooting

### Function won't deploy

- Check API server logs: `kubectl logs -n kube-serverless deployment/kube-serverless-api`
- Verify function YAML syntax
- Ensure runtime is supported

### Function not scaling

- Check HPA status: `kubectl get hpa -n kube-serverless`
- Verify metrics-server is installed
- Check resource requests/limits

### Cold starts taking too long

- Increase minReplicas to keep instances warm
- Optimize function initialization code
- Use lighter base images

## Next Steps

- Explore [examples/](../examples/) for more function templates
- Read [ARCHITECTURE.md](./ARCHITECTURE.md) to understand the platform
- Check [API.md](./API.md) for API documentation
