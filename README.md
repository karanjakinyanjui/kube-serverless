# Kube-Serverless Platform

A production-ready Function-as-a-Service (FaaS) platform built on Kubernetes, providing serverless capabilities with auto-scaling, event-driven triggers, and comprehensive monitoring.

## Features

- **Function Deployment** - Deploy via API, CLI, or web dashboard
- **Auto-Scaling to Zero** - Cost-efficient scaling based on demand
- **Multiple Runtimes** - Node.js 18, Python 3.9, Go 1.19
- **Event-Driven Triggers** - HTTP, Cron schedules, message queues
- **Comprehensive Monitoring** - Prometheus metrics, cold start tracking, cost estimation
- **Management Dashboard** - React-based UI for function management
- **Production Ready** - RBAC, health checks, resource limits

## Architecture

```
┌─────────────┬──────────────┬─────────────┐
│     CLI     │  Dashboard   │  HTTP API   │
└──────┬──────┴──────┬───────┴──────┬──────┘
       │             │              │
       └─────────────┴──────────────┘
                     │
              ┌──────▼──────┐
              │  API Server │
              └──────┬──────┘
                     │
       ┌─────────────┼─────────────┐
       │             │             │
  ┌────▼───┐    ┌───▼────┐   ┌───▼────┐
  │Node.js │    │ Python │   │   Go   │
  │Runtime │    │Runtime │   │Runtime │
  └────────┘    └────────┘   └────────┘
```

## Quick Start

### Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl configured
- Docker (for building images)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/karanjakinyanjui/kube-serverless.git
cd kube-serverless
```

2. Install the platform:
```bash
cd k8s
./install.sh
```

3. Install the CLI:
```bash
cd cli
make build
sudo make install
```

### Deploy Your First Function

Create a function file `hello.js`:
```javascript
module.exports.handler = async (event) => {
  return {
    statusCode: 200,
    body: JSON.stringify({ message: 'Hello from Kube-Serverless!' })
  };
};
```

Deploy it:
```bash
ksls deploy hello-world \
  --runtime nodejs18 \
  --handler index.handler \
  --code hello.js
```

Invoke it:
```bash
ksls invoke hello-world --payload '{}'
```

## Use Cases

### Webhook Handlers
Handle incoming webhooks from external services with auto-scaling based on traffic.

### Scheduled Data Processing
Process data on a schedule using cron triggers, scaling to zero between runs.

### API Endpoints
Deploy lightweight API endpoints that scale automatically based on demand.

### Event Processing
React to messages from queues and process events asynchronously.

## Components

### API Server
- Go-based REST API
- Function lifecycle management
- Metrics collection
- Health monitoring

### CLI Tool (`ksls`)
- Deploy, list, get, delete functions
- Invoke functions
- View metrics and logs
- YAML-based configuration

### Dashboard UI
- React-based web interface
- Real-time metrics visualization
- Function management
- Test invocation interface

### Runtimes
- **Node.js 18**: Express-based runtime with dynamic code loading
- **Python 3.9**: Flask-based runtime with module loading
- **Go 1.19**: Native performance with plugin support

### Monitoring
- Prometheus metrics collection
- Cold start tracking
- Execution duration histograms
- Cost estimation per function

## Documentation

- [Getting Started Guide](docs/GETTING_STARTED.md) - Detailed setup and usage
- [Architecture Documentation](docs/ARCHITECTURE.md) - System design and components
- [API Reference](docs/API.md) - Complete API documentation

## Examples

Check the [examples/](examples/) directory for:
- Node.js HTTP function
- Python scheduled data processor
- Go webhook handler

## Project Structure

```
kube-serverless/
├── api/              # API server (Go)
├── cli/              # CLI tool (Go)
├── k8s/              # Kubernetes manifests
│   ├── triggers/     # Event trigger configurations
│   └── install.sh    # Installation script
├── runtimes/         # Function runtimes
│   ├── nodejs/       # Node.js runtime
│   ├── python/       # Python runtime
│   └── go/           # Go runtime
├── ui/               # Dashboard (React)
├── examples/         # Example functions
├── docs/             # Documentation
└── Makefile          # Build automation
```

## Key Features Explained

### Auto-Scaling to Zero
Functions automatically scale down to zero replicas when idle, eliminating costs during periods of no activity. When a request arrives, Kubernetes scales up the function within seconds.

### Event-Driven Triggers
- **HTTP**: Direct HTTP invocation via API or ingress
- **Cron**: Scheduled execution using Kubernetes CronJobs
- **Message Queues**: RabbitMQ integration for async processing

### Monitoring & Metrics
Track key metrics for each function:
- Total invocations
- Cold start count
- Average execution duration
- Error rate
- Estimated cost

### Cost Efficiency
- Scale to zero when idle
- Resource limits prevent over-provisioning
- Per-function cost tracking
- Efficient resource utilization

## Development

Build all components:
```bash
make all
```

Run locally with Docker Compose:
```bash
docker-compose up
```

Run tests:
```bash
make test
```

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests.

## License

MIT License - see LICENSE file for details

## Showcase Value

This project demonstrates:
- **Modern Cloud-Native Architecture** - Kubernetes, containers, microservices
- **Serverless Computing** - FaaS implementation with auto-scaling
- **Full-Stack Development** - Go backend, React frontend, CLI tools
- **DevOps Best Practices** - IaC, monitoring, CI/CD-ready
- **Production Readiness** - RBAC, health checks, resource management
- **Innovation** - Trendy technology stack showcasing technical depth

## Support

For issues, questions, or contributions:
- Open an issue on GitHub
- Check the documentation in [docs/](docs/)
- Review examples in [examples/](examples/)

## Roadmap

- [ ] GraphQL API support
- [ ] WebSocket triggers
- [ ] Multi-region deployment
- [ ] Function versioning and rollback
- [ ] A/B testing support
- [ ] Distributed tracing integration
- [ ] Custom domain mapping
- [ ] API Gateway integration

---

Built with ❤️ using Kubernetes, Go, React, and Node.js
