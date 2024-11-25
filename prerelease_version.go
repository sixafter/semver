// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Error messages for prerelease validation.
var (
	ErrEmptyPrerelease        = errors.New("prerelease is empty")
	ErrLeadingZeroInNumeric   = "numeric prerelease version must not contain leading zeroes: %q"
	ErrInvalidPrereleaseChars = "invalid character(s) found in prerelease: %q"
)

// PrereleaseVersion represents a semantic version prerelease identifier.
//
// A prerelease version can be either numeric or alphanumeric.
// Numeric prerelease versions have lower precedence than alphanumeric ones.
type PrereleaseVersion struct {
	partString  string
	partNumeric uint64
	isNumeric   bool
}

// NewPrereleaseVersion creates a new valid PrereleaseVersion from a string.
//
// Returns an error if the string is empty, contains invalid characters, or if
// a numeric prerelease version has leading zeroes.
//
// Example:
//
//	v, err := semver.NewPrereleaseVersion("alpha.1")
//	if err != nil {
//	  fmt.Println("Error:", err)
//	} else {
//	  fmt.Println(v.String()) // Output: alpha.1
//	}
//
//	v2, err := semver.NewPrereleaseVersion("01")
//	if err != nil {
//	  fmt.Println("Error:", err) // Output: numeric prerelease version must not contain leading zeroes: "01"
//	}
func NewPrereleaseVersion(s string) (PrereleaseVersion, error) {
	if len(s) == 0 {
		return PrereleaseVersion{}, errors.New("prerelease is empty")
	}

	// Check if the string contains only numbers
	if containsOnlyNumbers(s) {
		// Check for leading zeroes
		if len(s) > 1 && s[0] == '0' {
			return PrereleaseVersion{}, fmt.Errorf("numeric prerelease version must not contain leading zeroes: %q", s)
		}

		// Parse numeric string
		number, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return PrereleaseVersion{}, err
		}

		return PrereleaseVersion{
			partNumeric: number,
			isNumeric:   true,
		}, nil
	}

	// Check if the string contains only alphanumeric characters
	if containsOnlyAlphanumeric(s) {
		return PrereleaseVersion{
			partString: s,
			isNumeric:  false,
		}, nil
	}

	// If neither numeric nor alphanumeric, return an error
	return PrereleaseVersion{}, fmt.Errorf("invalid character(s) found in prerelease: %q", s)
}

// IsNumeric checks if the prerelease version is numeric.
//
// Example:
//
//	v, _ := semver.NewPrereleaseVersion("123")
//	fmt.Println(v.IsNumeric()) // Output: true
//
//	v2, _ := semver.NewPrereleaseVersion("alpha")
//	fmt.Println(v2.IsNumeric()) // Output: false
func (v PrereleaseVersion) IsNumeric() bool {
	return v.isNumeric
}

// Compare compares two PrereleaseVersion instances.
//
// Returns:
//   - -1 if v < o
//   - 0 if v == o
//   - 1 if v > o
//
// Numeric prerelease versions are always less than non-numeric ones.
// Numeric versions are compared numerically; alphanumeric versions are compared lexicographically.
//
// Example:
//
//	v1, _ := semver.NewPrereleaseVersion("123")
//	v2, _ := semver.NewPrereleaseVersion("alpha")
//	fmt.Println(v1.Compare(v2)) // Output: -1
//
//	v3, _ := semver.NewPrereleaseVersion("alpha")
//	v4, _ := semver.NewPrereleaseVersion("beta")
//	fmt.Println(v3.Compare(v4)) // Output: -1
func (v PrereleaseVersion) Compare(o PrereleaseVersion) int {
	// Numeric identifiers have lower precedence than non-numeric identifiers
	if v.isNumeric != o.isNumeric {
		if v.isNumeric {
			return -1
		}
		return 1
	}

	// If both are numeric, compare numerically
	if v.isNumeric {
		if v.partNumeric < o.partNumeric {
			return -1
		} else if v.partNumeric > o.partNumeric {
			return 1
		} else {
			return 0
		}
	}

	// If both are non-numeric, compare lexicographically (ASCII sort order)
	return strings.Compare(v.partString, o.partString)
}

// String returns the string representation of the PrereleaseVersion.
//
// Example:
//
//	v, _ := semver.NewPrereleaseVersion("123")
//	fmt.Println(v.String()) // Output: 123
//
//	v2, _ := semver.NewPrereleaseVersion("alpha")
//	fmt.Println(v2.String()) // Output: alpha
func (v PrereleaseVersion) String() string {
	if v.isNumeric {
		return strconv.FormatUint(v.partNumeric, 10)
	}
	return v.partString
}

// containsOnlyNumbers checks if the string contains only numeric characters.
//
// Example:
//
//	fmt.Println(containsOnlyNumbers("123")) // Output: true
//	fmt.Println(containsOnlyNumbers("abc")) // Output: false
func containsOnlyNumbers(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// containsOnlyAlphanumeric checks if the string contains only ASCII letters and numbers.
//
// Example:
//
//	fmt.Println(containsOnlyAlphanumeric("abc123")) // Output: true
//	fmt.Println(containsOnlyAlphanumeric("abc-123")) // Output: false
func containsOnlyAlphanumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}

// hasLeadingZeroes checks if the string has leading zeroes.
//
// Example:
//
//	fmt.Println(hasLeadingZeroes("0123")) // Output: true
//	fmt.Println(hasLeadingZeroes("123"))  // Output: false
func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}
