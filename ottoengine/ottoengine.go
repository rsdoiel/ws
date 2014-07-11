/**
 * ottoengine.go - ottoengine module provides a way to define route processing
 * using the Otto JavaScript virutal machine.
 *
 * Otto is written by Robert Krimen, see https://github.com/robertkrimen/otto
 * Otto Engine is written by Robert Doiel, see https://github.com/rsdoiel/ws
 */
package ottoengine

import (
	"../wslog"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Program struct {
	Route    string
	Filename string
	Source   []byte
	VM       *otto.Otto
	Script   *otto.Script
}

func Load(root string) ([]Program, error) {
	var programs []Program

	err := filepath.Walk(root, func(filename string, file_info os.FileInfo, err error) error {
		// Trim the leading path from the path string Trim ext from path string, save this as route.
		ext := path.Ext(filename)
		if file_info != nil && file_info.IsDir() != true && ext == ".js" {
			if ext == ".js" {
				route := strings.TrimSuffix(strings.TrimPrefix(filename, root), ".js")
				log.Printf("Reading %s\n", filename)
				source, err := ioutil.ReadFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				vm := otto.New()
				full_path, err := filepath.Abs(filename)
				if err != nil {
					log.Fatal(err)
				}

				// Attempt to compile source and abort is there is a problem
				script, err := vm.Compile(full_path, source)
				if err != nil {
					log.Fatalf("File: %s, %s\n", full_path, err)
				}
				programs = append(programs, Program{Route: route, Filename: filename, Source: source, VM: vm, Script: script})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func createRequestLiteral(r *http.Request) string {
	var src string

	buf, err := json.Marshal(r)
	if err != nil {
		src = "{}"
	} else {
		src = string(buf)
	}

	return src
}

type Response struct {
	Code    int
	Status  string
	Headers map[string]string //interface{}
	Content string
}

func createResponseLiteral() string {
	src := `{
		Code: 200,
		Status: "OK",
		Headers: {"content-type": "text/plain"},
		Content: "",
		setHeader: function (key, value) {
			this.Headers[key.toLowerCase()] = value;
			return (this.Headers[key.toLowerCase()] === value);
		},
		getHeader: function (key) {
			if (this.Headers[key.toLowerCase()] !== undefined) {
				return this.Headers[key.toLowerCase()];
			}
			return false;
		},
		setContent: function (content) {
			this.Content = content;
			return (this.Content === content);
		},
		getContent: function () {
			return this.Content;
		}
	}`
	return src
}

func IsJSON(value otto.Value) bool {
	blob, _ := value.ToString()
	if (strings.HasPrefix(blob, "[\"") == true &&
		strings.HasSuffix(blob, "\"]") == true) ||
		(strings.HasPrefix(blob, "{\"") == true &&
			strings.HasSuffix(blob, "\"}") == true) {
		return true
	}
	return false
}

func IsHTML(value otto.Value) bool {
	blob, _ := value.ToString()
	return strings.HasPrefix(blob, "<!DOCTYPE html>")
}

func Engine(program Program) {
	http.HandleFunc(program.Route, func(w http.ResponseWriter, r *http.Request) {
		var (
			vm               *otto.Otto
			script           string
			closure_script   string
			run_script       string
			request_literal  string
			response_literal string
			go_response      Response
		)

		// 1. Create fresh Request object literal.
		request_literal = createRequestLiteral(r)

		// 2. Create a fresh Response object literal.
		response_literal = createResponseLiteral()

		// 3. Setup the VM for the Route with our closure
		vm = program.VM
		script = string(program.Source)
		closure_script = `JSON.stringify((function(Request,Response){var value = %s;if (value) { Response.setContent(value); };return Response;}(%s,%s)));`
		run_script = fmt.Sprintf(closure_script, script, request_literal, response_literal)

		// 4. Run the VM wrapped with a closure containing`Request, Response
		output, err := vm.Run(run_script)
		if err != nil {
			msg := fmt.Sprintf("Script: %s", err)
			wslog.LogResponse(500, "Internal Server Error",
				r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		json_src, err := output.ToString()
		if err != nil {
			msg := fmt.Sprintf("Conversion to JSON: %s", err)
			wslog.LogResponse(500, "Internal Server Error",
				r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		json_err := json.Unmarshal([]byte(json_src), &go_response)
		if json_err != nil {
			content_type := "text/plain"
			// 5. Calc headers
			if IsJSON(output) {
				content_type = "application/json"
			} else if IsHTML(output) {
				content_type = "text/html"
			}
			w.Header().Set("Content-Type", "text/html")
			// 6. send the output to the browser.
			fmt.Fprintf(w, "%s", output)
			wslog.LogResponse(200, "OK", r.Method, r.URL, r.RemoteAddr, program.Filename, content_type)
			return

		}
		fmt.Printf("DEBUG go_response: code: %d, status: %s, content: %s, headers: %v\n",
			go_response.Code, go_response.Status, go_response.Content, go_response.Headers)

		// 5. update headers from responseObject
		content_type := "text/plain"
		for key, value := range go_response.Headers {
			w.Header().Set(key, value)
		}

		fmt.Fprintf(w, "%s", go_response.Content)
		wslog.LogResponse(go_response.Code, go_response.Status, r.Method, r.URL, r.RemoteAddr, program.Filename, content_type)
	})
}

func AddRoutes(programs []Program) {
	for i := range programs {
		log.Printf("Adding route (%d) %s from %s\n", i, programs[i].Route, programs[i].Filename)
		Engine(programs[i])
	}
}
