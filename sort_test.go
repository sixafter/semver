// Copyright (c) 2024-2025 Six After, Inc
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
		"1.2.3-alpha+build.123",
		"2.0.0-beta.1",
		"3.0.0-rc.1",
		"1.0.0+build.1",
		"1.0.0-alpha.beta",
		"2.1.3",
		"0.1.0",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"2.0.1-beta.2",
		"4.0.0-alpha.3+exp.sha.5114f85",
		"5.1.0+build.5678",
		"3.3.3-rc.2",
		"6.2.0-beta+ci.789",
		"1.1.1-alpha.2.3",
		"7.0.0+build.1234",
		"8.0.0-alpha.1.5+meta.data.001",
		"2.4.5+build.meta.sha256",
		"9.1.2-beta-unstable",
	}

	var versions []*Version
	for _, vs := range versionStrings {
		v := MustParse(vs)
		versions = append(versions, &v)
	}

	Sort(versions)

	expectedOrder := []string{
		"0.1.0",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0",
		"1.0.0+build.1",
		"1.1.1-alpha.2.3",
		"1.2.3-alpha+build.123",
		"2.0.0-beta.1",
		"2.0.1-beta.2",
		"2.1.3",
		"2.4.5+build.meta.sha256",
		"3.0.0-rc.1",
		"3.3.3-rc.2",
		"4.0.0-alpha.3+exp.sha.5114f85",
		"5.1.0+build.5678",
		"6.2.0-beta+ci.789",
		"7.0.0+build.1234",
		"8.0.0-alpha.1.5+meta.data.001",
		"9.1.2-beta-unstable",
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
		"1.2.3-alpha+build.123",
		"2.0.0-beta.1",
		"3.0.0-rc.1",
		"1.0.0+build.1",
		"1.0.0-alpha.beta",
		"2.1.3",
		"0.1.0",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"2.0.1-beta.2",
		"4.0.0-alpha.3+exp.sha.5114f85",
		"5.1.0+build.5678",
		"3.3.3-rc.2",
		"6.2.0-beta+ci.789",
		"1.1.1-alpha.2.3",
		"7.0.0+build.1234",
		"8.0.0-alpha.1.5+meta.data.001",
		"2.4.5+build.meta.sha256",
		"9.1.2-beta-unstable",
	}

	var versions []*Version
	for _, vs := range versionStrings {
		v := MustParse(vs)
		versions = append(versions, &v)
	}

	Sort(versions)
	Reverse(versions)

	expectedOrder := []string{
		"9.1.2-beta-unstable",
		"8.0.0-alpha.1.5+meta.data.001",
		"7.0.0+build.1234",
		"6.2.0-beta+ci.789",
		"5.1.0+build.5678",
		"4.0.0-alpha.3+exp.sha.5114f85",
		"3.3.3-rc.2",
		"3.0.0-rc.1",
		"2.4.5+build.meta.sha256",
		"2.1.3",
		"2.0.1-beta.2",
		"2.0.0-beta.1",
		"1.2.3-alpha+build.123",
		"1.1.1-alpha.2.3",
		"1.0.0+build.1",
		"1.0.0",
		"1.0.0-alpha.beta",
		"1.0.0-alpha.1",
		"1.0.0-alpha",
		"0.1.0",
	}

	for i, v := range versions {
		is.Equal(expectedOrder[i], v.String(), "Versions should be reverse sorted correctly")
	}
}
