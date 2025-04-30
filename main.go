package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Kan_17578"
	dbname   = "gorm-demo"
)

func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := "TestSecret"

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Successfully connected!")

	print(db)
	db.AutoMigrate(&Book{}, &User{})
	fmt.Println("database migration connected!")

	app := fiber.New()
	app.Use("books", authRequired)

	app.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(GetBooks(db))
	})

	app.Get("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(GetBook(db, id))
	})

	app.Post("/books", func(c *fiber.Ctx) error {
		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = CreateBook(db, book)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Create book successful!",
		})

	})

	app.Put("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book := new(Book)

		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book.ID = uint(id)

		err = UpdateBook(db, book)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"message": "Update book successful",
		})
	})

	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = DeleteBook(db, id)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Delete book successful",
		})
	})

	// User API
	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		err = CreateUser(db, user)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"message": "Register user successful",
		})

	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		token, err := LoginUser(db, user)

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})
		return c.JSON(fiber.Map{
			"message": "login successful",
		})
	})

	log.Fatal(app.Listen(":8000"))

	// ---------------------------CreateBook
	// newBook := &Book{
	// 	Name:        "Karn",
	// 	Author:      "Kornchai",
	// 	Descriction: "Test3",
	// 	Price:       670,
	// }
	// CreateBook(db, newBook)
	// ---------------------------GetBook
	// currentBook := GetBook(db, 1)
	// fmt.Println(currentBook)
	// ---------------------------GetBooks
	// currentBooks := GetBooks(db)
	// fmt.Println(currentBooks)
	// ---------------------------UpdateBook
	// currentBook.Name = "New Karn"
	// currentBook.Price = 300

	// UpdateBook(db, currentBook)
	// ---------------------------DeleteBook
	// DeleteBook(db, 2)
	// ---------------------------SearchBook
	// searchBook := SearchBook(db, "Karn")
	// for _, book := range searchBook {
	// 	fmt.Println(book.ID, book.Name, book.Author, book.Price)
	// }

}
