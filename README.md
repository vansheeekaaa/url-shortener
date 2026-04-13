# SnipURL

SnipURL is a fast URL shortening service that converts long URLs into short, shareable links with optional expiration times and real-time click analytics.

## Live Deployment

- **Frontend:** Deployed on **Netlify**
- **Backend API:** Hosted on **Render** (kept awake via **UptimeRobot**)
- **Database:** Serverless Postgres via **Neon**

## Features

- **Fast API**: Built with Go and Gin framework for minimal overhead and fast routing.
- **Cryptographically Secure**: Utilizes Go's `crypto/rand` to generate cryptographically unpredictable short codes (protecting against predictable `math/rand` sequences).
- **Clean Architecture**: Strictly isolates functionality across Repositories, Services, and Handlers, resulting in highly testable, decoupled code.
- **Error Handling & Retries**: Catches PostgreSQL unique constraint collisions (error `23505`) and handles automated retries for short code generation without breaking the application flow.
- **Link Expiration**: Option to set explicit expiry durations (e.g., 1 Hour, 1 Day, 30 Days). Expired links automatically return `410 Gone` responses.
- **URL Normalization**: Handles variations of the same URL (e.g., `google.com` vs `google.com/`) to avoid creating duplicate entries.
- **Click Analytics**: Tracks total click count and last accessed time per link. Click recording is fired asynchronously (goroutine) so it never adds latency to the redirect response.
- **Live Stats UI**: The frontend polls the stats endpoint every 10 seconds while the stats card is open, giving a real-time view of incoming clicks.

## Tech Stack & Infrastructure

**Frontend**
- HTML5 & CSS3
- Vanilla JavaScript (Async/Await, Clipboard API, Polling)
- **Deployment:** Netlify

**Backend**
- [Go](https://go.dev/) (1.21+)
- [Gin Web Framework](https://github.com/gin-gonic/gin) (Fast routing & CORS management)
- **Deployment:** Render (pinged via UptimeRobot to prevent cold starts)

**Database**
- [PostgreSQL](https://www.postgresql.org/)
- `database/sql` + `lib/pq` driver (native SQL, no ORM)
- **Hosting:** Neon Serverless Postgres

## API Documentation

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
- `expiry_seconds` (number, optional): Seconds until the link expires. `0` or omitted means the link never expires.

**Success Response (200 OK):**
```json
{
    "short_url": "http://localhost:8080/jK8pLq",
    "short_code": "jK8pLq"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body, malformed URL, or negative `expiry_seconds`.
- `500 Internal Server Error`: Failed to generate a unique short code after maximum retries.

---

### `GET /:code`
Redirects to the original URL. Increments the click counter asynchronously.

**Example:** `GET /jK8pLq`

**Success Response:**
- `307 Temporary Redirect` to the original URL.

**Error Responses:**
- `404 Not Found`: Short code does not exist.
- `410 Gone`: Link has expired.

---

### `GET /stats/:code`
Returns click analytics for a short URL.

**Example:** `GET /stats/jK8pLq`

**Success Response (200 OK):**
```json
{
    "short_code": "jK8pLq",
    "original_url": "https://your-very-long-url.com/",
    "click_count": 42,
    "created_at": "2026-04-13T18:00:00Z",
    "last_accessed_at": "2026-04-13T19:15:00Z",
    "expires_at": null
}
```

**Error Responses:**
- `404 Not Found`: Short code does not exist.

---

### `GET /ping`
Health check endpoint, also used by UptimeRobot to prevent cold starts.

**Success Response (200 OK):**
```json
{
    "message": "pong",
    "status": "alive"
}
```

## Running Locally

### 1. Database Setup
Ensure PostgreSQL is running locally or you have access to a cloud Postgres instance (like Neon). Run `schema.sql` to create the table and indexes:
```bash
psql -U your_postgres_user -d your_db_name -f backend/db/schema.sql
```

### 2. Configure Environment
Inside the `/backend` folder, create a `.env` file:
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
*The API will begin listening at `http://localhost:8080`.*

### 4. Launch the Frontend
1. Open `frontend/app.js` and confirm `API_URL` points to your local backend:
   ```javascript
   const API_URL = "http://localhost:8080";
   ```
2. Open `/frontend/index.html` directly in any browser, no build step required.

## Testing

The `URLService` is fully covered by unit tests using an in-memory mock repository. Tests cover URL normalization, expiry logic, invalid inputs, and click analytics.

```bash
cd backend/services
go test -v
```
