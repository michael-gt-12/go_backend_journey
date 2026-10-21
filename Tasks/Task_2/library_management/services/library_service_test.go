package services

import (
	"library_management/models"
	"testing"
)

func TestListAvailableBooks(t *testing.T) {
	lib := NewLibrary()
	lib.AddBook(models.Book{ID: 1, Title: "A", Status: models.StatusAvailable})
	lib.AddBook(models.Book{ID: 2, Title: "B", Status: models.StatusBorrowed})
	lib.AddBook(models.Book{ID: 3, Title: "C", Status: models.StatusAvailable})

	got := lib.ListAvailableBooks()

	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 3 {
		t.Fatalf("expected books 1 and 3 in order, got %+v", got)
	}
}

func TestReturnBook(t *testing.T) {
	setup := func() *Library {
		lib := NewLibrary()
		lib.AddBook(models.Book{ID: 1, Title: "A", Status: models.StatusAvailable})
		lib.AddMember(models.Member{ID: 1, Name: "Abel"})
		lib.AddMember(models.Member{ID: 2, Name: "Sara"})
		_ = lib.BorrowBook(1, 1) // Abel borrows book 1
		return lib
	}

	tests := []struct {
		name     string
		bookID   int
		memberID int
		wantErr  error
	}{
		{"success", 1, 1, nil},
		{"book not found", 99, 1, ErrorBookNotFound},
		{"member not found", 1, 99, ErrorMemberNotFound},
		{"not borrowed by this member", 1, 2, ErrorNotBorrowedByThisMember},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib := setup()
			err := lib.ReturnBook(tt.bookID, tt.memberID)
			if err != tt.wantErr {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}
