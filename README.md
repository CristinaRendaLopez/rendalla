# Rendalla
Web platform for managing and distributing sheet music.

## Features
- Upload and store sheet music in different formats.
- Search and download music sheets for various instruments.
- Authentication. Only authenticated users can create, update, or delete resources. Public endpoints are read-only.

## Setup
These instructions will get the project running locally for development and testing purposes.
### Prerequisites
- [Go](https://go.dev/dl/) 1.20+
- [Node.js](https://nodejs.org/) 18+
- [npm](https://www.npmjs.com/)
- DynamoDB Local for local development and tests.
You can run DynamoDB Local using Docker:

```bash
docker run -d -p 8000:8000 amazon/dynamodb-local
```

### Backend setup
```bash
# Navigate to backend
cd backend

# Run the application
go run app/main.go
```
Make sure you have a .env file in the backend/ directory following this example:
```env
ENV=test

# AWS
AUTH_SECRET_NAME=your_auth_secret
AWS_REGION=your-aws-region
AWS_ACCESS_KEY_ID=your_key_id
AWS_SECRET_ACCESS_KEY=your_secret_key

# DynamoDB
SONGS_TABLE=your_songs_table
DOCUMENTS_TABLE=your_documents_table

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION_HOURS=1

AUTH_USERNAME=test
AUTH_PASSWORD=hashed_password_here

# Server
APP_PORT=8080
```
### Frontend setup
```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## Running tests
**Backend:**

```bash
cd backend
go test ./...
```
**Frontend:**

```
cd frontend
npm run test
```


## Project Structure

- **backend/** → Contains the **Go** backend code, which runs on **AWS Lambda** and exposes a REST API for managing songs and sheet music.
  - `app/` → Entry point for building the Lambda handler.
  - `bootstrap/` → Loads environment configuration and initializes the DynamoDB client for the application.
  - `dto/` → Data Transfer Objects for request/response handling.
  - `errors/` → Centralized application error types.
  - `handlers/` → HTTP request handlers (songs, documents, search, authentication).
  - `integration_tests/` → Integration tests for validating complete API behavior.
  - `middleware/` → Custom middleware (JWT authentication and validation).
  - `mocks/` → Mock implementations of services and repositories used in unit testing.
  - `models/` → Internal data structures representing songs and documents.
  - `repository/` → Database access layer using DynamoDB.
  - `router/` → Gin router configuration and route registration.
  - `services/` → Business logic for song/document management, search, and authentication.
  - `utils/` → Utility functions (e.g., UUIDs, time, normalization).
  - `go.mod` / `go.sum` → Go dependency declarations.

- **frontend/** → **Vue.js** frontend code.
  - `src/components/` → Reusable UI components.
  - `src/views/` → Main views (Home, Sheet Music, etc.).
  - `main.js` → Vue entry point.
  - `package.json` → Frontend dependencies.

- **.github/** → **GitHub Actions** setup for CI/CD.
  - `workflows/` → CI/CD pipelines.

- **README.md** → General project explanation.
- **LICENSE** → GPLv3 License.
- **.gitignore** → Files to be ignored by Git.

## Tech Stack
- **Backend**: Go, AWS Lambda, DynamoDB, API Gateway.
- **Frontend**: Vue.js.
- **Infrastructure**: AWS CloudFormation.
