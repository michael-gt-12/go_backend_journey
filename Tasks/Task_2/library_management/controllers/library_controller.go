// Package controllers handles console input and output.
package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"library_management/models"
	"library_management/services"
)

// LibraryController runs the console menu on top of a LibraryManager.
type LibraryController struct {
	library services.LibraryManager
	reader  *bufio.Reader
}

// NewLibraryController builds a controller that reads from the keyboard.
func NewLibraryController(lib services.LibraryManager) *LibraryController {
	return &LibraryController{library: lib, reader: bufio.NewReader(os.Stdin)}
}

// readLine prints a prompt and returns the trimmed line the user typed.
func (c *LibraryController) readLine(prompt string) string {
	fmt.Print(prompt)
	text, err := c.reader.ReadString('\n')
	if err != nil { // input closed (e.g. Ctrl+Z / Ctrl+D)
		fmt.Println("\nInput closed. Goodbye!")
		os.Exit(0)
	}
	return strings.TrimSpace(text)
}

// readInt asks for a number; ok is false if the input was not a valid number.
func (c *LibraryController) readInt(prompt string) (int, bool) {
	n, err := strconv.Atoi(c.readLine(prompt))
	if err != nil {
		fmt.Println("Please enter a valid number.")
		return 0, false
	}
	return n, true
}

func (c *LibraryController) printMenu() {
	fmt.Println("\n=== Library Management System ===")
	fmt.Println("1. Add a book")
	fmt.Println("2. Remove a book")
	fmt.Println("3. Borrow a book")
	fmt.Println("4. Return a book")
	fmt.Println("5. List available books")
	fmt.Println("6. List books borrowed by a member")
	fmt.Println("7. Exit")
}

// Run shows the menu repeatedly until the user chooses Exit.
func (c *LibraryController) Run() {
	for {
		c.printMenu()
		choice, ok := c.readInt("Choose an option: ")
		if !ok {
			continue
		}
		switch choice {
		case 1:
			c.addBook()
		case 2:
			c.removeBook()
		case 3:
			c.borrowBook()
		case 4:
			c.returnBook()
		case 5:
			c.listAvailable()
		case 6:
			c.listBorrowed()
		case 7:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option, choose a number from 1 to 7.")
		}
	}
}

func (c *LibraryController) addBook() {
	id, ok := c.readInt("Book ID: ")
	if !ok {
		return
	}
	title := c.readLine("Title: ")
	author := c.readLine("Author: ")
	if title == "" || author == "" {
		fmt.Println("Title and author cannot be empty.")
		return
	}
	c.library.AddBook(models.Book{ID: id, Title: title, Author: author, Status: models.StatusAvailable})
	fmt.Println("Book added.")
}

// removeBook asks for a book ID and removes it. RemoveBook returns no error
// (per the spec's interface), so we cannot report whether the ID existed.
func (c *LibraryController) removeBook() {
	bookID, ok := c.readInt("Book ID: ")
	if !ok {
		return
	}
	c.library.RemoveBook(bookID)
	fmt.Println("Remove request done (nothing happens if the ID did not exist).")
}

func (c *LibraryController) borrowBook() {
	bookID, ok := c.readInt("Book ID: ")
	if !ok {
		return
	}
	memberID, ok := c.readInt("Member ID: ")
	if !ok {
		return
	}
	if err := c.library.BorrowBook(bookID, memberID); err != nil {
		fmt.Println("Could not borrow:", err)
		return
	}
	fmt.Println("Book borrowed.")
}

func (c *LibraryController) returnBook() {
	bookID, ok := c.readInt("Book ID: ")
	if !ok {
		return
	}
	memberID, ok := c.readInt("Member ID: ")
	if !ok {
		return
	}
	if err := c.library.ReturnBook(bookID, memberID); err != nil {
		fmt.Println("Could not return:", err)
		return
	}
	fmt.Println("Book returned.")
}

func (c *LibraryController) listAvailable() {
	books := c.library.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("No books are available.")
		return
	}
	for _, b := range books {
		fmt.Printf("[%d] %s by %s\n", b.ID, b.Title, b.Author)
	}
}

// listBorrowed asks for a member ID and prints the books that member has.
// An empty result means either no borrowed books or an unknown member,
// because ListBorrowedBooks cannot tell the two apart.
func (c *LibraryController) listBorrowed() {
	memberID, ok := c.readInt("Member ID: ")
	if !ok {
		return
	}
	books := c.library.ListBorrowedBooks(memberID)
	if len(books) == 0 {
		fmt.Println("No borrowed books found for this member.")
		return
	}
	for _, b := range books {
		fmt.Printf("[%d] %s by %s\n", b.ID, b.Title, b.Author)
	}
}
