package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

// PostStruct wraps a buffer string
type PostStruct struct {
	Buffer string
}

func serve(port int, context string) {
	http.Handle("/", handler(context))
	http.Handle("/api/", apiHandler(context))
	http.Handle("/api/v1/", apiHandler(context))
	http.Handle("/health/", healthHandler(context))
	http.Handle("/api/v1/metrics/", metricsHandler(context))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buffer, dot, output string

		getJSONData(&buffer)

		output = "dot"

		dot, percentageIsolated, percentageIsolatedNamespace, percentageCovered, err := processBytes([]byte(buffer), &output)
		if err != nil {
			sData := fmt.Sprintf("<p>Can't process input data: %s</p>", err)
			fmt.Fprintf(w, page(context, sData))
			return
		}

		cmd := exec.Command("dot", "-Tsvg")
		cmd.Stdin = strings.NewReader(dot)
		var svg bytes.Buffer
		cmd.Stdout = &svg
		err = cmd.Run()
		if err != nil {
			sDot := fmt.Sprintf("<p>Graphviz conversion failed: %s</p>", err)
			fmt.Fprintf(w, page(context, sDot))
			return
		}

		colorClassIsolation := "bg-success"
		if percentageIsolated < 50 {
			colorClassIsolation = "bg-danger"
		} else if percentageIsolated < 75 {
			colorClassIsolation = "bg-warning"
		}

		colorClassNamespaceIsolation := "bg-success"
		if percentageIsolatedNamespace < 50 {
			colorClassNamespaceIsolation = "bg-danger"
		} else if percentageIsolatedNamespace < 75 {
			colorClassNamespaceIsolation = "bg-warning"
		}

		colorClassCoverage := "bg-success"
		if percentageCovered < 50 {
			colorClassCoverage = "bg-danger"
		} else if percentageCovered < 75 {
			colorClassCoverage = "bg-warning"
		}

		svgString := string(svg.Bytes())

		patternW := regexp.MustCompile(`width=[^ ]+`)
		svgString = patternW.ReplaceAllLiteralString(svgString, "")
		patternH := regexp.MustCompile(`height=[^ ]+`)
		svgString = patternH.ReplaceAllLiteralString(svgString, "")
		svgString = strings.Replace(svgString, "Times,serif", "sans-serif", -1)

		buffer = fmt.Sprintf(`
	<div>%s</div>
	<br/>
	<div class="progress">
		<div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% isolation</div>
	</div>
	<br/>
	<div class="progress">
		<div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% namespace isolation</div>
	</div>
	<br/>
	<div class="progress">
		<div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% namespace coverage</div>
	</div>`,
			svgString,
			colorClassIsolation,
			percentageIsolated,
			percentageIsolated,
			percentageIsolated,
			colorClassNamespaceIsolation,
			percentageIsolatedNamespace,
			percentageIsolatedNamespace,
			percentageIsolatedNamespace,
			colorClassCoverage,
			percentageCovered,
			percentageCovered,
			percentageCovered)
		fmt.Fprintf(w, page(context, buffer))
	})
}

func apiHandler(context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buffer := fmt.Sprintf(`
<div class="row">
  <div class="col-sm-2"><a href="/">/</a></div>
  <div class="col-sm-10">show graph</div>
</div>
<div class="row">
  <div class="col-sm-2"><a href="/health/">/health/</a></div>
  <div class="col-sm-10">health endpoint</div>
</div>
<div class="row">
  <div class="col-sm-2"><a href="/api/v1/metrics/">/api/v1/metrics/</a></div>
  <div class="col-sm-10">metrics endpoint</div>
</div>`)
		fmt.Fprintf(w, page(context, buffer))
		return
	})
}

func healthHandler(context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"status\":\"ok\"}")
	})
}

func metricsHandler(context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buffer, output string

		getJSONData(&buffer)

		output = "dot"

		_, percentageIsolated, percentageIsolatedNamespaceToNamespace, percentageNamespaceCoverage, err := processBytes([]byte(buffer), &output)
		if err != nil {
			sData := fmt.Sprintf("<p>Can't process input data: %s</p>", err)
			fmt.Fprintf(w, page(context, sData))
			return
		}

		fmt.Fprintf(w, "{\"percentageIsolated\":%d,\"percentageIsolatedNamespaceToNamespace\":%d,\"percentageNamespaceCoverage\":%d}", percentageIsolated, percentageIsolatedNamespaceToNamespace, percentageNamespaceCoverage)
	})
}

func handleGet(context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buffer, dot, output string

		getJSONData(&buffer)

		output = "dot"

		dot, percentageIsolated, percentageIsolatedNamespaceToNamespace, percentageCovered, err := processBytes([]byte(buffer), &output)
		if err != nil {
			sData := fmt.Sprintf("<p>Can't process input data: %s</p>", err)
			fmt.Fprintf(w, page(context, sData))
			return
		}

		cmd := exec.Command("dot", "-Tsvg")
		cmd.Stdin = strings.NewReader(dot)
		var svg bytes.Buffer
		cmd.Stdout = &svg
		err = cmd.Run()
		if err != nil {
			sDot := fmt.Sprintf("<p>Graphviz conversion failed: %s</p>", err)
			fmt.Fprintf(w, page(context, sDot))
			return
		}

		colorClassIsolation := "progress-bar-success"
		if percentageIsolated < 50 {
			colorClassIsolation = "progress-bar-danger"
		} else if percentageIsolated < 75 {
			colorClassIsolation = "progress-bar-warning"
		}

		colorClassNamespaceIsolation := "progress-bar-success"
		if percentageIsolatedNamespaceToNamespace < 50 {
			colorClassNamespaceIsolation = "progress-bar-danger"
		} else if percentageIsolatedNamespaceToNamespace < 75 {
			colorClassNamespaceIsolation = "progress-bar-warning"
		}

		colorClassCoverage := "progress-bar-success"
		if percentageCovered < 50 {
			colorClassCoverage = "progress-bar-danger"
		} else if percentageCovered < 75 {
			colorClassCoverage = "progress-bar-warning"
		}

		svgString := string(svg.Bytes())

		patternW := regexp.MustCompile(`width=[^ ]+`)
		svgString = patternW.ReplaceAllLiteralString(svgString, "")
		patternH := regexp.MustCompile(`height=[^ ]+`)
		svgString = patternH.ReplaceAllLiteralString(svgString, "")
		svgString = strings.Replace(svgString, "Times,serif", "sans-serif", -1)

		buffer = fmt.Sprintf(`
	<div>%s</div>
	<div class="progress">
	<div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% isolation</div>
	</div>
	<div class="progress">
	<div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% namespace isolation</div>
	</div>
	<div class="progress">
	<div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% namespace coverage</div>
	</div>`,
			svgString,
			colorClassIsolation,
			percentageIsolated,
			percentageIsolated,
			percentageIsolated,
			colorClassNamespaceIsolation,
			percentageIsolatedNamespaceToNamespace,
			percentageIsolatedNamespaceToNamespace,
			percentageIsolatedNamespaceToNamespace,
			colorClassCoverage,
			percentageCovered,
			percentageCovered,
			percentageCovered)
		fmt.Fprintf(w, page(context, buffer))
	})
}

func handlePost(context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Fprintf(w, "Can't read POST request body '%s': %s", body, err)
			return
		}
	})
}
