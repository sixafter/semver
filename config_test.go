// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetConfig tests the Config() method of the generator.
func TestGetConfig(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	gen, err := NewParser(WithStrictAdherence(true))
	is.NoError(err, "NewGenerator() should not return an error with the default alphabet")

	// Assert that generator implements Configuration interface
	config, ok := gen.(Configuration)
	is.True(ok, "Parser should implement Configuration interface")

	runtimeConfig := config.Config()

	is.True(runtimeConfig.StrictAdherence(), "Config.StrictAdherence should be true")
}
