# Gost - Official Documentation

## 1. Introduction and Philosophy

Welcome to the **Gost** documentation! This library (or boilerplate) was designed to solve the main pain point of Go developers coming from the Node.js ecosystem: **the lack of standards in complex projects.**

While in Go we have the freedom to structure things any way we want, NestJS and Angular taught us that strict standards, separation of concerns, and dependency injection save scaling projects and drastically reduce a team's maintenance costs.

Gost is built on top of the **Gin Web Framework** for extreme C10K performance, but orchestrated identically to NestJS.

---

## 2. Gost Architectural Patterns

Gost doesn't reinvent the wheel. It adopts the fundamental concepts of Nest and maps them as follows:

1. **Modules (`.module.go`)**
   In a clean project, a domain (e.g., _Users_) shouldn't know how _Authentication_ is done. In Gost, you create "black boxes" called Modules. A Module is responsible for orchestrating its own parts (Controllers, Services, and Repositories) and attaching them to the network (Gin Router).
2. **Controllers (`.controller.go`)**
   The Controller in Gost is dumb. Its only function is to **receive HTTP messages (JSON, FormData, Queries)**, validate the syntax of what clients sent using _Pipes_, delegate execution to the Service, and return a `c.JSON()` at the end. No interacting with the SQL Database here!

3. **Services (`.service.go`)**
   They act as the providers with `@Injectable()` in Nest. The Service is where the business rules live: "_If the user clocks in late, deduct from their balance_". The service doesn't know if it was called by a REST API, a CLI bot, or Kafka/RabbitMQ messaging.

4. **Repositories (`.repository.go`)**
   They act as intermediaries between your business rules and the Oracle (The Database). The Repository is restricted to speaking in GORM or pure Queries (SQL). Any data access goes through it, abstracting access so that if tomorrow you stop using Postgres and switch to MongoDB, **all your changes will only happen in the Repository.**

5. **Entities & DTOs**
   They are the Gost contracts.
   - **Entities (`entities/`)**: How is my data mapped in the database?
   - **DTOs (`dto/`)**: What does my API expect to receive from the Front-End (via Pipes)?

---

## 3. The Lifecycle of a Request

When Angular requests `POST /api/v1/users`:

1. The call passes through the `Cors Middleware` (checks HTTP Security).
2. Passes through the `LoggerInterceptor` (logs what time the request started).
3. Enters the `UsersController`.
4. The Pipe validates its `CreateUserDto`. If a valid email is missing, it explodes into a 400 Bad Request that is caught and masked by the `ErrorHandler`. If OK, it proceeds.
5. The Controller sends the clean DTO to the `UsersService`.
6. The `UsersService` thinks ("Is the email unique?"), passing it to the `UsersRepository`.
7. The `UsersRepository` saves the data via GORM in PostgreSQL, returning the saved object.
8. The `UsersController` sends the response to the Client.
9. The `LoggerInterceptor` finishes, printing on the console that the request took 15ms.

Continue reading the next files to dive deep `file by file` into the system!
