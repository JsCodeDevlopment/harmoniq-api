# 02 - Bootstrap and Configurations

Gost is self-contained and its startup is triggered centrally. Let's x-ray how the orchestration happens before the system is "Listening on port 3000".

---

## 1. `main.go` - The Ignition Point

Location: `/main.go`

This is the native and mandatory file in Go. Its sole role is to call the base `app` module and it shouldn't contain business rules or route definitions.

```go
package main

import "gost/src/app"

func main() {
	app.Bootstrap()
}
```

---

## 2. `src/app/app.module.go` - The Heart (`AppModule`)

Location: `/src/app/app.module.go`

Compared to Nest's `app.module.ts`, this is the root of your project. The `Bootstrap()` method:

1. **Pulls `.env`**: Through `config.LoadEnv()`.
2. **Connects peripherals**: Connects to GORM (`ConnectDatabase()`) and Redis (`ConnectRedis()`).
3. **Starts the Framework**: `gin.Default()` wraps the HTTP Engine.
4. **Configures CORS**: Adds rules declared in `.env`.
5. **Configures Middlewares (Interceptors and Filters)**:
   ```go
   router.Use(interceptors.LoggerInterceptor())
   router.Use(filters.ErrorHandler())
   ```
   **Didactic Tip**: By doing a `router.Use()`, **EVERY** request entering the api will pass through them. Since the Exception Filter is here at the top, it prevents your entire API from crashing ("Panic") in case of a _Data Race_ or Nil Pointer in a distant Module.
6. **Configures API Groups (Versioning)**:
   ```go
   api := router.Group("/api/v1")
   ```
   It ensures all modules below won't be at the root "/", but inside "/api/v1/". The V1/V2 practice avoids breaking Mobile apps if we update routes.
7. **Module Injection**: Here SubModules are instantiated, e.g.: `users.InitModule(api)`.
8. **And finally the Wheel of Life**: Executes the OS configured local loop (_port_): `router.Run(":" + port)`.

---

## 3. Environment and `.env` Files

Location: `/.env` / `src/config/env.go`

The utility function of the `config` package loads the file from the Disk and makes the `os.Getenv()` function usable for the entire API.
Always define the keys you'll use there in the root file (`.env.example`), as it helps other devs know what local configurations they need if they clone from Github!

---

## 4. `src/config/cors.go` - Cross-Domain Security

Location: `/src/config/cors.go`

_Cross-Origin Resource Sharing_ is responsible for not letting a malicious hacker panel steal cookies or make JavaScript requests to your API on behalf of other users.
The middleware generator function `SetupCors()`:

- Takes the `ALLOWED_CORS` env, splits it by comma (`,`), generating an array and allows only **Safe Front-End Origins** to use your API.

---

## 5. `src/config/database.go` and `redis.go` - Persistence

Location: `/src/config/database.go` and `/src/config/redis.go`

These files define the Global Variables (Singleton) of your Engines.
Go doesn't have full metaprogramming support like TypeORM or Prisma, but GORM comes very close.
`ConnectDatabase` formats a `DSN` (Data Source Name string) unifying your credentials from `.env` and opens an asynchronous connection Pool to Postgres kept in the `config.DB` variable.

**Architecture Tip**: Note that we will never use `config.DB` in a controller. The Repository will receive it injected!

Similarly, the `Context` (`config.Ctx`) and `RedisClient` are populated to support cache transactions in RAM, running smoothly without blocks thanks to local Go-Routines.
