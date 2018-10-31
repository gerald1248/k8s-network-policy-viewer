package main

import (
	//"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func serve(certificate string, key string, hostname string, selfsigned bool, port int) {}

//Visulize Kubernetes NetworkPolicy objects
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [<JSON/YAML file> [<JSON/YAML file>]]\nAlternatively, pipe input to STDIN: kubectl get networkpolicy,po --all-namespaces -o json | %s\n", filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	certificate := flag.String("c", "cert.pem", "TLS server certificate")
	key := flag.String("k", "key.pem", "TLS server key")
	host := flag.String("n", "localhost", "hostname")
	port := flag.Int("p", 8080, "listen on port")
	selfsigned := flag.Bool("s", false, "Self-signed certificate")
	//output := flag.String("o", "dot", "output format (dot, md)")

	flag.Parse()
	args := flag.Args()

	//use case [A]: STDIN handling
	stdinFileInfo, _ := os.Stdin.Stat()
	if stdinFileInfo.Mode()&os.ModeNamedPipe != 0 {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		Result, err := processBytes(stdin)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		fmt.Println(Result.Buffer)
		os.Exit(0)
	}

	//use case [B]: server
	if len(args) == 0 {
		serve(*certificate, *key, *host, *selfsigned, *port)
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
		buffer, err := processFile(arg)
		secs := time.Since(start).Seconds()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v (%.2fs)\n", arg, err, secs)
			os.Exit(1)
		}
		fmt.Println(buffer)
		os.Exit(len(buffer))
	}
}
