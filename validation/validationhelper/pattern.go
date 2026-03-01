// Package validationhelper provides validation helper functions for govalid.
package validationhelper

import (
	"regexp"
	"sync"
)

var (
	patternCache   = make(map[string]*regexp.Regexp)
	patternCacheMu sync.RWMutex
)

// MatchPattern checks if value matches the given regex pattern.
// Patterns are compiled once and cached for performance.
//
// This function uses a thread-safe cache to avoid recompiling the same pattern
// multiple times. The pattern is compiled on first use and stored for subsequent calls.
//
// Parameters:
//   - pattern: A valid Go regex pattern (re2 syntax)
//   - value: The string to match against the pattern
//
// Returns:
//   - true if the value matches the pattern
//   - false if the value doesn't match or the pattern is invalid
//
// Example patterns:
//   - "^[a-z]+$"      matches lowercase letters only
//   - "^\\d{3}-\\d{4}$" matches phone numbers like "123-4567"
//   - "^[A-Z]{2}\\d{4}$" matches codes like "AB1234"
func MatchPattern(pattern, value string) bool {
	// Try to get cached compiled pattern with read lock
	patternCacheMu.RLock()
	re, ok := patternCache[pattern]
	patternCacheMu.RUnlock()

	if !ok {
		// Pattern not cached, acquire write lock to compile and cache
		patternCacheMu.Lock()

		// Double-check after acquiring write lock (another goroutine may have cached it)
		re, ok = patternCache[pattern]
		if !ok {
			var err error

			re, err = regexp.Compile(pattern)
			if err != nil {
				// Invalid pattern - treat as no match
				patternCacheMu.Unlock()

				return false
			}

			patternCache[pattern] = re
		}

		patternCacheMu.Unlock()
	}

	return re.MatchString(value)
}
