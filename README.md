# Distributed File Processing System

A high-performance, event-driven microservices architecture built with Go v1.21+, gRPC, Redis, and PostgreSQL, fully containerized with Docker and orchestrated on Kubernetes.

## System Architecture

The system utilizes an asynchronous "Producer-Consumer" pattern to ensure high availability:

1. Upload Service: Receives files via gRPC, saves them to a Persistent Volume, and publishes tasks to Redis.

2. Worker Service: Consumes tasks from the Redis queue and performs file analysis (MIME type, word count, and SHA256 hashing).

3. Metadata Service: Manages the source of truth for file states (Pending → Processing → Completed).

4. Result Service: Persists the final analytical output into PostgreSQL.

5. Infrastructure: PostgreSQL for relational metadata and Redis for the high-speed task queue.

## Tech Stack

- Language: Go (Golang) v1.21+

- Communication: gRPC / Protocol Buffers

- Database: PostgreSQL 18

- Queue/Cache: Redis 7

- Orchestration: Kubernetes / Docker

## Getting Started

### 1. Kubernetes Deployment

Apply all manifests in the correct order to ensure infrastructure is available before the application services start:

Bash

kubectl apply -f k8s/infrastructure.yaml
kubectl apply -f k8s/databases.yaml
kubectl apply -f k8s/services.yaml
kubectl apply -f k8s/upload.yaml 2. Database Initialization
Since the cluster starts fresh, you must create the database and schema manually or via migrations:

Bash

# Create the database

kubectl exec -it <postgres-pod-name> -- psql -U file_processor -d postgres -c "CREATE DATABASE metadata_db;"

# Apply the schema

kubectl exec -i <postgres-pod-name> -- psql -U file_processor -d metadata_db <<EOF
CREATE TABLE files (
id UUID PRIMARY KEY,
filename TEXT NOT NULL,
size BIGINT NOT NULL,
mime_type TEXT NOT NULL,
status TEXT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE results (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
output_data TEXT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
EOF
🐳 Docker Support
Build Images
Each service includes a multi-stage Dockerfile:

Bash

docker build -t your-username/upload-service:v1 -f build/package/upload.Dockerfile .

# Repeat for worker, metadata, and result services

⚡️ Stress Testing
A script is provided to test the system's concurrency by launching 10 simultaneous upload requests.

Bash

# Run the stress test

./stress_test.sh
Script Content (stress_test.sh):

Bash

#!/bin/bash
for i in {1..10}
do
go run cmd/client/main.go testdata/sample.txt &
done
wait
echo "All uploads sent!"
📊 Verification & Monitoring
Check Logs
Watch the workers process the stress test tasks in real-time:

Bash

kubectl logs -l app=worker -f --tail=20
Query Results
Verify that all 10 files were processed and saved:

Bash

kubectl exec -it <postgres-pod-name> -- psql -U file_processor -d metadata_db -c "SELECT \* FROM results;"
✅ Features
Asynchronous Workflow: Decoupled file uploads from heavy processing.

Auto-Healing: Kubernetes manages pod restarts and service discovery.

Concurrent Scaling: Multiple worker pods can process the Redis queue simultaneously.

Data Integrity: Foreign key constraints between file metadata and analysis results.
