# 08 - Testing Strategies (Unit & E2E)

Testing is a core component to ensuring a robust architecture. Drawing inspiration from NestJS conventions, Gost separates "Unit Tests" (testing isolated methods) from "End-To-End (E2E) Tests" (testing full module/API intersections).

We utilize Go's extremely fast native `testing` engine and pair it with `github.com/stretchr/testify` to gain access to assertion functionalities and dependency mocking reminiscent of `Jest`.

---

## 1. Unit Testing (Inside Modules)

Unit Tests are strictly kept within the layer they belong to. When you create a Module, you place a `tests/` folder inside of it so its unit validations never pollute global imports.

**Location:** e.g., `src/modules/users/tests/users.service_test.go`

Unit Testing heavily revolves around mock implementations to guarantee speed and pure logic testing:

- **Mocking Repositories:** Testing a `Service` relies on replacing its `Repository` interface parameter with a mock (like `testify/mock`).
- **Zero Disk Activity:** A Unit Test SHOULD NOT read logic from real Postgres DBs. If a test is talking to Postgres, it's an Integration Test, not a Unit Test.

**Example Structure:**

```go
mockRepo := new(MockUserRepository)
service := users.NewUserService(mockRepo)

// Forcing mock result:
mockRepo.On("Create", mock.Anything).Return(nil)

// Run
user, err := service.Create(fakeDto)
```

Running tests on a module:

```bash
go test ./src/modules/users/tests/... -v
```

---

## 2. End-To-End (E2E) Integration Testing (Root Level)

E2E testing proves the true network capability of the API. It assumes the roles of the client applications, firing raw JSON bodies globally to verify your entire pipeline: Request -> Pipe Validators -> Controllers -> Services -> Repositories -> Data Store.

**Location:** `/test/e2e/users_e2e_test.go`

To accomplish this effectively, the framework `app.module.go` has been refactored separating the server instantiation logic (`SetupApp()`) from the listener loop (`Bootstrap()`). E2E fetches `SetupApp()` directly out-of-the-box.

**E2E Lifecycle in Gost:**

1. Call `app.SetupApp()`. This builds `Gin`, wires database contexts, and triggers modules' registrations.
2. Initialize an `httptest.NewRecorder()` (a simulated buffer equivalent to what a navigator expects).
3. Push fake standard Go Http requests `http.NewRequest` containing bytes pointing directly at endpoints (e.g. `/api/v1/users`).
4. Read the buffered responses checking JSON structure keys and asserting status headers!

Running all Global tests:

```bash
go test ./test/... -v
```

---

## 3. Libraries

- `github.com/stretchr/testify/assert`: Simplifies traditional `if result != expected { t.Error() }` blocks with elegant variants like `assert.NoError(t, err)`.
- `github.com/stretchr/testify/mock`: Toolset that records method calls allowing you to inject behaviors without modifying your business code.

> **Tip on Test Environments:** While running E2E tests, the `config.LoadEnv()` triggers and the codebase uses your `.env` connection (pointing to your local Docker). In CI/CD pipelines (like Github Actions) simply declare separate step environments variables (Spin up a fresh dummy Postgres Database before testing pipeline triggers!).
