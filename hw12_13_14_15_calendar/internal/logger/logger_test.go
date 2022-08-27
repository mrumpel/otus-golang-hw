package logger

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("stdout clean run", func(t *testing.T) {
		_, err := New("INFO", os.Stdout.Name())
		require.NoError(t, err)
	})

	t.Run("main test", func(t *testing.T) {
		file, err := ioutil.TempFile("", "test")
		require.NoError(t, err)
		defer func() {
			err := file.Close()
			require.NoError(t, err)
			err = os.Remove(file.Name())
			require.NoError(t, err)
		}()

		lg, err := New("INFO", file.Name())
		require.NoError(t, err)

		lg.Info("message1")
		b, err := io.ReadAll(file)
		require.NoError(t, err)
		s := string(b)

		require.Contains(t, s, "message1")
		require.Contains(t, s, "info")

		_, err = file.Seek(0, 0)
		require.NoError(t, err)

		lg.Error("message2")
		b, err = io.ReadAll(file)
		require.NoError(t, err)
		s = string(b)
		require.Contains(t, s, "message2")
		require.Contains(t, s, "error")
	})

	t.Run("wrong file", func(t *testing.T) {
		_, err := New("INFO", "/")
		require.Error(t, err)
	})
}
