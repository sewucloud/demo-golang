## Dependencies

Project ini menggunakan beberapa package berikut:

- [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2) — Web framework
- [github.com/gofiber/fiber/v2/middleware/logger](https://pkg.go.dev/github.com/gofiber/fiber/v2/middleware/logger) — Middleware logging
- [github.com/gofiber/fiber/v2/middleware/recover](https://pkg.go.dev/github.com/gofiber/fiber/v2/middleware/recover) — Middleware recover panic
- [github.com/google/uuid](https://pkg.go.dev/github.com/google/uuid) — UUID generator
- [github.com/gofiber/swagger](https://github.com/gofiber/swagger) — Swagger UI untuk Fiber
- [github.com/swaggo/swag/cmd/swag](https://github.com/swaggo/swag) — CLI untuk generate dokumentasi Swagger

## Installation

Jalankan perintah berikut untuk menginstall dependencies:

```bash
go get github.com/gofiber/fiber/v2
go get github.com/gofiber/fiber/v2/middleware/logger
go get github.com/gofiber/fiber/v2/middleware/recover
go get github.com/google/uuid
go get github.com/gofiber/swagger
go install github.com/swaggo/swag/cmd/swag@latest