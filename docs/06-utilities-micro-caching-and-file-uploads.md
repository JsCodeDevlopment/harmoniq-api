# 06 - Utilities, Upload and Local Cache

On a daily basis we need more than just trafficking JSON. Sometimes the flow contains external media (Multipart Binaries) or metrics stuck in bottlenecks, where we introduce utilities to ease the journeys of the DEV and the User.

---

## 1. The Multer Substitute: File Upload Utility

Location: `/src/common/utils/file-upload.util.go`

In Node, we suffered to learn to implement "Multipart" interception. The Gost lib introduces into its tool drawers (`utils/`) the simple magic of `UploadImage`, which imitates the native behavior of downloading Binary Payloads.

When hitting the endpoints head-on (like our avatar example: `POST /users/:id/avatar`), an Android or Front-end client fires its binary form on network.
In the controller we call the engine:

```go
func UploadImage(c *gin.Context, fieldName string, destFolder string)
```

**How it works Didactically:**

1. The tool retrieves the buffer with Gin `c.FormFile("fieldName")`. (In frontend the file input must be called the same field).
2. If the local backend (in container) does not have the referred folder (e.g., `./uploads`), it builds the folder dynamically (`os.MkdirAll`) using global permissions support (`os.ModePerm`).
3. Instead of using playful names, it creates cryptographic cache overwrite protection converting the original name to a UNIX millisecond and returning (e.g., `170293023901.png`).
4. It allocates on Secondary Memory (Disk, I/O) and returns to database in your Controller DTO (`filePath, err`). The Service now saves this absolute link with ease! You have all the data processed in a clean native function!

---

## 2. Abstracting Micro-Caching via REDIS Database

Location: `/src/config/redis.go` and Global instantiations.

In heavy APIs, direct listing requests (like `GET /users/super-query`) violently hurt the Database and generate instance costs! What if you want to do caching like in nest providers with `@Cache()` interceptors?

We injected from the native Core ecosystem and activated with the parallel _containers_ supported through a built-in Golang SDK: `github.com/go-redis/redis/v8`.

**Configured on Boot, ready to use!**
Whenever your system needs to store ready results without triggering Repositories (_Read/Write Intensive Flows_), use from any file in your architecture (In _Services_):

```go
import "gost/src/config"
import "time"

// Saving the Payload that was already serialized from JSON to the Temporary Redis Db, expiring in hours
config.RedisClient.Set(config.Ctx, "my_query_v1", listResults, time.Hour)

// Reloading (Fetching to return in view even before hitting the SQL API):
val, err := config.RedisClient.Get(config.Ctx, "my_query_v1").Result()
```

_Power in hands without headache of configuring clients in pool every time, already modular and singleton_.

The End! With these utilities in hand + Severe Validation (Pipes) + Acentric Injection (Modules), **Gost** consolidates itself as a fantastic architectural casing for teams migrating or implementing robustness without the "overhead" that would be expected in Go! Happy Coding!
