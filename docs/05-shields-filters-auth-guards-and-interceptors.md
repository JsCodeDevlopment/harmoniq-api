# 05 - Security Middlewares and Advanced Interceptors

The real web request flow is almost never static from point "A" to "B" in the controller. We require logs, authentication, metrics. This is shaped by _Guards_ and _Interceptors_.

---

## 1. Guard Traffic: The Guards (Route Protection)

Location: `/src/common/guards/auth.guard.go`

If you need a restriction by Session or **JWT** Token, this is the interception place, so that restrictions can be massively validated without dirtying the Controller's logic.

**Guard Example (`AuthGuard`)**:
The function acts at the exact moment the Network calls your route, reading from the Request:

```go
token := c.GetHeader("Authorization") // Checks presence of `Bearer EyJ....`
if token == "" {
    c.AbortWithStatusJSON(http.StatusUnauthorized, ... ) // Blocked, Aborts subsequent Controller
    return
}
```

**How to lock Routes?**
Go back to the initialization of your module in the `users.module.go` file or on singular routes:

```go
// Example locking a specific route of the module:
usersGroup.POST("/super-secret", guards.AuthGuard(), controller.CreateSecret)
```

In this guard class you implement your local parser JWT engine (`jwt.Parse(...)`). Validate, inject `c.Set("user_id", payload)` for the controller to read, or block with 403 error `Abort()`.

---

## 2. Visual Logistics: Interceptors (Logger)

Location: `/src/common/interceptors/logger.interceptor.go`

Inspired by Exception Filters and Request Logging in Nest, the Interceptor wraps your controller "Before and After" network processing!

The processing flow works due to Gin's `c.Next()`.

```go
// BEFORE (Calculation Starts)
start := time.Now()
log.Printf("incoming...") // Prints who hit the api

// MIDDLE
c.Next() // Returns scope for your app to calculate (Controllers and DB)

// AFTER
duration := time.Since(start)
log.Printf("outgoing...") // At the end of the process C returns from NEXT stack, finishing the metric
```

You can create middlewares here with _tracing telemetry_ for example, using OpenTelemetry! Totally globally coupled with extreme lightness!

---

## 3. The Final Shield: Exception Filter

Location: `/src/common/filters/http-exception.filter.go`

Without a doubt one of the greatest clean structuring powers of Nest and now of the Gost lib.
In simple APIs people handle JSONs with giant chaotic code blocks containing `if err != nil...`.
In Gost we use Delegation:

If an SQL Database Exception occurs when trying to list in your controller: `c.Error(err)`. Just that.

The **global `ErrorHandler`** runs at the processing tail. It watches the request to see if anyone used the `c.Errors` bucket:

```go
if len(c.Errors) > 0 {
    err := c.Errors.Last()
    // Prevents sensitive data from leaking, wrapping in 500 Internal Server Error with Frontend Platform Standards!
    c.JSON(...)
}
```

This will give you peace of mind, ensuring that under no collapse an ugly "Go Dump Trace" message appears on your end users' screens, always returning a concise, validated and secure `statusCode, error, message` key contract!
