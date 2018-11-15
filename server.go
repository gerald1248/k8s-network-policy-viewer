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

	dot, percentage, err := processBytes([]byte(buffer), &output)
	if err != nil {
		sData := fmt.Sprintf("<p>Can't process input data: %s</p>", err)
		fmt.Fprintf(*w, page(sData))
		return
	}

	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dot)
	var svg bytes.Buffer
	cmd.Stdout = &svg
	err = cmd.Run()
	if err != nil {
		sDot := fmt.Sprintf("<p>Graphviz conversion failed: %s</p>", err)
		fmt.Fprintf(*w, page(sDot))
		return
	}

	isolatedPercentage := 100 - percentage
	colorClass := "progress-bar-success"
	if isolatedPercentage < 50 {
		colorClass = "progress-bar-danger"
	} else if isolatedPercentage < 75 {
		colorClass = "progress-bar-warning"
	}

	buffer = fmt.Sprintf(`
<div>%s</div>
<div class="progress">
  <div class="progress-bar %s" style="width: %d%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%% isolated</div>
</div>`,
		strings.Replace(svg.String(), "Times,serif", "sans-serif", -1),
		colorClass,
		isolatedPercentage,
		isolatedPercentage,
		isolatedPercentage)
	fmt.Fprintf(*w, page(buffer))
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body '%s': %s", body, err)
		return
	}
}
