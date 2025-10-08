# LogiQ - Logic Quiz Platform

LogiQ is a modern, full-stack platform for intelligent logic quiz. Its core capability is automatically generating diverse logic reasoning questions from predefined rules and configurations. The platform includes a cross-platform client for learners, a powerful backend service, and a web-based admin console.

## Core Features

* **Intelligent question generation**: The backend ships with a question generation engine that dynamically creates logic questions from configuration.
* **Cross-platform learner client**: Built with Flutter, a single codebase runs on iOS, Android, the web, and desktop to deliver a consistent experience.
* **Visual admin console**: A modern React + Vite web interface helps administrators manage the question bank, trigger generation, and inspect metrics.
* **Containerized deployment**: The backend stack (Go service, MongoDB, Redis) is orchestrated with Docker and Docker Compose for one-command startup and easy deployments.

## System Architecture

The LogiQ platform uses a decoupled front-end/back-end architecture with three main components:

1. **`frontend` (learner client)**: A Flutter mobile/web application where learners answer questions, review results, and track personal statistics.
2. **`admin_site` (admin console)**: A React-powered web admin dashboard for question management and data exploration.
3. **`backend` (API service)**: A Go (Gin) API server responsible for:
   * Handling requests from both the learner client and the admin console.
   * Executing the core question generation logic.
   * Interacting with MongoDB for persistence and Redis for caching.
   * Managing user authentication and authorization.

```
+----------------+      +----------------+
|  Flutter App   |      |   React Admin  |
|   (frontend)   |      |   (admin_site) |
+-------+--------+      +--------+-------+
        |                        |
        |         HTTP/S         |
        +-----------+------------+
                    |
+-------------------+-------------------+
|              Go Backend               |
|                (Gin)                  |
+-------------------+-------------------+
|         |         |                   |
+---------+--+   +--+---------+   +-----+------+
|  MongoDB   |   |   Redis    |   |  Generation  |
| (Database) |   |  (Cache)   |   |    Engine    |
+------------+   +------------+   +--------------+
```

## Tech Stack

* **Backend**: Go, Gin
* **Learner client**: Flutter, Dart
* **Admin console**: React, TypeScript, Vite
* **Database**: MongoDB
* **Cache**: Redis
* **Deployment**: Docker, Docker Compose

## Quick Start

### 1. Launch the backend services (Docker)

Create a `.env` file inside the `backend/` directory and fill in the required environment variables.

**`.env` example:**
```env
MONGO_URI=mongodb://mongodb:27017
REDIS_ADDR=redis:6379
GIN_MODE=debug
```

From the project root, start all backend services with:
```bash
docker-compose -f docker/docker-compose.yml up --build -d
```

### 2. Launch the admin console (React)

```bash
cd admin_site
npm install
npm run dev
```
The admin console runs at `http://localhost:5173` by default.

### 3. Launch the learner client (Flutter)

```bash
cd frontend
flutter pub get
flutter run
```
The Flutter application runs on the emulator, physical device, or browser you choose.
