package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"

	_ "demo-golang/docs"

	fiberSwagger "github.com/gofiber/swagger"
)

// @title Fiber CRUD API
// @version 1.0
// @description This is a simple CRUD API with Fiber.

// @contact.name API Support
// @contact.url https://sewucloud.com
// @contact.email support@sewucloud.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year,omitempty"`
}

var (
	storeMu sync.RWMutex
	store   = map[string]Book{}
)

func validateBookPayload(b *Book) error {
	if b.Title == "" {
		return errors.New("title is required")
	}
	if b.Author == "" {
		return errors.New("author is required")
	}
	return nil
}

// getAllBooks godoc
// @Summary Get all books
// @Description Get list of books with optional pagination
// @Tags books
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} map[string]interface{}
// @Router /books/ [get]
func getAllBooks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	storeMu.RLock()
	defer storeMu.RUnlock()

	books := make([]Book, 0, len(store))
	for _, v := range store {
		books = append(books, v)
	}

	start := (page - 1) * limit
	if start > len(books) {
		start = len(books)
	}
	end := start + limit
	if end > len(books) {
		end = len(books)
	}
	paged := books[start:end]

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data":  paged,
		"page":  page,
		"limit": limit,
		"total": len(books),
	})
}

// getBookByID godoc
// @Summary Get a book by ID
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} Book
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func getBookByID(c *fiber.Ctx) error {
	id := c.Params("id")
	storeMu.RLock()
	defer storeMu.RUnlock()
	b, ok := store[id]
	if !ok {
		return fiber.NewError(http.StatusNotFound, "book not found")
	}
	return c.Status(http.StatusOK).JSON(b)
}

// createBook godoc
// @Summary Create a new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body Book true "Create book"
// @Success 201 {object} Book
// @Failure 400 {object} map[string]string
// @Router /books/ [post]
func createBook(c *fiber.Ctx) error {
	var payload Book
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid JSON body")
	}
	if err := validateBookPayload(&payload); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	payload.ID = uuid.New().String()

	storeMu.Lock()
	store[payload.ID] = payload
	storeMu.Unlock()

	return c.Status(http.StatusCreated).JSON(payload)
}

// updateBook godoc
// @Summary Partially update a book
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body Book true "Update book"
// @Success 200 {object} Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [patch]
func updateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	storeMu.RLock()
	_, ok := store[id]
	storeMu.RUnlock()
	if !ok {
		return fiber.NewError(http.StatusNotFound, "book not found")
	}

	var payload Book
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid JSON body")
	}

	storeMu.Lock()
	existing := store[id]
	if payload.Title != "" {
		existing.Title = payload.Title
	}
	if payload.Author != "" {
		existing.Author = payload.Author
	}
	if payload.Year != 0 {
		existing.Year = payload.Year
	}
	store[id] = existing
	storeMu.Unlock()

	return c.Status(http.StatusOK).JSON(existing)
}

// replaceBook godoc
// @Summary Replace a book (PUT)
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body Book true "Replace book"
// @Success 200 {object} Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [put]
func replaceBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var payload Book
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid JSON body")
	}
	if err := validateBookPayload(&payload); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	payload.ID = id

	storeMu.Lock()
	if _, exists := store[id]; !exists {
		storeMu.Unlock()
		return fiber.NewError(http.StatusNotFound, "book not found")
	}
	store[id] = payload
	storeMu.Unlock()

	return c.Status(http.StatusOK).JSON(payload)
}

// deleteBook godoc
// @Summary Delete a book by ID
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /books/{id} [delete]
func deleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	storeMu.Lock()
	defer storeMu.Unlock()
	if _, ok := store[id]; !ok {
		return fiber.NewError(http.StatusNotFound, "book not found")
	}
	delete(store, id)
	return c.SendStatus(http.StatusNoContent)
}

func seedData() {
	b1 := Book{ID: uuid.New().String(), Title: "Clean Architecture", Author: "Robert C. Martin", Year: 2017}
	b2 := Book{ID: uuid.New().String(), Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Year: 2015}
	storeMu.Lock()
	store[b1.ID] = b1
	store[b2.ID] = b2
	storeMu.Unlock()
}

func main() {
	app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		if e, ok := err.(*fiber.Error); ok {
			return c.Status(e.Code).JSON(fiber.Map{"error": e.Message})
		}
		log.Println("internal error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}})

	app.Use(recover.New())
	app.Use(logger.New())

	// Swagger docs
	app.Get("/swagger/*", fiberSwagger.New())

	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("ok") })

	r := app.Group("/api")
	books := r.Group("/books")
	books.Get("/", getAllBooks)
	books.Get(":id", getBookByID)
	books.Post("/", createBook)
	books.Patch(":id", updateBook)
	books.Put(":id", replaceBook)
	books.Delete(":id", deleteBook)

	seedData()

	log.Println("listening on http://localhost:3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
