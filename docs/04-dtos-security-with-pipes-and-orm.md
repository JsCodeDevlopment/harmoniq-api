# 04 - DTOs, Entities and Secure Validations

A world-class foundation ensures that Garbage cannot enter the database (Garbage In / Garbage Out).
Gost uses `Pipes` and `DTOs` imported from the NestJS-oriented world.

---

## 1. Database Entities (`entities/`)

Location: `/src/modules/users/entities/user.entity.go`

The **Entity** is the physical blueprint of your data in the SQL base.

```go
type User struct {
	gorm.Model           // This adds ID, CreatedAt, UpdatedAt and DeletedAt fields
	Name   string `json:"name"`
	Email  string `json:"email" gorm:"unique"` // Creates the rule that it doesn't repeat in base (Constraint)
	Avatar string `json:"avatar"`
}
```

**GORM AutoMigrate:**
In your `users.module.go`, the _AutoMigrate()_ line reads this file to rebuild this table without the need for raw SQL migrations. The _gorm.Model_ field makes the entity Soft Deleteable (`DeletedAt`); when running Delete in the Service, the row is not removed from the base, it just acquires a deleted watermark via timestamp.

---

## 2. Restricted Inputs (DTO - Data Transfer Object)

Location: `/src/modules/users/dto/create-user.dto.go`

We isolate what enters from what is physical because **not everything from the JSON received in the POST should be saved as an entity** (For example, an _I Agree to Terms_ checkbox will not be in my Entity but the DTO will be tested).

**Rules (Struct Tags):**
In Go, we can force the Engine to validate inputs on DTOs through Tags:

```go
type CreateUserDto struct {
	Name  string `json:"name" binding:"required,min=3"`
	Email string `json:"email" binding:"required,email"`
}
```

This requires not only that an input JSON payload arrives, but checks that _Name_ has 3 or more letters, and fails the regex if _Email_ is not a syntactic construct of an alias (`blah@blah.com`).

---

## 3. The Garbage Interceptor: PIPES

Location: `/src/common/pipes/validation.pipe.go`

In NestJS we'd have `ValidationPipe`. In Gost we adopt a Pipe with **Go Generics (1.18+)**.

```go
func ValidateBody[T any](c *gin.Context) (*T, error) { ... }
```

Here is the ultimate architecture: We pass via base class `T` that it should compose for `dto.CreateUserDto`.
The pipe reads an array of loose bytes (`JSON`) sent by the frontend, checks field by field via the "binding" engine the validations (`min=3`, `email`), couples in C Memory to the struct and checks.

**Didactic Use:**
In your controller:

```go
d, err := pipes.ValidateBody[dto.CreateUserDto](c)
if err != nil {
	return // The DTO failed. The Pipe *ALREADY ANSWERED* 400 Bad Request to the Client (Angular)! We finish.
}
// From this line onwards, `d` (from payload) is 100% pure and there is no more chance of Null Pointers!
```

Thanks to the generics `[T any]` it will accept any file in your DTO sub-module without us needing to reinvent unit validators. Fantastic for ensuring robustness!

---

*PS: To explore the framework's pool of validators just see the references of Tags via Gin (*go-playground/validator*). Supports mathematical, logical comparisons of fields such as IPs, links, CPF etc.*
