// Package services holds the library business logic.
package services

import "library_management/models"

// LibraryManager describes everything the library can do.
type LibraryManager interface {
	// AddBook adds a new book to the library.
	AddBook(book models.Book)
	// RemoveBook removes the book with the given ID.
	RemoveBook(bookID int)
	// BorrowBook lets a member borrow an available book.
	BorrowBook(bookID int, memberID int) error
	// ReturnBook lets a member return a book they borrowed.
	ReturnBook(bookID int, memberID int) error
	// ListAvailableBooks returns every book that is not borrowed.
	ListAvailableBooks() []models.Book
	// ListBorrowedBooks returns the books a member currently has.
	ListBorrowedBooks(memberID int) []models.Book
}

// Library is the in-memory implementation of LibraryManager.
// Methods are added on Day 4.
type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

// NewLibrary creates an empty library.
func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}
