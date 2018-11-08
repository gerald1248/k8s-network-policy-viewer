package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// visualize Kubernetes NetworkPolicy objects
// accept input from files, stdin, API calls
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [<JSON/YAML file> [<JSON/YAML file>]]\nAlternatively, pipe input to STDIN: kubectl get networkpolicy,po --all-namespaces -o json | %s\n", filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		printPodCount()
		os.Exit(0)
	}

	port := flag.Int("p", 8080, "listen on port")
	output := flag.String("o", "dot", "output format (dot, json, yaml)")

	flag.Parse()
	args := flag.Args()

	//use case [A]: STDIN handling
	stdinFileInfo, _ := os.Stdin.Stat()
	if stdinFileInfo.Mode()&os.ModeNamedPipe != 0 {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		result, err := processBytes(stdin, output)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		fmt.Println(result)
		os.Exit(0)
	}

	//use case [B]: server
	if len(args) == 0 {
		serve(*port)
		return
	} else if len(args) == 1 {
		switch args[0] {
		default:
			break
		}
	}

	// use case [C]: file input
	for _, arg := range args {
		start := time.Now()
		buffer, err := processFile(arg, output)
		secs := time.Since(start).Seconds()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v (%.2fs)\n", arg, err, secs)
			os.Exit(1)
		}
		fmt.Println(buffer)
		os.Exit(len(buffer))
	}
}
