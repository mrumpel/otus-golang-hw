package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("no such file", func(t *testing.T) {
		c, err := NewConfig("file_not_exist")
		require.Nil(t, c)
		require.Error(t, err)
	})

	t.Run("default config", func(t *testing.T) {
		c, err := NewConfig("../../" + defaultConfigPath)

		require.NoError(t, err)

		require.Equal(t, "postgres", c.Storage.Type)
		require.Equal(t, "INFO", c.Logger.Level)
		require.Equal(t, "1984", c.Server.Port)
	})
}
