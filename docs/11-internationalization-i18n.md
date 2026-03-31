# 11 - Internationalization (i18n)

**Gost** includes a robust Internationalization (i18n) system that allows your application to support multiple languages seamlessly. It handles locale detection, message translation, and localized validation errors.

---

## 📂 Architecture and Files

The i18n system is located in `src/common/i18n` and consists of three main components:

### 1. Translation Provider (`provider.go`)
This is the core engine that uses `github.com/nicksnyder/go-i18n/v2`.
- **`Initialize(localesPath string, defaultLang language.Tag)`**: Loads all `.json` files from the specified directory and sets the fallback language.
- **`T(localizer *i18n_lib.Localizer, messageID string, templateData interface{})`**: The main function to translate a message. If the key is not found, it returns the key itself as a fallback.

### 2. Gin Middleware (`middleware.go`)
Automatically manages locale detection for every incoming request.
- Detects the preferred language from the `Accept-Language` header.
- Injects a `Localizer` instance into the Gin context.
- Provides `i18n.Translate(c, "key")` as a shorthand for controllers.

### 3. Validator Localization (`validator.go`)
Integrates with `go-playground/validator/v10` to provide human-readable validation errors in the user's language.
- Registers default translations for supported locales (currently English and Portuguese-BR).
- **`FormatValidationError(c, err)`**: Safely converts raw validator errors into a translated string.

---

## 🌍 Adding New Languages

1. **Create a Locale File**: Add a new JSON file in the `locales/` directory (e.g., `locales/es.json`).
   ```json
   {
     "welcome": "¡Bienvenido a Gost!",
     "errors": {
       "unauthorized": "No está autorizado."
     }
   }
   ```

2. **Register Validator Translation** (Optional): If you want localized struct validation for the new language, update `src/common/i18n/validator.go` to include the new locale and its corresponding `go-playground/validator` translations.

---

## 🛠️ How to Use

### 1. In Controllers
You can use the `i18n.Translate` helper which automatically extracts the localizer from the request context.

```go
func (ctrl *UserController) Welcome(c *gin.Context) {
    message := i18n.Translate(c, "welcome")
    c.JSON(200, gin.H{"message": message})
}
```

### 2. With Template Data (Variables)
If your translation string contains placeholders like `"hello": "Hello {{.Name}}!"`:

```go
message := i18n.Translate(c, "hello", map[string]interface{}{
    "Name": "Jonatas",
})
```

### 3. Localized Validation Errors
Gost's `ValidateBody` pipe already uses `FormatValidationError` under the hood. Any validation failure will automatically return errors in the language specified in the `Accept-Language` header.

**Request Header:** `Accept-Language: pt-BR`
**Response:**
```json
{
  "statusCode": 400,
  "error": "Bad Request",
  "message": "Nome é um campo obrigatório; "
}
```

### 4. Localized Global Exceptions
The `ErrorHandler` and `FormattedErrorGenerator` are also wired to the i18n system. If you pass a translation key as the message, it will be translated before reaching the client.

```go
utils.FormattedErrorGenerator(c, 403, "errors.unauthorized", "errors.unauthorized_detail")
```

---

## ⚙️ Configuration

The system is initialized in `src/app/app.module.go`:

```go
// Initialize translation bundle
i18n.Initialize("locales", language.English)

// Initialize validator translations
i18n.InitValidator()

// Register middleware
router.Use(i18n.Middleware())
```

---

## ✅ Best Practices

- **Keys Hierarchy**: Use dot notation in your JSON files for better organization (e.g., `user.errors.not_found`).
- **Fallback**: Always provide an `en.json` as it is the default fallback language in the framework.
- **Header usage**: Clients should send `Accept-Language` (e.g., `pt-BR, en;q=0.9`) to get the best experience.
