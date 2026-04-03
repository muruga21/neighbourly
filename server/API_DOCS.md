# Neighbourly Server API Documentation

This document maintains a list of all available REST API endpoints for the Neighbourly backend, including their required parameters and expected responses.

---

## 1. System Endpoints

### 1.1 Root
- **Endpoint**: `GET /`
- **Description**: A basic index route.
- **Request Parameters**: None
- **Response**:
  - `200 OK`
  ```text
  Hello, World!
  ```
  - `404 Not Found` (If accessing any unhandled path that falls through to root)

### 1.2 Health Check
- **Endpoint**: `GET /health`
- **Description**: Used to check if the server is up and running.
- **Request Parameters**: None
- **Response**:
  - `200 OK`
  ```text
  im alive
  ```

---

## 2. Authentication

### 2.1 User Signup
- **Endpoint**: `POST /api/auth/signup`
- **Description**: Registers a new user into the system using standard email/password authentication and returns a JWT token.
- **Headers**:
  - `Content-Type: application/json`
- **Request Body Parameters**:
  | Field | Type | Description | Required |
  |-------|------|-------------|----------|
  | `fullName` | string | The user's full name | Yes |
  | `email` | string | The user's email address | Yes |
  | `phoneNumber` | string | The user's phone number | No |
  | `password` | string | The user's chosen password (min 6 chars recommended) | Yes |
  | `role` | string | The user's role (e.g., "seeker" or "provider") | No |

- **Example Request JSON**:
  ```json
  {
      "fullName": "John Doe",
      "email": "john.doe@example.com",
      "phoneNumber": "123-456-7890",
      "password": "mySecurePassword123",
      "role": "seeker"
  }
  ```

- **Responses**:
  - **Success (`201 Created`)**
    ```json
    {
        "success": true,
        "message": "Account created successfully.",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
    }
    ```
  - **Error (`400 Bad Request`)** - Triggered by missing required fields or invalid JSON.
    ```json
    {
        "success": false,
        "message": "Missing required fields",
        "token": ""
    }
    ```
  - **Error (`409 Conflict`)** - Triggered if a user with that email already exists in the database.
    ```json
    {
        "success": false,
        "message": "Email already exists",
        "token": ""
    }
    ```
  - **Error (`500 Internal Server Error`)** - Triggered by hashing failures or database connection/insertion issues.
    ```json
    {
        "success": false,
        "message": "Error creating user",
        "token": ""
    }
    ```

### 2.2 User Login
- **Endpoint**: `POST /api/auth/login`
- **Description**: Authenticates an existing user using their phone number and password, returning a JWT token.
- **Headers**:
  - `Content-Type: application/json`
- **Request Body Parameters**:
  | Field | Type | Description | Required |
  |-------|------|-------------|----------|
  | `phoneNumber` | string | The user's phone number | Yes |
  | `password` | string | The user's password | Yes |
  | `role` | string | The user's expected role (e.g., "seeker") | No |

- **Example Request JSON**:
  ```json
  {
      "phoneNumber": "+15550000000",
      "password": "userPassword123!",
      "role": "seeker"
  }
  ```

- **Responses**:
  - **Success (`200 OK`)**
    ```json
    {
        "success": true,
        "message": "Login successful",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
    }
    ```
  - **Error (`400 Bad Request`)** - Triggered by missing fields or invalid JSON payload.
    ```json
    {
        "success": false,
        "message": "Missing required fields",
        "token": ""
    }
    ```
  - **Error (`401 Unauthorized`)** - Triggered by an incorrect phone number or password match.
    ```json
    {
        "success": false,
        "message": "Invalid phone number or password",
        "token": ""
    }
    ```
