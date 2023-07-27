package sandbox_test

import (
	"io"
	"testing"

	"github.com/stealthrocket/timecraft/internal/assert"
	"github.com/stealthrocket/timecraft/internal/sandbox"
)

func TestRootFS(t *testing.T) {
	t.Run("fstest", func(t *testing.T) {
		testFS(t, sandbox.FS(sandbox.RootFS(sandbox.DirFS("testdata/fstest"))))
	})

	t.Run("ReadFile", func(t *testing.T) {
		tests := []struct {
			scenario string
			path     string
			flags    int
			want     string
			err      error
		}{
			{
				scenario: "regular file in the top-level directory",
				path:     "answer",
				want:     "42\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "regular file in a sub-directory",
				path:     "tmp/message",
				want:     "hello world\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "cannot escape the root directory via relative path",
				path:     "../../tmp/message",
				want:     "hello world\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows symlinks to files in the same directory",
				path:     "symlink-to-answer",
				want:     "42\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows symlinks to files in a sub-directory",
				path:     "symlink-to-message",
				want:     "hello world\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows symlinks to files in a parent directory",
				path:     "tmp/symlink-to-answer",
				want:     "42\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows absolute symlinks to files in the same directory",
				path:     "absolute-symlink-to-answer",
				want:     "42\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows absolute symlinks to files in a sub-directory",
				path:     "absolute-symlink-to-message",
				want:     "hello world\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows absolute symlinks to directories",
				path:     "absolute-symlink-to-tmp/message",
				want:     "hello world\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "does not follow symlinks when O_NOFOLLOW is set",
				path:     "symlink-to-answer",
				flags:    sandbox.O_RDONLY | sandbox.O_NOFOLLOW,
				err:      sandbox.ELOOP,
			},

			{
				scenario: "does not follow dangling symlinks",
				path:     "symlink-to-nowhere",
				flags:    sandbox.O_RDONLY,
				err:      sandbox.ENOENT,
			},

			{
				scenario: "does not follow absolute dangling symlinks",
				path:     "absolute-symlink-to-nowhere",
				flags:    sandbox.O_RDONLY,
				err:      sandbox.ENOENT,
			},

			{
				scenario: "follows relative symlinks to files in the same directory",
				path:     "relative-symlink-to-answer",
				want:     "42\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows relative symlinks to files in a sub-directory",
				path:     "relative-symlink-to-message",
				want:     "hello world\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "follows relative symlinks to a sub-directory",
				path:     "relative-symlink-to-tmp/message",
				want:     "hello world\n",
				flags:    sandbox.O_RDONLY,
			},

			{
				scenario: "does not follow relative dangling symlinks",
				path:     "relative-symlink-to-nowhere",
				flags:    sandbox.O_RDONLY,
				err:      sandbox.ENOENT,
			},
		}

		for _, test := range tests {
			t.Run(test.scenario, func(t *testing.T) {
				rootFS := sandbox.RootFS(sandbox.DirFS("testdata/rootfs"))
				b, err := readFile(rootFS, test.path, test.flags)
				if test.err != nil {
					assert.Error(t, err, test.err)
				} else {
					assert.OK(t, err)
					assert.Equal(t, string(b), test.want)
				}
			})
		}
	})
}

func readFile(fsys sandbox.FileSystem, name string, flags int) ([]byte, error) {
	f, err := fsys.Open(name, flags, 0)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	b := make([]byte, s.Size())
	n, err := io.ReadFull(f, b)
	return b[:n], err
}
