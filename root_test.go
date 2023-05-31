package main_test

import (
	"testing"

	"github.com/stealthrocket/timecraft/internal/assert"
)

var root = tests{
	"invoking timecraft without a command prints the introduction message": func(t *testing.T) {
		stdout, stderr, err := timecraft(t)
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "timecraft - WebAssembly Time Machine\n")
		assert.Equal(t, stderr, "")
	},

	"show the timecraft help with the short option": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "-h")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft <command> ")
		assert.Equal(t, stderr, "")
	},

	"show the timecraft help with the long option": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "--help")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft <command> ")
		assert.Equal(t, stderr, "")
	},
}
