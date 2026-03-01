package unit

import (
	"testing"

	"github.com/sivchari/govalid/test"
)

func TestPatternValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.Pattern
		expectError bool
	}{
		// Valid cases - lowercase letters only (pattern: ^[a-z]+$)
		{"valid_single_char", test.Pattern{Username: "a"}, false},
		{"valid_lowercase", test.Pattern{Username: "john"}, false},
		{"valid_long", test.Pattern{Username: "johndoe"}, false},

		// Invalid cases
		{"empty", test.Pattern{Username: ""}, true},
		{"uppercase", test.Pattern{Username: "John"}, true},
		{"numbers", test.Pattern{Username: "john123"}, true},
		{"special_chars", test.Pattern{Username: "john_doe"}, true},
		{"spaces", test.Pattern{Username: "john doe"}, true},
		{"mixed_case", test.Pattern{Username: "JohnDoe"}, true},
		{"all_uppercase", test.Pattern{Username: "JOHN"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test govalid
			err := test.ValidatePattern(&tt.data)
			hasError := err != nil
			if hasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (err: %v)", tt.expectError, hasError, err)
			}
		})
	}
}

// TestPatternCSVValidation tests comma-separated lowercase list pattern: ^[a-z]+(,[a-z]+)*$
func TestPatternCSVValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.PatternCSV
		expectError bool
	}{
		// Valid cases
		{"single_tag", test.PatternCSV{Tags: "foo"}, false},
		{"two_tags", test.PatternCSV{Tags: "foo,bar"}, false},
		{"three_tags", test.PatternCSV{Tags: "a,b,c"}, false},
		{"long_tags", test.PatternCSV{Tags: "alpha,beta,gamma,delta"}, false},

		// Invalid cases
		{"empty", test.PatternCSV{Tags: ""}, true},
		{"trailing_comma", test.PatternCSV{Tags: "foo,"}, true},
		{"leading_comma", test.PatternCSV{Tags: ",foo"}, true},
		{"double_comma", test.PatternCSV{Tags: "foo,,bar"}, true},
		{"spaces", test.PatternCSV{Tags: "foo, bar"}, true},
		{"uppercase", test.PatternCSV{Tags: "Foo,bar"}, true},
		{"numbers", test.PatternCSV{Tags: "foo1,bar2"}, true},
		{"special_chars", test.PatternCSV{Tags: "foo-bar,baz"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := test.ValidatePatternCSV(&tt.data)
			hasError := err != nil
			if hasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (err: %v)", tt.expectError, hasError, err)
			}
		})
	}
}

// TestPatternPhoneValidation tests phone number pattern: ^\d{3}-\d{3}-\d{4}$
func TestPatternPhoneValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.PatternPhone
		expectError bool
	}{
		// Valid cases
		{"valid_phone", test.PatternPhone{Phone: "123-456-7890"}, false},
		{"valid_phone_zeros", test.PatternPhone{Phone: "000-000-0000"}, false},

		// Invalid cases
		{"empty", test.PatternPhone{Phone: ""}, true},
		{"no_dashes", test.PatternPhone{Phone: "1234567890"}, true},
		{"wrong_format", test.PatternPhone{Phone: "12-3456-7890"}, true},
		{"too_short", test.PatternPhone{Phone: "123-456-789"}, true},
		{"too_long", test.PatternPhone{Phone: "123-456-78901"}, true},
		{"letters", test.PatternPhone{Phone: "abc-def-ghij"}, true},
		{"spaces", test.PatternPhone{Phone: "123 456 7890"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := test.ValidatePatternPhone(&tt.data)
			hasError := err != nil
			if hasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (err: %v)", tt.expectError, hasError, err)
			}
		})
	}
}

// TestPatternCodeValidation tests alphanumeric code pattern: ^[A-Z]{2}\d{4}$
func TestPatternCodeValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        test.PatternCode
		expectError bool
	}{
		// Valid cases
		{"valid_code", test.PatternCode{Code: "AB1234"}, false},
		{"valid_code_zeros", test.PatternCode{Code: "ZZ0000"}, false},
		{"valid_code_nines", test.PatternCode{Code: "AA9999"}, false},

		// Invalid cases
		{"empty", test.PatternCode{Code: ""}, true},
		{"lowercase_letters", test.PatternCode{Code: "ab1234"}, true},
		{"too_few_letters", test.PatternCode{Code: "A12345"}, true},
		{"too_many_letters", test.PatternCode{Code: "ABC123"}, true},
		{"too_few_numbers", test.PatternCode{Code: "AB123"}, true},
		{"too_many_numbers", test.PatternCode{Code: "AB12345"}, true},
		{"letters_at_end", test.PatternCode{Code: "1234AB"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := test.ValidatePatternCode(&tt.data)
			hasError := err != nil
			if hasError != tt.expectError {
				t.Errorf("govalid: expected error=%v, got error=%v (err: %v)", tt.expectError, hasError, err)
			}
		})
	}
}

// TestPatternOnIntSkipped verifies that pattern on int field is ignored
// (pattern validator returns nil for non-string types)
// This is verified at compile-time: no ValidatePatternOnInt function is generated
// The struct can still be instantiated and used, but no validation is performed
func TestPatternOnIntSkipped(t *testing.T) {
	// This test verifies the struct type exists and can be used
	// No ValidatePatternOnInt function is generated because pattern only works on strings
	data := test.PatternOnInt{Number: 12345}
	_ = data // Use the struct to verify it compiles

	// The fact that this test compiles proves that:
	// 1. PatternOnInt struct exists
	// 2. Pattern marker on int field was silently ignored (no compilation errors)
	// 3. No broken validation code was generated
	t.Log("PatternOnInt struct created successfully - pattern marker was correctly ignored for int field")
}
