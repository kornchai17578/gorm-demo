package main

import (
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func CreateBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetBook(db *gorm.DB, id int) *Book {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error get book: %v", result.Error)
	}
	return &book
}

func GetBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		log.Fatalf("Error get book: %v", result.Error)
	}
	return books
}

func UpdateBook(db *gorm.DB, book *Book) error {
	result := db.Model(&book).Updates(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteBook(db *gorm.DB, id int) error {
	var book Book
	result := db.Delete(&book, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SearchBook(db *gorm.DB, bookName string) []Book {
	var books []Book
	result := db.Where("name = ?", bookName).Order("price desc").Find(&books)
	if result.Error != nil {
		log.Fatalf("Error search book: %v", result.Error)
	}
	return books
}
