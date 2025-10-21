## Project Structure

ecommerce-platform/
├── services/
│   ├── user-service/
│   ├── product-service/
│   ├── order-service/
│   ├── payment-service/
│   ├── inventory-service/
│   ├── notification-service/
│   └── api-gateway/
├── proto/              # gRPC Proto files (shared)
├── pkg/                # Shared libraries
│   ├── database/
│   ├── redis/
│   ├── kafka/
│   └── middleware/
├── docker-compose.yml
└── k8s/               # Kubernetes manifests

### **technology:**
- **Golang** (Gin, gRPC, GraphQL)
- **Postgres** (แต่ละ service มี DB ของตัวเอง)
- **Redis** (Cache, Session)
- **Kafka** (Event-driven communication)
- **Docker** + **Docker Compose**
- **Kubernetes** (Deploy production)

### **Microservices: **
1. **User Service** - manage users, authentication (gRPC)
2. **Product Service** - (REST API)
3. **Order Service** - (gRPC) 
4. **Payment Service** - (gRPC)
5. **Inventory Service** - (gRPC)
6. **Notification Service** - (Kafka Consumer)
7. **API Gateway** - GraphQL + REST gateway
