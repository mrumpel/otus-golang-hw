package main

import (
	"log"
	"os"
)

func main() {
	// 1: Check args count
	if len(os.Args) < 2 {
		log.Fatal("Not enough args")
	}

	// 2: parse dir to Envs
	vars, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// 3: run command with os.Envs + Envs from dir, send result to OS
	os.Exit(RunCmd(os.Args[2:], vars))
}
