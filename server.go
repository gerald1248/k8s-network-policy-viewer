package main

import (
        "bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
        "os/exec"
        "strings"
)

type PostStruct struct {
	Buffer string
}

func serve(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePost(&w, r)
	case "GET":
		handleGet(&w, r)
	}
}

func handleGet(w *http.ResponseWriter, r *http.Request) {
	var buffer, dot, output string

	getJsonData(&buffer)

	output = "dot"

	dot, err := processBytes([]byte(buffer), &output)
	if err != nil {
		dot = err.Error()
	}

	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dot)
	var svg bytes.Buffer
	cmd.Stdout = &svg
	err = cmd.Run()
	if err != nil {
		log.Printf("Graphviz conversion failed: %s\n", err)
	}
	fmt.Fprintf(*w, "GET request\nRequest struct = %v\n\nJSON data:\n%s\nDot:\n%s\nSVG:\n%s\n", r, buffer, dot, svg.String())
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body '%s': %s", body, err)
		return
	}
}
