# Event Trigger Platform

A Golang-based API to create and manage event triggers (scheduled & API-based).

## 🚀 Features
- Create, Edit, Delete, and Execute Triggers
- Schedule events using cron
- API triggers with JSON payloads
- Event logs with 48-hour lifecycle
- Redis caching for logs
- **Dockerized & Deployed on Fly.io**

---

## 🛠 Tech Stack
- **Golang + Gin** (API)
- **PostgreSQL** (Database)
- **Redis** (Caching)
- **Docker & Fly.io** (Deployment)

---

## 📖 **API Documentation**
### 🔹 **Swagger UI**
- **Local**: `http://localhost:4000/swagger/index.html`
- **Deployed**: `https://event-trigger-app-bitter-bird-4607.fly.dev/swagger/index.html`

### 🔹 **Endpoints**
#### **Trigger Management**
| Method | Endpoint | Description |
|--------|---------|-------------|
| `POST` | `/api/triggers` | Create a new trigger |
| `GET` | `/api/triggers` | Get all triggers |
| `GET` | `/api/triggers/:id` | Get trigger details |
| `PUT` | `/api/triggers/:id` | Update a trigger |
| `DELETE` | `/api/triggers/:id` | Delete a trigger |

#### **Trigger Execution**
| Method | Endpoint | Description |
|--------|---------|-------------|
| `POST` | `/api/triggers/:id/execute` | Manually execute a trigger |

#### **Event Logs**
| Method | Endpoint | Description |
|--------|---------|-------------|
| `GET` | `/api/events` | Get recent event logs (last 2 hours) |
| `GET` | `/api/events/archived` | Get archived event logs |
| `DELETE` | `/api/events/purge` | Purge old events |


## 🛠 Local Setup

### 1️⃣ Prerequisites
- Install **Docker**: [Docker Install Guide](https://docs.docker.com/get-docker/)

### 2️⃣ Clone Repository
```sh
git clone https://github.com/shivamsouravjha/segwise.git
cd segwise
```

### 3️⃣ Run with Docker
```sh
docker compose up --build
```

The API will run on `http://localhost:4000`.

---

## 📡 API Endpoints

### **1️ Register a New User**
```sh
curl -X POST "http://localhost:4000/api/triggers"      -H "Content-Type: application/json"      -d '{
           "type": "scheduled",
           "schedule": "*/10 * * * *", 
           "endpoint": null,
           "payload": null
         }'
```

### **2️ Login and Get JWT Token**
```sh
curl -X POST "http://localhost:4000/api/triggers"      -H "Content-Type: application/json"      -d '{
           "type": "api",
           "schedule": null,
           "endpoint": "http://example.com/webhook",
           "payload": "{ \"message\": \"Hello, world!\" }"
         }'
```

### **3️ Create a Scheduled Trigger**
```sh
curl -X POST "http://localhost:4000/api/triggers" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
     -d '{
           "type": "scheduled",
           "schedule": "in 10 seconds", 
           "endpoint": null,
           "payload": null,
           "oneTime":false,
         }'
```

### **4️ Create an API Trigger**
```sh
curl -X POST "http://localhost:4000/api/triggers" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
     -d '{
           "type": "api",
           "schedule": null,
           "endpoint": "http://example.com/webhook",
           "payload": "{ \"message\": \"Hello, world!\" }"
         }'
```

### **5️ Get All Triggers**
```sh
curl -X GET "http://localhost:4000/api/triggers" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### **6 Get a Trigger by ID**
```sh
curl -X GET "http://localhost:4000/api/triggers/{trigger_id}" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### **7 Update a Trigger**
```sh
curl -X PUT "http://localhost:4000/api/triggers/{trigger_id}" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
     -d '{
           "type": "scheduled",
           "schedule": "*/5 * * * *", 
           "endpoint": null,
           "payload": null
         }'
```

### **8 Delete a Trigger**
```sh
curl -X DELETE "http://localhost:4000/api/triggers/{trigger_id}" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### **9 Execute a Trigger Manually**
```sh
curl -X POST "http://localhost:4000/api/triggers/{trigger_id}/execute" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### **10 Get All Active Events**
```sh
curl -X GET "http://localhost:4000/api/events" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### **11 Get All Archived Events**
```sh
curl -X GET "http://localhost:4000/api/events/archived" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### **12 Purge Old Archived Events**
```sh
curl -X DELETE "http://localhost:4000/api/events/purge" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### **13 Test a One-Time Scheduled Trigger**
```sh
curl -X POST "http://localhost:4000/api/triggers/test/scheduled" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
     -d '{
           "delay": 5
         }'
```

### **14 Test a One-Time API Trigger**
```sh
curl -X POST "http://localhost:4000/api/triggers/test/api" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
     -d '{
           "endpoint": "http://example.com/webhook",
           "payload": { "test": "API payload test" }
         }'
```


---

## 🚀 Deployment (Fly.io)

The project is deployed at:  
🔗 **[Event Trigger App](https://event-trigger-app-bitter-bird-4607.fly.dev/swagger/index.html)**

---

## 💰 Cost Estimation
Fly.io Free Tier Usage:

* Compute (256MB RAM, 1 vCPU)	Free
* PostgreSQL (LiteDB)	Free
* Redis (Upstash)	Free
* Bandwidth (5GB limit)	Free
* Estimated Cost for 30 Days (24x7, 5 queries/day): $0 (Free Tier)


## ✅ **Testing**

Run tests:
```sh
go test -v ./... -cover
```

---

## 📌 **Assumptions**
- Redis is **only for caching event logs**.
- API triggers **use flat JSON payloads**.
- Events **expire in 48 hours** (2 hours active, 46 hours archived).
- Fly.io **free tier** is used, so **scaling is limited**.

---


## 📌 **Credits where due**
- ChatGPT and CoPilot(dev related help)
- Keploy VSCode extension to generate testcases

---
