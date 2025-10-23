// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersion(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// Define test cases
	testCases := []struct {
		name          string
		major         uint64
		minor         uint64
		patch         uint64
		preRelease    []PrereleaseVersion
		buildMetadata []string
		expected      string
	}{
		{
			name:          "Version with prerelease and build metadata",
			major:         1,
			minor:         2,
			patch:         3,
			preRelease:    []PrereleaseVersion{{partString: "alpha", isNumeric: false}, {partNumeric: 1, isNumeric: true}},
			buildMetadata: []string{"build", "2024"},
			expected:      "1.2.3-alpha.1+build.2024",
		},
		{
			name:          "Version without prerelease and build metadata",
			major:         1,
			minor:         0,
			patch:         0,
			preRelease:    nil,
			buildMetadata: nil,
			expected:      "1.0.0",
		},
		{
			name:          "Version with only prerelease",
			major:         2,
			minor:         1,
			patch:         4,
			preRelease:    []PrereleaseVersion{{partString: "beta", isNumeric: false}},
			buildMetadata: nil,
			expected:      "2.1.4-beta",
		},
		{
			name:          "Version with only build metadata",
			major:         0,
			minor:         9,
			patch:         7,
			preRelease:    nil,
			buildMetadata: []string{"build123"},
			expected:      "0.9.7+build123",
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := New(tc.major, tc.minor, tc.patch, tc.preRelease, tc.buildMetadata)
			is.Equal(tc.expected, v.String(), "Version string should match expected value")
		})
	}
}

func TestSupportedVersion(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v := Version{
		Major: 2,
		Minor: 0,
		Patch: 0,
	}
	is.True(SupportedVersion.Equal(v), "Version should be supported")
}

func TestParseVersion(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		input       string
		expected    Version
		shouldError bool
	}{
		{
			input: "1.2.3",
			expected: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			shouldError: false,
		},
		{
			input: "0.1.0",
			expected: Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			shouldError: false,
		},
		{
			input: "1.2.3-alpha.1+build.123",
			expected: Version{
				Major:         1,
				Minor:         2,
				Patch:         3,
				PreRelease:    []PrereleaseVersion{{partString: "alpha", partNumeric: 0, isNumeric: false}, {partString: "", partNumeric: 1, isNumeric: true}},
				BuildMetadata: []string{"build", "123"},
			},
			shouldError: false,
		},
		{
			input:       "invalid.version",
			expected:    Version{},
			shouldError: true,
		},
		{
			input:       "1.2",
			expected:    Version{},
			shouldError: true,
		},
		{
			input:       "",
			expected:    Version{},
			shouldError: true,
		},
	}

	for _, test := range tests {
		v, err := Parse(test.input)
		if test.shouldError {
			is.Error(err, "Expected an error for input: %s", test.input)
		} else {
			is.NoError(err, "Did not expect an error for input: %s", test.input)
			is.Equal(test.expected, v, "Versions should match")
		}
	}
}

func TestVersionString(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v := Version{
		Major:         1,
		Minor:         2,
		Patch:         3,
		PreRelease:    []PrereleaseVersion{{partString: "alpha", partNumeric: 0, isNumeric: false}, {partString: "", partNumeric: 1, isNumeric: true}},
		BuildMetadata: []string{"build", "123"},
	}
	expected := "1.2.3-alpha.1+build.123"
	is.Equal(expected, v.String(), "Version string should match expected value")
}

func TestVersionComparison(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.1.0", "1.0.1", 1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.0-alpha", "1.0.0-beta", -1},
		{"1.0.0-alpha.1", "1.0.0-alpha", 1},
		{"1.0.0+build1", "1.0.0+build2", 0}, // Build metadata is ignored in comparison
		{"1.0.0+build1", "1.0.0", 0},        // Build metadata is ignored in comparison
	}

	for _, test := range tests {
		v1 := MustParse(test.v1)
		v2 := MustParse(test.v2)
		result := v1.Compare(v2)
		is.Equal(test.expected, result, "Comparison between %s and %s", test.v1, test.v2)
	}
}

func TestMustParse(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.NotPanics(func() {
		v := MustParse("1.2.3")
		is.NotNil(v)
	}, "MustParse should not panic on valid version")

	is.Panics(func() {
		MustParse("invalid")
	}, "MustParse should panic on invalid version")
}

func TestVersionEqual(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v1 := MustParse("1.2.3")
	v2 := MustParse("1.2.3")
	v3 := MustParse("1.2.4")

	is.True(v1.Equal(v2), "Versions should be equal")
	is.False(v1.Equal(v3), "Versions should not be equal")
}

func TestVersionLessThan(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v1 := MustParse("1.2.3")
	v2 := MustParse("1.2.4")

	is.True(v1.LessThan(v2), "%s should be less than %s", v1, v2)
	is.False(v2.LessThan(v1), "%s should not be less than %s", v2, v1)
}

func TestVersionGreaterThan(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v1 := MustParse("2.0.0")
	v2 := MustParse("1.9.9")

	is.True(v1.GreaterThan(v2), "%s should be greater than %s", v1, v2)
	is.False(v2.GreaterThan(v1), "%s should not be greater than %s", v2, v1)
}

func TestVersionPreReleaseComparison(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v1 := MustParse("1.0.0-alpha")
	v2 := MustParse("1.0.0-alpha.1")
	v3 := MustParse("1.0.0-alpha.beta")
	v4 := MustParse("1.0.0-beta")
	v5 := MustParse("1.0.0-beta.2")
	v6 := MustParse("1.0.0-beta.11")
	v7 := MustParse("1.0.0-rc.1")
	v8 := MustParse("1.0.0")

	versions := []Version{v1, v2, v3, v4, v5, v6, v7, v8}
	for i := 0; i < len(versions)-1; i++ {
		is.True(versions[i].LessThan(versions[i+1]), "%s should be less than %s", versions[i], versions[i+1])
	}
}

func TestParseInvalidVersions(t *testing.T) {
	// Define a list of invalid versions to test.
	invalidVersions := []string{
		"1.0",             // Incomplete version (missing patch version)
		"v1.0.0",          // Prefix with `v` is not allowed in Semver
		"1.0.0-alpha..1",  // Double dots are invalid
		"1.0.0+build+123", // Invalid multiple `+` in build metadata
		"1.0.0-01",        // Leading zeros in numeric pre-release identifiers are not allowed
		"1.0.0-",          // Ends with a dash
		"1.0.0+build.!",   // Invalid character `!` in build metadata
		"1.0.0-beta_$",    // Invalid character `$` in pre-release identifier
		"1..0.0",          // Double dots in the version components
	}

	for _, version := range invalidVersions {
		t.Run(version, func(t *testing.T) {
			_, err := Parse(version)
			if err == nil {
				t.Errorf("expected error for invalid version: %s, but got none", version)
			}
		})
	}
}

func TestStrictAdherence(t *testing.T) {
	strictParser, err := NewParser(WithStrictAdherence(true))
	if err != nil {
		t.Fatalf("Failed to create strict parser: %v", err)
	}

	nonStrictParser, err := NewParser(WithStrictAdherence(false))
	if err != nil {
		t.Fatalf("Failed to create non-strict parser: %v", err)
	}

	// Test cases for strict adherence
	tests := []struct {
		parser      Parser
		input       string
		expectError bool
	}{
		// Strict adherence enabled - should fail for leading zeros
		{parser: strictParser, input: "1.01.0", expectError: true},
		{parser: strictParser, input: "1.0.00", expectError: true},
		{parser: strictParser, input: "01.0.0", expectError: true},

		// Strict adherence enabled - should succeed for valid versions
		{parser: strictParser, input: "1.0.0", expectError: false},
		{parser: strictParser, input: "1.2.3-alpha.1", expectError: false},

		// Non-strict adherence - should allow leading zeros
		{parser: nonStrictParser, input: "1.01.0", expectError: false},
		{parser: nonStrictParser, input: "1.0.00", expectError: false},
		{parser: nonStrictParser, input: "01.0.0", expectError: false},

		// Non-strict adherence - should succeed for valid versions
		{parser: nonStrictParser, input: "1.0.0", expectError: false},
		{parser: nonStrictParser, input: "1.2.3-alpha.1", expectError: false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := tt.parser.Parse(tt.input)
			if tt.parser == nonStrictParser && errors.Is(err, ErrLeadingZeroInNumericIdentifier) {
				err = nil // Allow leading zeros in non-strict mode
			}
			if (err != nil) != tt.expectError {
				t.Errorf("Parse(%s) strict=%v: expected error: %v, got: %v", tt.input, tt.parser == strictParser, tt.expectError, err)
			}
		})
	}
}

func TestInitPanicsOnParserFailure(t *testing.T) {
	// No t.Parallel(): this test mutates package-level state.

	// Save and restore the original indirection.
	orig := newParserFunc
	defer func() {
		newParserFunc = orig
		// Restore DefaultParser after test so other tests don’t crash
		var err error
		DefaultParser, err = newParserFunc()
		if err != nil {
			t.Fatalf("failed to restore DefaultParser: %v", err)
		}
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic during initDefaultParser, got none")
		}
	}()

	// Force failure to simulate init panic
	newParserFunc = func(...Option) (Parser, error) {
		return nil, errors.New("forced failure for test")
	}

	initDefaultParser() // Should panic
}
