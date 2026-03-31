The **Gost CLI** is the framework's automation engine. It is designed to eliminate boilerplate code and speed up the development process by scaffolding projects and generating full-stack modules following Gost's architectural standards.

Starting from version 1.1.0, the CLI is **Standalone**, meaning it carries the entire framework within its binary using Go's embedding feature.

---

## 📂 Standalone Architecture

The CLI follows a command-pattern structure using the [Cobra](https://github.com/spf13/cobra) library.

### 1. Embedded Templates
The framework source code (`src/`, `locales/`, `docs/`, etc.) is embedded into the binary using `go:embed TemplateFS`. This allows the CLI to initialize new projects without needing to clone the repository or have the source code locally.

### 2. Project Initializer (`commands/init.go`)
The `gost init` command leverages the embedded filesystem:
- **Extraction**: It walks through `TemplateFS` and writes the files to the user's destination directory.
- **Module Pruning**: If the user chooses a "Basic" template, it physically removes the code and configurations for modules not selected (e.g., Auth, RabbitMQ, i18n).
- **Template Patching**: Replaces the generic `gost` module name with the new project name across all `.go` and `go.mod` files.

---

## 🚀 Installation & Distribution

The Gost CLI can be installed globally without cloning the repository.

### 1. One-liner (Fastest)
Ideal for Linux, macOS, and Git Bash:
```bash
curl -sSL https://gost.run/install.sh | sh
```

### 2. Go Global (Recommended for Developers)
Install directly into your `$GOPATH/bin`:
```bash
go install github.com/JsCodeDevlopment/gost/cmd/gost@latest
```

### 3. NPX (Node.js ecosystem)
Run without permanent installation:
```bash
npx gost-cli init my-project
```

---

## 🛠️ Usage Guide

### Initializing a Project
```bash
gost init my-api
```
Follow the interactive prompts to name your project and select the modules you need (Authentication, Messaging, i18n).

### Generating a Domain Module
```bash
gost make:module catalog
```

### Full Stack CRUD Generation
```bash
gost make:crud order
```
*Creates: Entity, DTOs, Repository, Service, Controller and registers them in `app.module.go`.*

---

## 🔄 Internal Logic Flows

### Project Scaffolding Flow
1. **Extraction**: Recursive walk of `gost.TemplateFS`.
2. **Prune**: If "Basic", delete unwanted folders from the newly created directory.
3. **Patch**: 
   - Update `go.mod` module name.
   - Update all imports in `.go` files.
   - Remove init calls in `app.module.go` for pruned modules.

### CRUD Generation Flow
1. **Templates**: Injected Go-string templates with placeholders.
2. **Project Detection**: Reads the local `go.mod` to ensure imports match your application name.
3. **File Creation**: Writes 6 distinct files.
4. **Registration**: 
   - Finds `import (` and injects the new module path.
   - Finds `ws.InitModule(api)` and injects the new `<module>.InitModule(api)` below it.

---

## ⚠️ Important Considerations

- **Standalone Mode**: The CLI carries a snapshot of the framework from when it was built. To get the latest framework updates, simply run `go install ...@latest` again.
- **Server Restart**: After running `make:crud`, you must restart your Go server for the new routes to be registered.
- **Database Migration**: The CLI generates the Entity, but you should ensure your DB is configured to handle the new table.
