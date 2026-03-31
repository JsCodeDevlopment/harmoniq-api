# 03 - Domain-Driven Modules and Dependency Injection

Every time a business scope Entity emerges (e.g., _Orders, Invoices, Customers_), it will gain what we call a "Domain Module" within Gost.

Let's detail the anatomy of a module taking the `Users` Example folder (`src/modules/users`).

---

## 1. The Entry Point: `users.module.go`

The module itself is essentially a class or initialization function. All **Fixed Dependency Injection** occurs in it.

```go
func InitModule(router *gin.RouterGroup) {
	// The "Wiring" Magic Happens Here:
	repository := NewUserRepository(config.DB)
	service := NewUserService(repository)
	controller := NewUserController(service)
	// ...
}
```

NestJS does this using reflection (`@Injectable()`). Gost achieves **extreme performance painlessly** by doing it proactively and typed in go.

- The Repository was born and received the Database.
- The Service received the Repository, so the Service doesn't know "_GORM_", it knows the "Contract that talks to the base".
- The Controller received the Service and hooked itself to the Group Routes (`router.Group("/users")`).

---

## 2. The Presentation Layer: `users.controller.go`

The Controller package acts as a "bouncer". It listens for requests on the Internet lane and takes Go packages, returning packages via network.

**Controller Structure:**
The struct contains `service UserService`.

### How to use in practice?

Always receive a network context via `(c *gin.Context)`.

1. **Retrieving URL Variables (`:id`):**
   Use `c.Param("id")`.
2. **Formatting Strict HTTP Responses:**
   Use `c.JSON(http.StatusOK, returnData)`. Never pass the magic string `200`, it foresees unseen bugs! Always use the internal HTTP constants library from the Go bank `net/http`!
3. **Escalating Errors:**
   On a request error (Like validation failure or non-existent user), stop the function right there!
   `c.JSON()` or use the error escalation to the Global Filter injecting and quitting:
   ```go
   c.Error(err) // Sends the exception to the global ErrorHandler!
   return       // And abort control of this memory stack pointer
   ```

---

## 3. The Brain: `users.service.go`

To fulfill the O in SOLID principles "Open/Close" all your Services are implemented from `Interfaces`.

The `UserService` Interface dictates the contract. Right below you have the `type userService struct` with the real layer.
Why this? **For Unit Testing**. If your controller explicitly requires the real struct, it couldn't be tested without a database because it would be tightly coupled. With an Interface, in the Front-end or CI/CD Unit Testing folders and routines you would pass a static class that fills fake queries taking zero I/O processing!

**Application Tip**: Your vital validations should stay here, before the repository intervenes.
_Is the invoice amount negative? Handle and emit the error! "HTTP rules should not enter this class", meaning... you wouldn't do JSON repassing here, only In-Memory Classes of your code (Entities, DTOs)!_

---

## 4. The Hidden Worker: `users.repository.go`

The repository exposed by the `UserRepository` interface manages transactions. It has the Pointer to the actual GORM "_DB Object_" (`db *gorm.DB`).

For simple base operations, it's recommended to invoke direct GORM functions. For difficult cases or searches of multiple complex relationships, you abstract to Go like the example:

```go
func (r *userRepository) FindById(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error // Returns only the selected Primary ID
	return &user, err
}
```

This way, your controller/Service rely purely on this bridge between the abstract machine and the physical Postgres Database.
