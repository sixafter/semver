// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"testing"
)

func BenchmarkParseVersionSerial(b *testing.B) {
	versions := []string{
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
	}

	// Preallocate slice to avoid allocation during benchmarking
	parsers := make([]Version, len(versions))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, v := range versions {
			var err error
			parsers[j], err = Parse(v)
			if err != nil {
				b.Errorf("Error parsing version %s: %v", v, err)
			}
		}
	}
}

func BenchmarkParseVersionConcurrent(b *testing.B) {
	versions := []string{
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
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, v := range versions {
				_, err := Parse(v)
				if err != nil {
					b.Errorf("Error parsing version %s: %v", v, err)
				}
			}
		}
	})
}

func BenchmarkParseVersionAllocations(b *testing.B) {
	version := "1.2.3-alpha.1+build.123"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(version)
		if err != nil {
			b.Errorf("Error parsing version %s: %v", version, err)
		}
	}
}

func BenchmarkParseVersionLargeSerial(b *testing.B) {
	// Generate a large number of versions
	var versions []string
	baseVersions := []string{
		"1.0.0",
		"1.2.3-alpha+build.123",
		"2.0.0-beta.1",
		"3.0.0-rc.1",
		"1.0.0+build.1",
		"1.0.0-alpha.beta",
	}
	for i := 0; i < 10000; i++ {
		versions = append(versions, baseVersions...)
	}

	parsers := make([]Version, len(versions))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, v := range versions {
			var err error
			parsers[j], err = Parse(v)
			if err != nil {
				b.Errorf("Error parsing version %s: %v", v, err)
			}
		}
	}
}

func BenchmarkParseVersionLargeConcurrent(b *testing.B) {
	// Generate a large number of versions
	var versions []string
	baseVersions := []string{
		"1.0.0",
		"1.2.3-alpha+build.123",
		"2.0.0-beta.1",
		"3.0.0-rc.1",
		"1.0.0+build.1",
		"1.0.0-alpha.beta",
	}
	for i := 0; i < 10000; i++ {
		versions = append(versions, baseVersions...)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, v := range versions {
				_, err := Parse(v)
				if err != nil {
					b.Errorf("Error parsing version %s: %v", v, err)
				}
			}
		}
	})
}
