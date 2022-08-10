package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1: parse flags and parameters
	var t time.Duration
	flag.DurationVar(&t, "timeout", 10, "setup tcp connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("error in set of parameters, need 2: host port")
		os.Exit(1)
	}

	addr := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	// 2: subscribe signal
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)

	// 3: create client + establish connection
	client := NewTelnetClient(addr, t, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "error in connection: ", err)
		return
	}
	defer client.Close()
	fmt.Fprintln(os.Stderr, "...Connected to", addr)

	// 4: run writer and reader
	go func() {
		defer cancel()
		err := client.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
	}()

	go func() {
		defer cancel()
		err := client.Send()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Fprintln(os.Stderr, "...EOF")
	}()

	// 5: wait till the end
	<-ctx.Done()
}
