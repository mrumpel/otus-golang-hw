package main

import "testing"
import "github.com/stretchr/testify/require"

func TestReadDir(t *testing.T) {

	t.Run("no directory", func(t *testing.T) {
		env, err := ReadDir("dir_does_not_exist")

		require.Error(t, err)
		require.Nil(t, env)
	})

	t.Run("empty directory", func(t *testing.T) {
		env, err := ReadDir("testdata/emptydir")

		require.Nil(t, err)
		require.Len(t, env, 0)
	})

	t.Run("dir inside", func(t *testing.T) {
		env, err := ReadDir("testdata/dirinside")

		require.Nil(t, err)
		require.Len(t, env, 1, "nested dir not skipped")
	})

	t.Run("read simple files", func(t *testing.T) {
		env, err := ReadDir("testdata/simple")

		require.Nil(t, err)
		require.Len(t, env, 3)
		require.True(t, env["todelete"].NeedRemove, "delete flag reading fail")
		require.Equal(t, "001", env["one"].Value, "one string reading fail")
		require.Equal(t, "002", env["two"].Value, "multiple string reading fail")
	})
}

func TestEvaluateEnvValue(t *testing.T) {
	tests := []struct {
		name, in, out string
	}{
		{
			name: "regular string",
			in:   "/usr/local/go",
			out:  "/usr/local/go",
		},
		{
			name: "0x00 replacement",
			in:   "fullpath\x00notincluded",
			out:  "fullpath\nnotincluded",
		},
		{
			name: "tabs and spaces",
			in:   "/usr/local/go\t\t \r\r  \t",
			out:  "/usr/local/go",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := evaluateEnvValue([]byte(test.in))
			require.Equal(t, test.out, res)
		})
	}

}
