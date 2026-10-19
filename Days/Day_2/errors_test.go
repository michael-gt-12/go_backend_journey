package day2

import (
	"errors"
	"testing"
)

func TestValidateAge(t *testing.T) {
	tests := []struct {
		name        string
		age         int
		wantErr     bool
		wantMessage string
	}{
		{"valid zero", 0, false, ""},
		{"valid adult", 25, false, ""},
		{"valid upper bound", 150, false, ""},
		{"negative", -1, true, "must not be negative"},
		{"too large", 151, true, "must not be greater than 150"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAge(tc.age)

			if !tc.wantErr {
				if err != nil {
					t.Fatalf("ValidateAge(%d) = %v, want nil", tc.age, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("ValidateAge(%d) = nil, want error", tc.age)
			}

			// errors.As checks whether the error is (or wraps) a *ValidationError.
			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("error %v is not a *ValidationError", err)
			}
			if ve.Field != "age" {
				t.Errorf("Field = %q, want %q", ve.Field, "age")
			}
			if ve.Message != tc.wantMessage {
				t.Errorf("Message = %q, want %q", ve.Message, tc.wantMessage)
			}
		})
	}
}

func TestValidationErrorString(t *testing.T) {
	err := &ValidationError{Field: "name", Message: "is required"}
	want := `validation failed on "name": is required`
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
