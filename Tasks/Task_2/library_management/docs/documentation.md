# Library Management System

A console-based library system written in Go. It manages books and members
in memory: books can be added, removed, borrowed and returned, and the
program can list available books and the books a member has borrowed.

## Folder structure

```
library_management/
├── main.go                         # entry point: seeds members, starts the menu
├── go.mod
├── models/
│   └── models.go                   # Book and Member structs, status constants
├── services/
│   ├── library_service.go          # LibraryManager interface + Library implementation
│   └── library_service_test.go     # table-driven tests for the service
├── controllers/
│   └── libray_controller.go        # console menu, input reading and validation
└── docs/
    └── documentation.md            # this file
```

## How it works

- **models**: `Book` (ID, Title, Author, Status) and `Member` (ID, Name,
  BorrowedBooks). A book's status is either `Available` or `Borrowed`.
- **services**: the `LibraryManager` interface defines six operations
  (AddBook, RemoveBook, BorrowBook, ReturnBook, ListAvailableBooks,
  ListBorrowedBooks). `Library` implements it using two maps:
  `map[int]models.Book` and `map[int]models.Member`.
  `var _ LibraryManager = (*Library)(nil)` makes the compiler verify that
  `Library` satisfies the interface.
- **controllers**: `LibraryController` depends on the `LibraryManager`
  interface (not on the concrete `Library`), reads input with `bufio`, and
  runs the menu loop.
- **main.go**: creates the library, seeds two members with `AddMember`
  (which is a method on `*Library`, not part of the interface), and starts
  the controller.

Go maps hold copies of structs, so every method that changes a `Book` or
`Member` writes the updated value back into the map.

## How to run

Requires Go (see `go.mod` for the version).

```bash
cd Tasks/Task_2/library_management
go run main.go
```

Seeded members: `1 - Abel`, `2 - Sara`. Books must be added from the menu
(option 1) before they can be borrowed.

### Menu

| Option | Action                              |
| ------ | ----------------------------------- |
| 1      | Add a book (ID, title, author)      |
| 2      | Remove a book by ID                 |
| 3      | Borrow a book (book ID, member ID)  |
| 4      | Return a book (book ID, member ID)  |
| 5      | List available books (sorted by ID) |
| 6      | List books borrowed by a member     |
| 7      | Exit                                |

## How to test

```bash
go test ./...
```

The service tests are table-driven and cover the success and failure cases
of `BorrowBook`, `ReturnBook` and `ListAvailableBooks`.

## Error cases

| Situation                                | Where handled              | Result                                                   |
| ---------------------------------------- | -------------------------- | -------------------------------------------------------- |
| Borrow/return with an unknown book ID    | `BorrowBook`, `ReturnBook` | `ErrorBookNotFound` ("Book not found")                   |
| Borrow/return with an unknown member ID  | `BorrowBook`, `ReturnBook` | `ErrorMemberNotFound` ("Member not found")               |
| Borrow a book that is already borrowed   | `BorrowBook`               | `ErrorAlreadyBorrowed`                                   |
| Return a book the member did not borrow  | `ReturnBook`               | `ErrorNotBorrowedByThisMember`                           |
| Non-numeric ID or menu choice            | controller `readInt`       | "Please enter a valid number." and the prompt repeats    |
| Menu choice outside 1-7                  | controller `Run`           | "Invalid option..." message                              |
| Empty title or author when adding a book | controller `addBook`       | "Title and author cannot be empty." and nothing is added |
| Input stream closed (Ctrl+Z / Ctrl+D)    | controller `readLine`      | Prints a goodbye message and exits cleanly               |

The controller prints the error message from the service, so the user sees
a friendly line such as `Could not borrow: Book is already Borrowed by another member`.

## Design decisions and known limitations

- **Interface from the spec:** `AddBook` and `RemoveBook` return no error, so
  the controller cannot tell the user whether an ID already existed or was
  missing. The menu prints a neutral message instead.
- **Duplicate book IDs:** `AddBook` currently overwrites an existing book with
  the same ID. This means re-adding the ID of a borrowed book would reset it to
  its new status. `AddMember` does the opposite and ignores duplicate IDs.
- **Removing a borrowed book:** `RemoveBook` deletes the book from the library
  but does not remove the copy in the borrowing member's `BorrowedBooks`.
- **Empty list ambiguity:** `ListBorrowedBooks` returns an empty list both for
  a member with no books and for an unknown member.
- **In-memory only:** all data is lost when the program exits. Persistence
  comes in later tasks.
