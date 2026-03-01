package benchmark

import (
	"regexp"
	"testing"

	"github.com/sivchari/govalid/test"
)

func BenchmarkGoValidPattern(b *testing.B) {
	instance := test.Pattern{Username: "johndoe"}
	b.ResetTimer()
	for b.Loop() {
		err := test.ValidatePattern(&instance)
		if err != nil {
			b.Fatal("unexpected error:", err)
		}
	}
	b.StopTimer()
}

func BenchmarkGoValidPatternCSV(b *testing.B) {
	instance := test.PatternCSV{Tags: "alpha,beta,gamma,delta"}
	b.ResetTimer()
	for b.Loop() {
		err := test.ValidatePatternCSV(&instance)
		if err != nil {
			b.Fatal("unexpected error:", err)
		}
	}
	b.StopTimer()
}

func BenchmarkGoValidPatternPhone(b *testing.B) {
	instance := test.PatternPhone{Phone: "123-456-7890"}
	b.ResetTimer()
	for b.Loop() {
		err := test.ValidatePatternPhone(&instance)
		if err != nil {
			b.Fatal("unexpected error:", err)
		}
	}
	b.StopTimer()
}

func BenchmarkGoValidPatternCode(b *testing.B) {
	instance := test.PatternCode{Code: "AB1234"}
	b.ResetTimer()
	for b.Loop() {
		err := test.ValidatePatternCode(&instance)
		if err != nil {
			b.Fatal("unexpected error:", err)
		}
	}
	b.StopTimer()
}

func BenchmarkGoRegexpPattern(b *testing.B) {
	re := regexp.MustCompile(`^[a-z]+$`)
	value := "johndoe"
	b.ResetTimer()
	for b.Loop() {
		if !re.MatchString(value) {
			b.Fatal("validation failed")
		}
	}
	b.StopTimer()
}

func BenchmarkGoRegexpPatternNoCache(b *testing.B) {
	value := "johndoe"
	b.ResetTimer()
	for b.Loop() {
		re := regexp.MustCompile(`^[a-z]+$`)
		if !re.MatchString(value) {
			b.Fatal("validation failed")
		}
	}
	b.StopTimer()
}
