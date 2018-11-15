package main

import (
	"bytes"
	"fmt"
	"html/template"
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
		fmt.Fprintf(*w, "Can't process input data: %s\n", err)
		return
	}

	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dot)
	var svg bytes.Buffer
	cmd.Stdout = &svg
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(*w, "Graphviz conversion failed: %s\n", err)
		return
	}

	isolatedPercentage := 100 - percentage
	colorClass := "progress-bar-success"
	if isolatedPercentage < 50 {
		colorClass = "progress-bar-danger"
	} else if isolatedPercentage < 75 {
		colorClass = "progress-bar-warning"
	}

	fields := map[string]string{
		"Svg": strings.Replace(svg.String(), "Times,serif", "sans-serif", -1),
		"ColorClass": colorClass,
		"Percentage": string(isolatedPercentage),
	}

	t := template.Must(template.New("index.tmpl").ParseFiles("tmpl/index.tmpl"))
	err = t.Execute(*w, fields)
	if err != nil {
		fmt.Fprintf(*w, "Can't apply template: %s\n", err)
		return
	}
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body '%s': %s", body, err)
		return
	}
}

func check(err error) {
	if err != nil {
		log.Printf("%s", err)
	}
}
