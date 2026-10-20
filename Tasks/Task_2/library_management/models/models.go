// Package models defines the data structures of the library system.
package models

// Allowed values for Book.Status.
const (
	StatusAvailable = "Available"
	StatusBorrowed  = "Borrowed"
)

// Book represents a book in the library.
type Book struct {
	ID     int
	Title  string
	Author string
	Status string // StatusAvailable or StatusBorrowed
}

// Member represents a library member.
type Member struct {
	ID            int
	Name          string
	BorrowedBooks []Book
}
