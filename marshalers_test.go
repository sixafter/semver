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
	is.EqualError(err, "unsupported type for Version")
}
