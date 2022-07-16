package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)

	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, item := range dirItems {
		if item.IsDir() {
			continue
		}
		if strings.Contains(item.Name(), "=") {
			return nil, errors.New("File name \"" + item.Name() + "\" contains \"=\"")
		}

		info, err := item.Info()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			envs[item.Name()] = EnvValue{
				NeedRemove: true,
			}
			continue
		}

		file, err := os.Open(dir + "/" + item.Name())
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		envs[item.Name()] = EnvValue{
			Value: evaluateEnvValue(scanner.Bytes()),
		}

	}

	return envs, nil
}

func evaluateEnvValue(content []byte) string {
	content = bytes.Replace(content, []byte{0x00}, []byte("\n"), 1)
	return string(bytes.TrimRight(content, " \t\r"))
}
