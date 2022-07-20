package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// 1: prepare args and exec
	var args []string

	if len(cmd) > 1 {
		args = cmd[1:]
	}
	comma := cmd[0]

	e := exec.Command(comma, args...)
	e.Stdout = os.Stdout
	e.Stdin = os.Stdin
	e.Stderr = os.Stderr

	// 2: edit envs
	e.Env = prepareEnvs(os.Environ(), env)

	// 3: run and exit code
	err := e.Start()
	if err != nil {
		return e.ProcessState.ExitCode()
	}

	err = e.Wait()
	if err != nil {
		return e.ProcessState.ExitCode()
	}

	return 0
}

func prepareEnvs(globalEnvs []string, editions Environment) []string {
	envs := make(map[string]string)

	for i := range globalEnvs {
		row := strings.Split(globalEnvs[i], "=")
		envs[row[0]] = row[1]
	}

	for k, v := range editions {
		if v.NeedRemove {
			delete(envs, k)
			continue
		}
		envs[k] = v.Value
	}

	res := make([]string, 0)
	for k, v := range envs {
		res = append(res, k+"="+v)
	}

	return res
}
