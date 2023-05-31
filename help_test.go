package main_test

import (
	"testing"

	"github.com/stealthrocket/timecraft/internal/assert"
)

var help = tests{
	"calling help with an unknown command causes an error": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "whatever")
		assert.ExitError(t, err, 2)
		assert.Equal(t, stdout, "")
		assert.Equal(t, stderr, "timecraft help whatever: unknown command\n")
	},

	"passing an unsupported flag to the command causes an error": func(t *testing.T) {
		_, _, err := timecraft(t, "help", "-_")
		assert.ExitError(t, err, 2)
	},

	"show the help command help with the short option": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "-h")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft <command> ")
		assert.Equal(t, stderr, "")
	},

	"show the help command help with the long option": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "--help")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft <command> ")
		assert.Equal(t, stderr, "")
	},

	"show the help command help after a command name": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "get", "--help")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft <command> ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help config": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "config")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft config ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help describe": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "describe")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft describe ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help export": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "export")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft export ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help get": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "get")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft get ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help help": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "help")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft <command> ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help profile": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "profile")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft profile ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help run": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "run")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft run ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help replay": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "replay")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft replay ")
		assert.Equal(t, stderr, "")
	},

	"timecraft help version": func(t *testing.T) {
		stdout, stderr, err := timecraft(t, "help", "version")
		assert.OK(t, err)
		assert.HasPrefix(t, stdout, "Usage:\ttimecraft version\n")
		assert.Equal(t, stderr, "")
	},
}
