// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"encoding/json"
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
			v := NewVersion(tc.major, tc.minor, tc.patch, tc.preRelease, tc.buildMetadata)
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

func TestVersionMarshalText(t *testing.T) {
	t.Parallel()
	v := MustParse("1.2.3-alpha+build.123")
	text, err := v.MarshalText()

	is := assert.New(t)
	is.NoError(err)
	is.Equal("1.2.3-alpha+build.123", string(text))
}

func TestVersionUnmarshalText(t *testing.T) {
	t.Parallel()
	var v Version
	err := v.UnmarshalText([]byte("1.2.3-alpha+build.123"))

	is := assert.New(t)
	is.NoError(err)
	is.Equal(MustParse("1.2.3-alpha+build.123"), v)
}

func TestVersionMarshalBinary(t *testing.T) {
	t.Parallel()
	v := MustParse("1.2.3-beta")
	binaryData, err := v.MarshalBinary()

	is := assert.New(t)
	is.NoError(err)
	is.Equal([]byte("1.2.3-beta"), binaryData)
}

func TestVersionUnmarshalBinary(t *testing.T) {
	t.Parallel()
	var v Version
	err := v.UnmarshalBinary([]byte("1.2.3+build.456"))

	is := assert.New(t)
	is.NoError(err)
	is.Equal(MustParse("1.2.3+build.456"), v)
}

func TestVersionMarshalJSON(t *testing.T) {
	t.Parallel()
	v := MustParse("1.2.3-alpha")
	jsonData, err := json.Marshal(v)

	is := assert.New(t)
	is.NoError(err)
	is.JSONEq(`"1.2.3-alpha"`, string(jsonData))
}

func TestVersionUnmarshalJSON(t *testing.T) {
	t.Parallel()
	var v Version
	err := json.Unmarshal([]byte(`"1.2.3-beta+build.789"`), &v)

	is := assert.New(t)
	is.NoError(err)
	is.Equal(MustParse("1.2.3-beta+build.789"), v)
}

func TestVersionValue(t *testing.T) {
	t.Parallel()
	v := MustParse("1.2.3-alpha")
	dbValue, err := v.Value()

	is := assert.New(t)
	is.NoError(err)
	is.Equal("1.2.3-alpha", dbValue)
}

func TestVersionScan(t *testing.T) {
	t.Parallel()
	var v Version
	is := assert.New(t)

	// Test with string
	err := v.Scan("1.2.3-alpha+build.123")
	is.NoError(err)
	is.Equal(MustParse("1.2.3-alpha+build.123"), v)

	// Test with []byte
	err = v.Scan([]byte("1.2.3-beta"))
	is.NoError(err)
	is.Equal(MustParse("1.2.3-beta"), v)

	// Test with unsupported type
	err = v.Scan(123)
	is.Error(err)
	is.EqualError(err, "unsupported type int for Version")
}
