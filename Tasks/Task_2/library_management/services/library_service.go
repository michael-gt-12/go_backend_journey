// Package services holds the library business logic.
package services

import (
	"errors"
	"library_management/models"
	"sort"
)

// 4 error types that i am going to use
var ErrorBookNotFound = errors.New("Book not found")
var ErrorMemberNotFound = errors.New("Member not found")
var ErrorAlreadyBorrowed = errors.New("Book is already Borrowed by another member")
var ErrorNotBorrowedByThisMember = errors.New("This Book is not Borrowed by this member")

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

var _ LibraryManager = (*Library)(nil)

// Library is the in-memory implementation of LibraryManager.
// Methods are added on Day 4.
type Library struct {
	books   map[int]models.Book
	members map[int]models.Member
}

// implementation of AddBook method

func (lib *Library) AddBook(book models.Book) {
	// don't i have to check if a book with the same id exists?
	id := book.ID
	lib.books[id] = book
}

// implementation of RemoveBook
func (lib *Library) RemoveBook(bookID int) {
	delete(lib.books, bookID)
}

// implementation of ListAvalibleBooks
func (lib *Library) ListAvailableBooks() []models.Book {
	var result []models.Book
	books := lib.books
	keys := make([]int, 0, len(books))
	for k := range books {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		if books[key].Status == models.StatusAvailable {
			result = append(result, books[key])
		}
	}
	return result

}

// implementation of BorrowBook
func (lib *Library) BorrowBook(bookID int, memberID int) error {
	books := lib.books
	members := lib.members
	_, bookExists := books[bookID]
	_, memberExists := members[memberID]
	if bookExists {
		book := books[bookID]
		if book.Status == models.StatusBorrowed {
			return ErrorAlreadyBorrowed
		}

		if memberExists {
			member := members[memberID]
			book.Status = models.StatusBorrowed
			member.BorrowedBooks = append(member.BorrowedBooks, book)
			books[bookID] = book
			members[memberID] = member
			return nil
		}
		return ErrorMemberNotFound
	}
	return ErrorBookNotFound
	// scratch work
	/* 	if bookExists && memberExists {
	   		book := books[bookID]
	   		member := members[memberID]
	   		member.BorrowedBooks = append(member.BorrowedBooks, book)
	   		book.Status = models.StatusBorrowed
	   		books[bookID] = book
	   		return nil
	   	}
	*/
}

// implementation of ReturnBook
func (lib *Library) ReturnBook(bookID int, memberID int) error {
	members := lib.members
	books := lib.books
	_, bookExists := books[bookID]
	_, memberExists := members[memberID]

	if !bookExists {
		return ErrorBookNotFound
	}
	if !memberExists {
		return ErrorMemberNotFound
	}

	member := members[memberID]
	book_ := books[bookID]
	BorrowedBooks := member.BorrowedBooks

	for i, book := range BorrowedBooks {
		id := book.ID
		if id == bookID {
			b_books := BorrowedBooks[i+1:]
			member.BorrowedBooks = append(BorrowedBooks[:i], b_books...)
			book_.Status = models.StatusAvailable

			books[bookID] = book_
			members[memberID] = member
			return nil
		}
	}
	return ErrorNotBorrowedByThisMember
}

// implementation of ListOfBorrowedBooks
func (lib *Library) ListBorrowedBooks(memberID int) []models.Book {
	members := lib.members
	return members[memberID].BorrowedBooks
}

// implementation of adding a member
func (lib *Library) AddMember(member models.Member) {
	members := lib.members
	id := member.ID
	_, exists := members[id]
	if !exists {
		members[id] = member
	}

}

// NewLibrary creates an empty library.
func NewLibrary() *Library {
	return &Library{
		books:   make(map[int]models.Book),
		members: make(map[int]models.Member),
	}
}
