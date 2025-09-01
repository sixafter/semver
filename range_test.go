// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRange(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		input       string
		shouldError bool
	}{
		// Existing Tests
		{input: ">1.0.0", shouldError: false},
		{input: "<=2.0.0", shouldError: false},
		{input: ">=1.2.3 <2.0.0 || >=3.0.0", shouldError: false},
		{input: "1.0.0", shouldError: false},
		{input: "!=1.0.0", shouldError: false},
		{input: "invalid", shouldError: true},

		// New Tests with Pre-release and Build Metadata
		{input: ">1.0.0-alpha", shouldError: false},                   // Greater than a pre-release version
		{input: "<=2.0.0-beta.1", shouldError: false},                 // Less than or equal to a beta pre-release
		{input: ">=1.0.0-alpha <2.0.0", shouldError: false},           // Range involving pre-release version
		{input: ">=1.2.3+build.123", shouldError: false},              // Version with build metadata
		{input: "1.0.0+build.1", shouldError: false},                  // Exact match with build metadata
		{input: "!=1.0.0-alpha", shouldError: false},                  // Not equal to a pre-release version
		{input: ">=1.0.0-alpha.1 <1.0.0-alpha.3", shouldError: false}, // Range involving pre-release identifiers
		{input: "<1.0.0+build.1", shouldError: false},                 // Less than a version with build metadata
		{input: "1.0.0-beta+exp.sha.5114f85", shouldError: false},     // Specific version with pre-release and build metadata
		{input: ">2.1.0-rc.1 <2.1.0+build.789", shouldError: false},   // Range between a release candidate and a version with build metadata
		{input: "1.0.0-alpha+build-metadata", shouldError: false},     // Pre-release version with build metadata
		{input: ">1.0.0-invalid", shouldError: false},                 // Valid pre-release identifier
		{input: "<=1.0.0+build...123", shouldError: true},             // Invalid build metadata
		{input: ">=1.2.3-pre-release <3.0.0", shouldError: false},     // Range with complex pre-release identifier
		{input: "1.0.0-beta.5+build-xyz", shouldError: false},         // Exact version with pre-release and build metadata
	}

	for _, test := range tests {
		_, err := ParseRange(test.input)
		if test.shouldError {
			is.Error(err, "Expected error for input: %s", test.input)
		} else {
			is.NoError(err, "Did not expect error for input: %s", test.input)
		}
	}
}

func TestVersionRangeMatches(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		rangeStr    string
		version     string
		shouldMatch bool
	}{
		{rangeStr: ">1.0.0", version: "1.0.1", shouldMatch: true},
		{rangeStr: "<=2.0.0", version: "2.0.0", shouldMatch: true},
		{rangeStr: ">=1.2.3 <2.0.0", version: "1.2.3", shouldMatch: true},
		{rangeStr: ">=1.2.3 <2.0.0", version: "2.0.0", shouldMatch: false},
		{rangeStr: ">=1.2.3 <2.0.0 || >=3.0.0", version: "3.1.0", shouldMatch: true},
		{rangeStr: "!=1.0.0", version: "1.0.0", shouldMatch: false},
	}

	for _, test := range tests {
		rng, err := ParseRange(test.rangeStr)
		is.NoError(err, "Range should parse correctly")
		v := MustParse(test.version)
		matches := rng.Contains(v)
		if test.shouldMatch {
			is.True(matches, "Version %s should match range %s", test.version, test.rangeStr)
		} else {
			is.False(matches, "Version %s should not match range %s", test.version, test.rangeStr)
		}
	}
}

func TestVersionRangeOR(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1, err := ParseRange(">1.0.0 <2.0.0")
	is.NoError(err)
	r2, err := ParseRange(">=3.0.0 !=3.5.0")
	is.NoError(err)
	combinedRange := r1.OR(r2)

	tests := []struct {
		version     string
		shouldMatch bool
	}{
		{"1.5.0", true},
		{"2.5.0", false},
		{"3.0.0", true},
		{"3.5.0", false},
		{"4.0.0", true},
	}

	for _, test := range tests {
		v := MustParse(test.version)
		matches := combinedRange.Contains(v)
		if test.shouldMatch {
			is.True(matches, "Version %s should match combined range", test.version)
		} else {
			is.False(matches, "Version %s should not match combined range", test.version)
		}
	}
}

func TestVersionRangeAND(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1, err := ParseRange(">1.0.0")
	is.NoError(err)
	r2, err := ParseRange("<2.0.0")
	is.NoError(err)
	combinedRange := r1.AND(r2)

	tests := []struct {
		version     string
		shouldMatch bool
	}{
		{"1.5.0", true},
		{"2.0.0", false},
		{"1.0.0", false},
	}

	for _, test := range tests {
		v := MustParse(test.version)
		matches := combinedRange.Contains(v)
		if test.shouldMatch {
			is.True(matches, "Version %s should match combined range", test.version)
		} else {
			is.False(matches, "Version %s should not match combined range", test.version)
		}
	}
}

func TestMustParseRange(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	t.Run("ValidRange", func(t *testing.T) {

		r := MustParseRange(">1.0.0 <2.0.0")
		v := MustParse("1.5.0")

		is.NotNil(r, "Expected a non-nil VersionRange")
		is.True(r.Contains(v), "Expected version to be within the range")
	})

	t.Run("InvalidRange", func(t *testing.T) {

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Expected panic for invalid range")
			}
		}()

		// This should panic
		MustParseRange("invalid range")
	})
}
