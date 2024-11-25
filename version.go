// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Error messages for version parsing
var (
	ErrEmptyVersionString        = errors.New("version string empty")
	ErrInvalidMajorNumber        = "invalid character(s) found in major number %q"
	ErrLeadingZeroMajor          = "major number must not contain leading zeroes %q"
	ErrInvalidMinorNumber        = "invalid character(s) found in minor number %q"
	ErrLeadingZeroMinor          = "minor number must not contain leading zeroes %q"
	ErrMissingVersionElements    = errors.New("missing major, minor, or patch elements")
	ErrInvalidPatchNumber        = "invalid character(s) found in patch number %q"
	ErrLeadingZeroPatch          = "patch number must not contain leading zeroes %q"
	ErrEmptyBuildMetadata        = errors.New("build metadata is empty")
	ErrInvalidBuildMetadataChars = "invalid character(s) found in build metadata %q"
)

// SupportedVersion is the latest fully supported specification version of semver.
//
// Example:
//
//	fmt.Println(semver.SupportedVersion) // Output: 2.0.0
var SupportedVersion = Version{
	Major: 2,
	Minor: 0,
	Patch: 0,
}

// Version represents a Semantic Versioning 2.0.0 version.
//
// A Version includes major, minor, and patch numbers, as well as optional pre-release and build metadata.
//
// Example:
//
//	v := semver.Version{
//	  Major: 1,
//	  Minor: 2,
//	  Patch: 3,
//	  PreRelease: []PrereleaseVersion{{partString: "alpha"}, {partNumeric: 1, isNumeric: true}},
//	  BuildMetadata: []string{"build", "123"},
//	}
//	fmt.Println(v.String()) // Output: 1.2.3-alpha.1+build.123
type Version struct {
	BuildMetadata []string
	PreRelease    []PrereleaseVersion
	Major         uint64
	Minor         uint64
	Patch         uint64
}

// NewVersion creates a new Version instance with the specified major, minor, patch components,
// optional prerelease identifiers, and optional build metadata.
//
// The prerelease identifiers determine the precedence of the version relative to other versions with the same
// major, minor, and patch numbers. Build metadata provides additional information but does not affect version precedence.
//
// Parameters:
//   - major: The major version number (uint64).
//   - minor: The minor version number (uint64).
//   - patch: The patch version number (uint64).
//   - preRelease: A slice of PrereleaseVersion structs representing the prerelease identifiers (optional).
//   - buildMetadata: A slice of strings representing the build metadata (optional).
//
// Returns:
//   - A Version struct populated with the provided values.
//
// Example:
//
//	package main
//
//	import (
//	    "fmt"
//	    "semver"
//	)
//
//	func main() {
//	    preRelease := []semver.PrereleaseVersion{
//	        {partString: "alpha", isNumeric: false},
//	        {partNumeric: 1, isNumeric: true},
//	    }
//	    buildMetadata := []string{"build", "2024"}
//
//	    v := semver.NewVersion(1, 2, 3, preRelease, buildMetadata)
//
//	    fmt.Println(v.String()) // Output: 1.2.3-alpha.1+build.2024
//	}
func NewVersion(major, minor, patch uint64, preRelease []PrereleaseVersion, buildMetadata []string) Version {
	return Version{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreRelease:    preRelease,
		BuildMetadata: buildMetadata,
	}
}

// Parse parses a version string into a Version struct.
//
// Returns an error if the version string is not a valid semantic version.
//
// Example:
//
//	v, err := semver.Parse("1.2.3-alpha.1+build.123")
//	if err != nil {
//	  fmt.Println("Error:", err)
//	} else {
//	  fmt.Println(v) // Output: 1.2.3-alpha.1+build.123
//	}
func Parse(version string) (Version, error) {
	if len(version) == 0 {
		return Version{}, ErrEmptyVersionString
	}

	v := Version{}
	var start, dotCount int

	// Parse Major, Minor, and Patch components by identifying '.' separators.
	for i := 0; i < len(version); i++ {
		if version[i] == '.' {
			switch dotCount {
			case 0: // Major
				if !containsOnlyNumbers(version[start:i]) {
					return Version{}, fmt.Errorf(ErrInvalidMajorNumber, version[start:i])
				}
				if hasLeadingZeroes(version[start:i]) {
					return Version{}, fmt.Errorf(ErrLeadingZeroMajor, version[start:i])
				}
				major, err := strconv.ParseUint(version[start:i], 10, 64)
				if err != nil {
					return Version{}, err
				}
				v.Major = major
				start = i + 1
			case 1: // Minor
				if !containsOnlyNumbers(version[start:i]) {
					return Version{}, fmt.Errorf(ErrInvalidMinorNumber, version[start:i])
				}
				if hasLeadingZeroes(version[start:i]) {
					return Version{}, fmt.Errorf(ErrLeadingZeroMinor, version[start:i])
				}
				minor, err := strconv.ParseUint(version[start:i], 10, 64)
				if err != nil {
					return Version{}, err
				}
				v.Minor = minor
				start = i + 1
			}
			dotCount++
		} else if version[i] == '-' || version[i] == '+' {
			break
		}
	}

	if dotCount != 2 {
		return Version{}, ErrMissingVersionElements
	}

	// Parse Patch
	i := start
	for i < len(version) && version[i] != '-' && version[i] != '+' {
		i++
	}
	patchStr := version[start:i]
	if !containsOnlyNumbers(patchStr) {
		return Version{}, fmt.Errorf(ErrInvalidPatchNumber, patchStr)
	}
	if hasLeadingZeroes(patchStr) {
		return Version{}, fmt.Errorf(ErrLeadingZeroPatch, patchStr)
	}
	patch, err := strconv.ParseUint(patchStr, 10, 64)
	if err != nil {
		return Version{}, err
	}
	v.Patch = patch

	// Parse Prerelease and Build Metadata
	for i < len(version) {
		if version[i] == '-' {
			start = i + 1
			i++
			for i < len(version) && version[i] != '+' {
				i++
			}
			prereleaseStr := version[start:i]

			// Split prerelease by '.'
			parts := strings.Split(prereleaseStr, ".")
			v.PreRelease = make([]PrereleaseVersion, 0, len(parts))
			for _, part := range parts {
				parsedPR, err := NewPrereleaseVersion(part)
				if err != nil {
					return Version{}, err
				}
				v.PreRelease = append(v.PreRelease, parsedPR)
			}
		} else if version[i] == '+' {
			start = i + 1
			i++
			for i < len(version) {
				i++
			}
			buildStr := version[start:i]

			// Split build metadata by '.'
			parts := strings.Split(buildStr, ".")
			v.BuildMetadata = make([]string, 0, len(parts))
			for _, part := range parts {
				if len(part) == 0 {
					return Version{}, ErrEmptyBuildMetadata
				}
				if !containsOnlyAlphanumeric(part) {
					return Version{}, fmt.Errorf(ErrInvalidBuildMetadataChars, part)
				}
				v.BuildMetadata = append(v.BuildMetadata, part)
			}
		} else {
			i++
		}
	}

	return v, nil
}

// String returns the string representation of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha.1+build.123")
//	fmt.Println(v.String()) // Output: 1.2.3-alpha.1+build.123
func (v Version) String() string {
	var sb strings.Builder
	sb.Grow(16) // Preallocate a reasonable capacity

	sb.WriteString(strconv.FormatUint(v.Major, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Minor, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Patch, 10))

	if len(v.PreRelease) > 0 {
		sb.WriteByte('-')
		for i, pr := range v.PreRelease {
			if i > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(pr.String())
		}
	}

	if len(v.BuildMetadata) > 0 {
		sb.WriteByte('+')
		for i, bm := range v.BuildMetadata {
			if i > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(bm)
		}
	}

	return sb.String()
}

// Compare compares two Version instances.
// Returns -1 if v < other, 0 if v == other, +1 if v > other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.4")
//	fmt.Println(v1.Compare(v2)) // Output: -1
func (v Version) Compare(other Version) int {
	// Compare Major
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}

	// Compare Minor
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}

	// Compare Patch
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	// Handle pre-release comparison
	if len(v.PreRelease) == 0 && len(other.PreRelease) == 0 {
		return 0
	} else if len(v.PreRelease) == 0 && len(other.PreRelease) > 0 {
		return 1
	} else if len(v.PreRelease) > 0 && len(other.PreRelease) == 0 {
		return -1
	}

	// Compare pre-release identifiers one by one
	minLen := len(v.PreRelease)
	if len(other.PreRelease) < minLen {
		minLen = len(other.PreRelease)
	}

	for i := 0; i < minLen; i++ {
		comp := v.PreRelease[i].Compare(other.PreRelease[i])
		if comp != 0 {
			return comp
		}
	}

	// If all compared identifiers are equal, the one with more identifiers has higher precedence
	if len(v.PreRelease) < len(other.PreRelease) {
		return -1
	}
	if len(v.PreRelease) > len(other.PreRelease) {
		return 1
	}
	return 0
}

// Equal checks if two versions are equal.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.3")
//	fmt.Println(v1.Equal(v2)) // Output: true
func (v Version) Equal(other Version) bool {
	return v.Compare(other) == 0
}

// LessThan checks if v is less than other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.4")
//	fmt.Println(v1.LessThan(v2)) // Output: true
func (v Version) LessThan(other Version) bool {
	return v.Compare(other) == -1
}

// LessThanOrEqual checks if v is less than or equal to other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.3")
//	fmt.Println(v1.LessThanOrEqual(v2)) // Output: true
func (v Version) LessThanOrEqual(other Version) bool {
	return v.Compare(other) <= 0
}

// GreaterThan checks if v is greater than other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.4")
//	v2 := semver.MustParse("1.2.3")
//	fmt.Println(v1.GreaterThan(v2)) // Output: true
func (v Version) GreaterThan(other Version) bool {
	return v.Compare(other) == 1
}

// GreaterThanOrEqual checks if v is greater than or equal to other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.2")
//	fmt.Println(v1.GreaterThanOrEqual(v2)) // Output: true
func (v Version) GreaterThanOrEqual(other Version) bool {
	return v.Compare(other) >= 0
}

// MustParse is a helper function that parses a version string and panics if invalid.
//
// Example:
//
//	v := semver.MustParse("1.2.3")
//	fmt.Println(v) // Output: 1.2.3
func MustParse(version string) Version {
	v, err := Parse(version)
	if err != nil {
		panic(err)
	}
	return v
}

// MarshalText implements encoding.TextMarshaler.
// It returns the string representation of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha+build.456")
//	text, err := v.MarshalText()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(text)) // Output: 1.2.3-alpha+build.456
func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It parses the given text into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.UnmarshalText([]byte("1.2.3-alpha+build.456"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3-alpha+build.456
func (v *Version) UnmarshalText(text []byte) error {
	parsed, err := Parse(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
// It returns the binary encoding of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha")
//	binaryData, err := v.MarshalBinary()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%s\n", binaryData) // Output: 1.2.3-alpha
func (v Version) MarshalBinary() ([]byte, error) {
	return v.MarshalText()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// It decodes the given binary data into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.UnmarshalBinary([]byte("1.2.3+build.123"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3+build.123
func (v *Version) UnmarshalBinary(data []byte) error {
	return v.UnmarshalText(data)
}

// MarshalJSON implements json.Marshaler.
// It returns the JSON encoding of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-beta")
//	jsonData, err := v.MarshalJSON()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(jsonData)) // Output: "1.2.3-beta"
func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON implements json.Unmarshaler.
// It decodes JSON data into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.UnmarshalJSON([]byte("\"1.2.3-beta+build\""))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3-beta+build
func (v *Version) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}
	return v.UnmarshalText([]byte(text))
}

// Value implements database/sql/driver.Valuer.
// It returns the string representation of the Version as a database value.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha")
//	dbValue, err := v.Value()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(dbValue) // Output: 1.2.3-alpha
func (v Version) Value() (driver.Value, error) {
	return v.String(), nil
}

// Scan implements database/sql.Scanner.
// It scans a database value into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.Scan("1.2.3-alpha+build")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3-alpha+build
func (v *Version) Scan(value interface{}) error {
	switch t := value.(type) {
	case string:
		return v.UnmarshalText([]byte(t))
	case []byte:
		return v.UnmarshalText(t)
	default:
		return fmt.Errorf("unsupported type %T for Version", value)
	}
}