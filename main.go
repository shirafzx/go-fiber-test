package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Book struct to hold book data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func checkMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	fmt.Printf(
		"URL = %s, Method = %s, Time = %s\n",
		c.OriginalURL(), c.Method(), start,
	)

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// name := claims["name"].(string)
	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

func main() {
  app := fiber.New()

	books = append(books, Book{ID: 1, Title: "Mikelopster", Author: "Mike"})
	books = append(books, Book{ID: 2, Title: "MM", Author: "Mike"})

	app.Post("/login", login)

	
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
  }))
	// app.Use(checkMiddleware)

	// app.Get("/books", getBooks)
	// app.Get("/books/:id", getBook)
	// app.Post("/books", createBook)
	// app.Put("/books/:id", updateBook)
	// app.Delete("/books/:id", deleteBook)

	bookGroup := app.Group("/books")

  // Apply the isAdmin middleware only to the /book routes
  bookGroup.Use(checkMiddleware)

  // Now, only authenticated admins can access these routes
  bookGroup.Get("/", getBooks)
  bookGroup.Get("/:id", getBook)
  bookGroup.Post("/", createBook)
  bookGroup.Put("/:id", updateBook)
  bookGroup.Delete("/:id", deleteBook)

	app.Listen(":8080")
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var memberUser = User {
	Email: "user@example.com",
	Password: "password123",
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user) ; err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email != memberUser.Email || user.Password != memberUser.Password {
		return fiber.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	 // Set claims
	 claims := token.Claims.(jwt.MapClaims)
	 claims["email"] = "user@example.com"
	 claims["role"] = "admin"
	 claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token
  t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
  if err != nil {
    return c.SendStatus(fiber.StatusInternalServerError)
  }

	return c.JSON(fiber.Map{"token": t})
}

