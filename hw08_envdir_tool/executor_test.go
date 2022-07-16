package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name         string
		cmd          []string
		expectedCode int
	}{
		{
			name:         "successful run",
			cmd:          []string{"pwd"},
			expectedCode: 0,
		},
		{
			name:         "general error",
			cmd:          []string{"cat", "file_dat_not_exist"},
			expectedCode: 1,
		},
		{
			name:         "command not exist",
			cmd:          []string{"/bin/bash", "command_not_exist"},
			expectedCode: 127,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			x := RunCmd(test.cmd, nil)
			require.Equal(t, test.expectedCode, x)
		})
	}
}

func TestPrepareEnvs(t *testing.T) {

	t.Run("no deletions", func(t *testing.T) {
		ed := make(Environment)
		envs := []string{"one=001", "two=002"}

		res := prepareEnvs(envs, ed)

		require.Len(t, res, 2)
		require.True(t, strInSlice(res, "one=001"))
		require.True(t, strInSlice(res, "two=002"))
	})

	t.Run("only delete", func(t *testing.T) {
		ed := make(Environment)
		ed["two"] = EnvValue{NeedRemove: true}
		envs := []string{"one=001", "two=002"}

		res := prepareEnvs(envs, ed)

		require.Len(t, res, 1)
		require.True(t, strInSlice(res, "one=001"))
		require.False(t, strInSlice(res, "two=002"))
	})

	t.Run("only add", func(t *testing.T) {
		ed := make(Environment)
		ed["two"] = EnvValue{NeedRemove: true}
		ed["three"] = EnvValue{Value: "003"}
		envs := []string{"one=001", "two=002"}

		res := prepareEnvs(envs, ed)

		require.Len(t, res, 2)
		require.True(t, strInSlice(res, "one=001"))
		require.False(t, strInSlice(res, "two=002"))
		require.True(t, strInSlice(res, "three=003"))
	})

	t.Run("add and delete", func(t *testing.T) {
		ed := make(Environment)
		ed["three"] = EnvValue{Value: "003"}
		envs := []string{"one=001", "two=002"}

		res := prepareEnvs(envs, ed)

		require.Len(t, res, 3)
		require.True(t, strInSlice(res, "one=001"))
		require.True(t, strInSlice(res, "two=002"))
		require.True(t, strInSlice(res, "three=003"))
	})
}

func strInSlice(slice []string, str string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}
	return false
}
