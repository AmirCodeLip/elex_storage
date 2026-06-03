# Elex Storage Cloud


[![ElexStorage](logo.png)](https://github.com/AmirCodeLip/elex_storage)

elex_storage is a high‑performance, scalable, and fully distributed storage solution designed for modern cloud‑native environments. Built with efficiency and reliability at its core, it provides a robust platform for managing large‑scale data workloads across diverse infrastructures.

- ⚡ High performance and low latency  
- 📈 Horizontally scalable distributed design  
- ☁️ Cloud‑friendly and infrastructure‑agnostic  
- 🧩 Lightweight, modular, and easy to integrate  
- 📜 Open and flexible under the MIT license


## Installation Prerequisites:


PostgreSQL – required as the primary database.

RabbitMQ – required for message queuing.

JWT Key Pair – you need to generate and configure both a private key and a public key for JWT authentication.


## Architecture:

The system follows a microservice architecture. Each service must be run separately, and they communicate with each other.

**🔎 API Gateway:**  
This project is single entry point for clients. And it's automatically map all the backends jrpc and api to the client.

**📝 File Metadata Service:**  
description

**☁️ File Storage Service:**  
description

## Configuration from source:
All environment variables (including database connections and JWT keys) must be set in a configs.yml inside each project but security items like database passwords are set directly from env.

Run each of the following from the project root directory after cloning.

API Gateway:
```sh
go run ./api_gateway/cmd/main.go
```

File Metadata Service
```sh
go run ./file_metadata/cmd/main.go
```

File Storage Service

```sh
cd file_storage
# Copy environment configuration. Secrets such as usernames and passwords are kept in this file.
cp .env.example .env

# Edit configs.yml inside file_storage (manually edit this file) and set db, rabbitmq configs

# Then run the application
go run ./cmd/main.go
```



