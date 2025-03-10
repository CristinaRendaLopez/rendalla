# Rendalla
**Web platform for managing and distributing sheet music.**  

## Features
- Upload and store sheet music in different formats.
- Search and download music sheets for various instruments.
- Secure admin panel with authentication.

## Project Structure

The repository is organized into different folders to separate the backend, frontend, and infrastructure.

- **backend/** → Contains the **Go** backend code, which runs on **AWS Lambda** and integrates with **DynamoDB** and **S3**.
  - `handlers/` → HTTP request handlers.
  - `models/` → Data structures and models.
  - `services/` → Business logic and AWS integration.
  - `main.go` → Backend entry point.
  - `go.mod` → Go dependencies.

- **frontend/** → **Vue.js** frontend code, allowing users to search and download sheet music.
  - `src/components/` → Reusable components.
  - `src/views/` → Main views (Home, Sheet Music, etc.).
  - `src/router/` → Vue Router configuration.
  - `src/store/` → State management (Vuex or Pinia).
  - `main.js` → Vue entry point.
  - `package.json` → Frontend dependencies.

- **infra/** → Configuration files to define and deploy infrastructure on **AWS** using **CloudFormation**.
  - `templates/` → AWS CloudFormation templates.
  - `deploy.sh` → Deployment script.

- **.github/** → **GitHub Actions** setup for **continuous integration and deployment (CI/CD)**.
  - `workflows/` → CI/CD pipelines.

- **docs/** → Project documentation.

- **README.md** → General project explanation.
- **LICENSE** → **GPLv3 License**, ensuring the code remains open-source.
- **.gitignore** → Files to be ignored by Git.

## Tech Stack
- **Backend**: Go, AWS Lambda, DynamoDB, API Gateway.
- **Frontend**: Vue.js, Vite.
- **Infrastructure**: AWS CloudFormation.
