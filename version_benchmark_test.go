// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"testing"
)

func BenchmarkParseVersionSerial(b *testing.B) {
	b.ReportAllocs()

	p, err := NewParser(WithStrictAdherence(true))
	if err != nil {
		b.Fatalf("Error creating parser: %v", err)
	}

	versions := []string{
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use modulo to select a different version in each loop
		version := versions[i%len(versions)]
		_, err = p.Parse(version)
		if err != nil {
			b.Errorf("Error parsing version %s: %v", version, err)
		}
	}
}

func BenchmarkParseVersionConcurrent(b *testing.B) {
	b.ReportAllocs()

	p, err := NewParser(WithStrictAdherence(true))
	if err != nil {
		b.Fatalf("Error creating parser: %v", err)
	}

	versions := []string{
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

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Use modulo to select a different version in each loop
			version := versions[i%len(versions)]
			_, err := p.Parse(version)
			if err != nil {
				b.Errorf("Error parsing version %s: %v", version, err)
			}
			i++
		}
	})
}

func BenchmarkParseVersionAllocations(b *testing.B) {
	b.ReportAllocs()
	version := "1.2.3-alpha.1+build.123"

	p, err := NewParser(WithStrictAdherence(true))
	if err != nil {
		b.Errorf("Error creating parser: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Parse(version)
		if err != nil {
			b.Errorf("Error parsing version %s: %v", version, err)
		}
	}
}
