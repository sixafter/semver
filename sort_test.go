// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortVersions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	versionStrings := []string{
		"1.0.0",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
		"2.0.0",
		"1.0.1",
		"1.1.0",
	}

	var versions []*Version
	for _, vs := range versionStrings {
		v := MustParse(vs)
		versions = append(versions, &v)
	}

	Sort(versions)

	expectedOrder := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
		"1.0.0",
		"1.0.1",
		"1.1.0",
		"2.0.0",
	}

	for i, v := range versions {
		is.Equal(expectedOrder[i], v.String(), "Versions should be sorted correctly")
	}
}

func TestReverseSortVersions(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	versionStrings := []string{
		"1.0.0",
		"1.0.1",
		"1.1.0",
		"2.0.0",
		"1.0.0-alpha",
		"1.0.0-beta",
	}

	var versions []*Version
	for _, vs := range versionStrings {
		v := MustParse(vs)
		versions = append(versions, &v)
	}

	Sort(versions)
	Reverse(versions)

	expectedOrder := []string{
		"2.0.0",
		"1.1.0",
		"1.0.1",
		"1.0.0",
		"1.0.0-beta",
		"1.0.0-alpha",
	}

	for i, v := range versions {
		is.Equal(expectedOrder[i], v.String(), "Versions should be reverse sorted correctly")
	}
}
