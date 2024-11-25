# semver

[![CI](https://github.com/sixafter/semver/workflows/ci/badge.svg)](https://github.com/sixafter/semver/actions)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/semver)](https://img.shields.io/github/go-mod/go-version/sixafter/semver)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_semver&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_semver)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/semver)](https://github.com/sixafter/semver/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/semver.svg)](https://pkg.go.dev/github.com/sixafter/semver)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/semver)](https://goreportcard.com/report/github.com/sixafter/semver)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
![CodeQL](https://github.com/sixafter/semver/actions/workflows/codeql-analysis.yaml/badge.svg)

A Semantic Versioning 2.0.0 compliant parser and utility library written in Go.

## Features

The `semver` library offers a comprehensive and efficient solution for working with Semantic Versioning 2.0.0. Key features include:

### Zero Dependencies
- Lightweight implementation with no external dependencies beyond the standard library.

### Full Compliance
- Parses and validates semantic versions according to the [Semantic Versioning 2.0.0](https://semver.org) specification.

### Version Parsing and Validation
- Parse semantic version strings into structured `Version` objects.
- Automatically validate version strings for correctness, including:
    - Major, minor, and patch components.
    - Pre-release identifiers (e.g., `alpha`, `beta`, `rc.1`).
    - Build metadata (e.g., `build.123`).
    - Enforces no leading zeroes in numeric components.

### Version Comparison
- Supports comparison of versions using Semantic Versioning rules:
    - `Compare`: Returns -1, 0, or 1 for less than, equal to, or greater than comparisons.
    - Convenient helper methods:
        - `LessThan`, `LessThanOrEqual`
        - `GreaterThan`, `GreaterThanOrEqual`
        - `Equal`
- Correctly handles precedence rules for pre-release versions and build metadata.

### Version Ranges
- Flexible range functionality for evaluating version constraints:
    - Define complex version ranges using a familiar syntax (e.g., `">=1.0.0 <2.0.0"`).
    - Determine whether a version satisfies a given range.
    - Combine multiple ranges for advanced constraints (e.g., `">=1.2.3 || <1.0.0-alpha"`).
- Useful for dependency management, release gating, and compatibility checks.

### Version Construction
- Create `Version` instances programmatically using the `NewVersion` constructor.
- Supports detailed customization of pre-release identifiers and build metadata.

### JSON Support
- Seamlessly marshal and unmarshal `Version` objects to and from JSON.
- Works with `encoding/json` for easy integration with APIs and configuration files.

### Database Support
- Compatible with `database/sql`:
    - Implements `driver.Valuer` to store `Version` in databases.
    - Implements `sql.Scanner` to retrieve `Version` from databases.

### Encoding and Decoding
- Implements standard Go interfaces:
    - `encoding.TextMarshaler` and `encoding.TextUnmarshaler` for text encoding.
    - `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler` for binary encoding.

### Performance Optimizations
- Efficient parsing and comparison with minimal memory allocations.
- Designed for high performance with concurrent workloads.

### Well-Tested
- Comprehensive test coverage, including:
    - Functional tests for all features.
    - Benchmarks to validate performance optimizations.
    - Detailed tests for range evaluation, parsing, and edge cases.

---

## Installation

### Using `go get`

To install the package, run the following command:

```sh
go get -u github.com/sixafter/semver
```

To use the package in your Go project, import it as follows:

```go
import "github.com/sixafter/semver"
```

---

## Usage

Here are some common use cases for the `semver` library:

### Parsing and Validating Versions
Parse a version string into a `Version` object and validate its correctness:

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/sixafter/semver"
)

func main() {
	v, err := semver.Parse("1.2.3-alpha.1+build.123")
	if err != nil {
		log.Fatalf("Error parsing version: %v", err)
	}
	fmt.Println(v) // Output: 1.2.3-alpha.1+build.123
}
```

### Comparing Versions

Compare two versions to determine their order:

```go
package main

import (
	"fmt"

	"github.com/sixafter/semver"
)

func main() {
	v1, _ := semver.Parse("1.2.3")
	v2, _ := semver.Parse("2.0.0")

	fmt.Println(v1.LessThan(v2)) // Output: true
	fmt.Println(v2.GreaterThan(v1)) // Output: true
}
```

### Version Ranges

Evaluate whether a version satisfies a given range:

```go
package main

import (
    "fmt"

    "github.com/sixafter/semver"
)

func main() {
	r, _ := semver.ParseRange(">=1.0.0 <2.0.0")
	v, _ := semver.Parse("1.5.0")

	fmt.Println(r.Contains(v)) // Output: true
}
```

### Working with JSON

Serialize and deserialize versions using JSON:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/sixafter/semver"
)

func main() {
	v := semver.MustParse("1.2.3-alpha.1+build.123")

	data, _ := json.Marshal(v)
	fmt.Println(string(data)) // Output: "1.2.3-alpha.1+build.123"

	var v2 semver.Version
	_ = json.Unmarshal(data, &v2)
	fmt.Println(v2) // Output: 1.2.3-alpha.1+build.123
}
```

---

## Performance

```shell
make bench
```

### Interpreting Results:

Sample output might look like this:

<details>
  <summary>Expand to see results</summary>

```shell
go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out
goos: darwin
goarch: arm64
pkg: github.com/sixafter/semver
cpu: Apple M2 Ultra
BenchmarkParseVersionSerial-24                   1642676               725.2 ns/op           544 B/op         14 allocs/op
BenchmarkParseVersionConcurrent-24               4059534               303.9 ns/op           544 B/op         14 allocs/op
BenchmarkParseVersionAllocations-24              7236553               162.4 ns/op           144 B/op          4 allocs/op
BenchmarkParseVersionLargeSerial-24                  208           5747548 ns/op         4160079 B/op     110000 allocs/op
BenchmarkParseVersionLargeConcurrent-24              656           1874606 ns/op         4160163 B/op     110000 allocs/op
PASS
ok      github.com/sixafter/semver      8.274s
```
</details>

* `ns/op`: Nanoseconds per operation. Lower values indicate faster performance.
* `B/op`: Bytes allocated per operation. Lower values indicate more memory-efficient code.
* `allocs/op`: Number of memory allocations per operation. Fewer allocations generally lead to better performance.

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
