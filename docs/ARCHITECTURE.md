# Kube-Serverless Architecture

## Overview

Kube-Serverless is a Function-as-a-Service (FaaS) platform built on Kubernetes, providing serverless capabilities with auto-scaling, event-driven triggers, and comprehensive monitoring.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Interfaces                          │
├──────────────┬──────────────────┬─────────────────────────────┤
│     CLI      │    Dashboard UI   │         HTTP API            │
└──────┬───────┴────────┬─────────┴──────────┬──────────────────┘
       │                │                     │
       └────────────────┴─────────────────────┘
                         │
                    ┌────▼────┐
                    │         │
                    │   API   │
                    │ Server  │
                    │         │
                    └────┬────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
    ┌────▼────┐    ┌────▼────┐    ┌────▼────┐
    │Function │    │Function │    │Function │
    │  Node   │    │  Python │    │   Go    │
    │ Runtime │    │ Runtime │    │ Runtime │
    └────┬────┘    └────┬────┘    └────┬────┘
         │               │               │
         └───────────────┼───────────────┘
                         │
                ┌────────┴────────┐
                │                 │
           ┌────▼────┐      ┌────▼────┐
           │Prometheus│      │  HPA    │
           │         │      │Auto-    │
           │Monitoring│      │Scaler   │
           └─────────┘      └─────────┘
```

## Components

### 1. API Server

**Technology**: Go with Gorilla Mux

**Responsibilities**:
- Function lifecycle management (CRUD operations)
- Function invocation routing
- Metrics collection and exposure
- Health checks and readiness probes

**Endpoints**:
- `POST /api/v1/functions` - Deploy function
- `GET /api/v1/functions` - List functions
- `GET /api/v1/functions/{name}` - Get function details
- `PUT /api/v1/functions/{name}` - Update function
- `DELETE /api/v1/functions/{name}` - Delete function
- `POST /api/v1/functions/{name}/invoke` - Invoke function
- `GET /api/v1/functions/{name}/metrics` - Get metrics

### 2. Runtime Servers

Each runtime provides a standardized interface for executing functions:

#### Node.js Runtime
- Base: Node.js 18 Alpine
- Framework: Express
- Features: Dynamic code loading, metric collection

#### Python Runtime
- Base: Python 3.9 Alpine
- Framework: Flask
- Features: Module loading, Prometheus integration

#### Go Runtime
- Base: Go 1.19 Alpine
- Framework: Gorilla Mux
- Features: Plugin support, native performance

**Common Runtime Features**:
- `/health` - Health check endpoint
- `/ready` - Readiness check endpoint
- `POST /` - Function invocation endpoint
- `/metrics` - Prometheus metrics endpoint

### 3. Kubernetes Resources

#### Custom Resource Definition (CRD)
- **Kind**: Function
- **Group**: serverless.kube.io
- **Version**: v1
- Defines function specifications and status

#### Per-Function Resources
- **Deployment**: Manages function pods
- **Service**: Exposes function internally
- **HPA**: Handles auto-scaling
- **ConfigMap**: Stores function code
- **CronJob** (optional): For scheduled triggers

### 4. Auto-Scaling

**Horizontal Pod Autoscaler (HPA)**:
- Metrics: CPU utilization (default 80%)
- Scale to zero capability
- Configurable min/max replicas
- Custom metrics support

**Scaling Behavior**:
```
minReplicas: 0  → Scale to zero when idle
maxReplicas: 10 → Maximum concurrent instances
```

### 5. Monitoring Stack

**Prometheus**:
- Scrapes metrics from API server and functions
- Retention: 30 days (configurable)
- Metrics stored: invocations, duration, cold starts, errors

**Custom Metrics**:
- `function_invocations_total` - Counter
- `function_duration_seconds` - Histogram
- `function_cold_starts_total` - Counter
- `function_deployments_total` - Counter

### 6. Event Triggers

#### HTTP Triggers
- Ingress routes requests to functions
- API server handles routing
- Support for custom paths

#### Cron Triggers
- Kubernetes CronJobs
- Configurable schedules
- Automatic function invocation

#### Message Queue Triggers
- RabbitMQ integration
- Topic-based routing
- Reliable message delivery

### 7. Dashboard UI

**Technology**: React

**Features**:
- Function management interface
- Real-time metrics visualization
- Function deployment wizard
- Test invocation interface

## Data Flow

### Function Deployment
```
CLI/UI → API Server → Kubernetes API
                ↓
        Creates: Deployment, Service, ConfigMap, HPA
                ↓
        Function Pod starts with runtime
```

### Function Invocation
```
HTTP Request → API Server → Function Service → Function Pod
                                                    ↓
                                            Runtime executes code
                                                    ↓
                                            Response ← ← ← ←
```

### Cold Start Flow
```
Request arrives → HPA scales from 0 → 1
                      ↓
              Pod starts (cold start)
                      ↓
              Runtime loads function code
                      ↓
              Request is processed
```

## Security

### RBAC
- ServiceAccount for API server
- ClusterRole with minimal permissions
- Namespace isolation

### Network Policies
- Functions isolated in namespace
- Controlled ingress/egress

### Code Execution
- Functions run in isolated containers
- Resource limits enforced
- No privilege escalation

## Scalability

### Horizontal Scaling
- Functions scale independently
- HPA based on metrics
- Fast scale-up/down

### Resource Efficiency
- Scale to zero when idle
- Shared runtime images
- Lightweight Alpine base images

### Performance Optimizations
- Connection pooling
- Metric caching
- Efficient code loading

## High Availability

- API server runs with 2+ replicas
- Functions can run multiple instances
- Health checks ensure reliability
- Automatic pod recovery

## Cost Optimization

### Scale to Zero
- No cost when function is idle
- Automatic scale-down after timeout

### Resource Limits
- CPU/Memory limits prevent over-provisioning
- Efficient resource utilization

### Metrics for Cost Analysis
- Execution duration tracking
- Invocation counting
- Cost estimation per function

## Extensibility

### Adding New Runtimes
1. Create runtime server (see `runtimes/`)
2. Build Docker image
3. Update API server runtime mapping
4. Deploy and test

### Custom Triggers
1. Create trigger controller
2. Watch Function CRD
3. Create trigger resources
4. Invoke function endpoint

### Custom Metrics
1. Add metric to runtime
2. Expose via `/metrics` endpoint
3. Configure Prometheus scraping
4. Use in HPA configuration

## Future Enhancements

- GraphQL API support
- WebSocket triggers
- Multi-region deployment
- Function versioning
- A/B testing support
- Distributed tracing
- Custom domain mapping
- API Gateway integration
