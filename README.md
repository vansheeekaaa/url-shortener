# SnipURL

SnipURL is a simple and efficient URL shortening service allowing users to convert long URLs into short and easy-to-share links with optional expiration times. 

## Live deployment

- **Frontend:** Deployed on **Netlify**
- **Backend API:** Hosted on **Render** (kept awake via **UptimeRobot**)
- **Database:** Serverless Postgres via **Neon**

## Features

- **Fast API**: Built with Go and Gin framework for minimal overhead and fast routing.
- **Cryptographically Secure**: Utilizes Go's `crypto/rand` to generate cryptographically unpredictable short codes (protecting against predictable `math/rand` sequences).
- **Clean Architecture**: Strictly isolates functionality natively across Repositories, Services, and Handlers resulting in highly testable code. 
- **Error Handling & Retries**: Catches PostgreSQL unique constraint collisions (error `23505`) and handles automated retries for code generation without breaking the application flow.
- **Link Expiration**: Option to set explicit expiry durations (e.g., 1 Hour, 1 Day, 30 Days) for URLs, automatically returning `410 Gone / Expired` responses once the time passes.
- **URL Normalization:** Handles variations of the same URL (e.g., `google.com` vs `google.com/`) to avoid creating unnecessary new links.

## Tech Stack & Infrastructure

**Frontend**
- HTML5 & CSS3 
- Vanilla JavaScript (Async API Fetching, Clipboard API)
- **Deployment:** Netlify

**Backend**
- [Go](https://go.dev/) (1.21+)
- [Gin Web Framework](https://github.com/gin-gonic/gin)  (Fast routing & CORS Management)
- **Deployment:** Render (Pinged via UptimeRobot to prevent cold starts)

**Database**
- [PostgreSQL](https://www.postgresql.org/)
- `database/sql` + `lib/pq` Driver (native SQL interfacing)
- **Hosting:** Neon Serverless Postgres

## API Documentation

The backend has the following endpoints.

### `POST /shorten`
Creates a new short URL.

**Request Body:**
```json
{
    "url": "https://your-very-long-and-unwieldy-url.com/with/extra/path",
    "expiry_seconds": 3600
}
```
- `url` (string, required): The original URL to shorten.
- `expiry_seconds` (number, optional): The number of seconds until the link expires. If `0` or omitted, the link will never expire.

**Success Response (200 OK):**
```json
{
    "short_url": "http://localhost:8080/jK8pLq"
}
```

**Error Responses:**
- `400 Bad Request`: If the request body is invalid, the URL is malformed, or `expiry_seconds` is negative.
- `500 Internal Server Error`: If the server fails to generate a unique code after multiple attempts or experiences other database issues.

### `GET /:code`
Redirects to the original URL associated with a short code.

**Example:**
`GET /jK8pLq`

**Success Response:**
- `307 Temporary Redirect` to the original URL.

**Error Responses:**
- `404 Not Found`: If the short code does not exist.
- `410 Gone`: If the link associated with the code has expired.

### `GET /ping`
A simple health check endpoint. This is also utilized by UptimeRobot to ensure the API remains active.

**Success Response (200 OK):**
```json
{
    "message": "pong",
    "status": "alive"
}
```

## Running Locally

### 1. Database Setup
Ensure PostgreSQL is running locally or you have access to a cloud Postgres instance (like Neon). Execute the predefined `schema.sql` to build the required table and fast secondary indexes.
```bash
# Reference /backend/db/schema.sql for the table architecture
psql -U your_postgres_user -d your_db_name -f backend/db/schema.sql
```

### 2. Configure Environment
Inside the `/backend` folder, create a `.env` file containing your Postgres connection string and base URL:
```env
DATABASE_URL="postgres://your_user:your_password@localhost:5432/your_db?sslmode=disable"
BASE_URL="http://localhost:8080"
```

### 3. Start the Backend API
```bash
cd backend
go mod tidy
go run main.go
```
*The API will begin listening for interactions at `http://localhost:8080`.*

### 4. Launch the Frontend
1. Open `frontend/app.js` and set the `API_URL` constant to your local backend:
   ```javascript
   const API_URL = "http://localhost:8080";
   ```
2. Because the frontend is entirely Vanilla HTML and JavaScript without the need for a sluggish npm build step, you can simply **double-click** `/frontend/index.html` from your computer folder to open it instantly in any web browser!

## Testing

The core generation and abstraction logic within the `URLService` is comprehensively covered by robust component-level unit tests. Ensuring high reliability across edge cases like URL normalization, active link expiry checks, and empty inputs.

To execute the automated test suite:
```bash
cd backend/services
go test -v
```
