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
		fmt.Fprintf(os.Stderr, `
Usage: %s [<JSON/YAML file> [<JSON/YAML file>]]

Set blacklist by exporting NETWORK_POLICY_VIEWER_BLACKLIST containing a comma-separated list of namespaces

Alternatively, pipe input to STDIN: kubectl get networkpolicy,po --all-namespaces -o json | %s
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	port := flag.Int("p", 8080, "listen on port")
	output := flag.String("o", "dot", "output format (dot, json, yaml)")
	server := flag.Bool("s", false, "launch server")
	name := flag.String("n", "", "context name")

	flag.Parse()
	args := flag.Args()

	// identify context/cluster name
	var context string
	// option #1: name param
	if len(*name) > 0 {
		context = *name
	} else {
		context = os.Getenv("CLUSTER_TESTS_CONTEXT") // option #2: custom context variable
		if len(context) == 0 {
			context = os.Getenv("KUBERNETES_PORT_443_TCP_ADDR") // option #3: IP address
		}
	}

	//use case [A]: server
	if *server == true {
		serve(*port, context)
		return
	}

	//use case [B]: STDIN handling
	stdinFileInfo, _ := os.Stdin.Stat()
	if stdinFileInfo.Mode()&os.ModeNamedPipe != 0 {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		result, _, _, _, err := processBytes(stdin, output)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		fmt.Println(result)
		os.Exit(0)
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
