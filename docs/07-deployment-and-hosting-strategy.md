# 07 - Deployment and Hosting Strategy

Deploying a Go application like Gost is straightforward because Go compiles down to a single, statically linked binary. This means you do not need a runtime environment (like Node.js or Python) installed on the target server.

---

## 1. Preparing for Production

Before deploying, ensure you compile the application. Go's cross-compilation is powerful.
If you are deploying to a standard Linux server (Ubuntu/Debian) via GitHub Actions or your local machine, run:

```bash
GOOS=linux GOARCH=amd64 go build -o gost-api main.go
```

This produces a file named `gost-api` which is all you need to execute on the server along with your `.env` file!

---

## 2. Hosting Recommendations

Because Gost is incredibly lightweight (usually consuming less than 30MB of RAM at idle), you have several excellent and cost-effective hosting options.

### Option A: Railway.app / Render (PaaS) - Recommended for Beginners/Startups

Platforms as a Service (PaaS) are the easiest way to deploy. They automatically detect Go environments.
**Why choosing them?**

- Zero server config.
- You push to GitHub, they build and deploy.
- Easy to attach managed PostgreSQL and Redis addons with 1-click.

**How to deploy (Render/Railway):**

1. Connect your Github Repository.
2. In the "Build Command", run: `go build -o app main.go`
3. In the "Start Command", run: `./app`
4. Copy your `.env` variables into the platform's "Environment Variables" tab.

### Option B: DigitalOcean Droplets / AWS EC2 (IaaS)

If you require maximum control and lower costs at scale, spinning up a Linux VPS is ideal.

**How to deploy via Docker (Recommended for VPS):**
Since we already have a `docker-compose.yml`, deploying is extremely easy. By adding a `Dockerfile` for the Go App, you can orchestrate everything.

1. **Create a `Dockerfile`** in your root:

```dockerfile
# Build Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# Run Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3000
CMD ["./main"]
```

2. **Update your `docker-compose.yml`** on the server to include the app:

```yaml
services:
  api:
    build: .
    ports:
      - "3000:3000"
    env_file: .env
    depends_on:
      - postgres
      - redis
  # ... (postgres and redis remain the same)
```

3. Run `docker-compose up -d --build` on your Droplet!

### Option C: Vercel (Serverless Functions)

**Overview:**
Vercel does not run traditional long-running Go servers. Instead, it runs serverless functions that execute when an HTTP request is received. While standard Vercel tutorials dictate creating multiple files inside the `/api` directory (where each file is an endpoint), doing so would break Gost's robust monolithic architecture and Dependency Injection.

To deploy Gost (which uses the Gin framework) on Vercel, we map **all** incoming traffic through a single Serverless "catch-all" handler.

**Step 1 — Configure Vercel Routing (`vercel.json`)**
Create a `vercel.json` file in your root folder:

```json
{
  "functions": {
    "api/index.go": {
      "runtime": "vercel-go@3.0.0"
    }
  },
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/api/index.go"
    }
  ]
}
```

_This catches all traffic and routes it to our single Go handler._

**Step 2 — Exposing the Router in Gost**
You must adjust `src/app/app.module.go`. Currently, `Bootstrap()` runs `router.Run()`. You need to extract the router setup so Vercel can consume it.

Modify your `app.module.go` (extracting the `router.Run` part) to return the `*gin.Engine`:

```go
package app
// ... imports

func Bootstrap() *gin.Engine {
    config.LoadEnv()
    // ConnectDatabase() -> (Warning: ensure your DB supports serverless connections/pooling)
    // ConnectRedis()

    router := gin.Default()
    router.Use(config.SetupCors())

    // ... middlewares

    api := router.Group("/api/v1")
    users.InitModule(api)

    return router
}
```

**Step 3 — Create the Serverless Entrypoint**
Create a new folder and file at the root of the project: `api/index.go`.

```go
package handler

import (
    "net/http"
    "gost/src/app"
)

// We define the GIN engine globally so it boots only once per warm serverless container
var engine = app.Bootstrap()

func Handler(w http.ResponseWriter, r *http.Request) {
    // We let Gin HTTP engine handle the raw Vercel request transparently
    engine.ServeHTTP(w, r)
}
```

**Step 4 — Test and Deploy**

- Install Vercel CLI: `npm install -g vercel`
- Test locally: `vercel dev`
- Deploy: `vercel login` and `vercel`

**⚠️ Critical Limitations for Gost on Vercel:**

- **Database Connections:** Serverless functions scale to 1000s of instances instantly, booting from zero. This will exhaust standard PostgreSQL connections quickly. You **MUST** use a connection pooler like PgBouncer or a Serverless DB (like Supabase, Neon.tech, or PlanetScale).
- **In-Memory states:** Go routines or global maps will not persist between requests as Vercel kills the container after completion. Redis is mandatory for caching.
- File uploads (`/uploads`) to the local disk will fail on Vercel because the file system is read-only and ephemeral. You must change your upload utility to use S3 (AWS) or standard Cloud Storages.

---

## 3. Reverse Proxies (Nginx / Caddy)

If you deploy on a VPS (Option B), never expose port 3000 directly. Use a reverse proxy to handle SSL/TLS (HTTPS).

**Example with Caddy (Extremely Easy Mux):**
Install Caddy and map your domain to the internal Gost port:

```text
api.yourdomain.com {
    reverse_proxy localhost:3000
}
```

Caddy will automatically provision and renew Let's Encrypt SSL certificates for you!

---

## 4. Environment Checklist for Prod

- Ensure `DB_HOST` and `REDIS_HOST` point to the correct production credentials.
- Ensure `ALLOWED_CORS` has your exact Frontend domains listed, removing `*`.
- Set Gin to Release Mode by adding `GIN_MODE=release` to your `.env` file to prevent Gin from printing debug logs into standard output in production.
