//
// Package ottoengine is a module providing a way to define route processing
// using the Otto JavaScript virutal machine.
//
// Otto is written by Robert Krimen, see https://github.com/robertkrimen/otto
// Otto Engine is written by Robert Doiel, see https://github.com/rsdoiel/ws
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the Simplified BSD License.
// See: http://opensource.org/licenses/BSD-2-Clause
//
package ottoengine

import (
	"../wslog"
	//"./reload"
	"bytes"
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

// Program keeps track of assigned route, the JS file, the source of
// the JS file as well as handles for the VM and Script elements in otto.
type Program struct {
	Route    string
	Filename string
	Source   []byte
	VM       *otto.Otto
	Script   *otto.Script
}

// LoadFile ingests a JavaScript file and returns a Program struct.
func LoadFile(root string, filename string, fileInfo os.FileInfo, err error) (*Program, error) {
	var (
		ext      string
		fullPath string
		route    string
		source   []byte
		vm       *otto.Otto
		script   *otto.Script
	)

	// Trim the leading path from the path string Trim ext from path string, save this as route.
	ext = path.Ext(filename)
	if fileInfo != nil && fileInfo.IsDir() != true && ext == ".js" {
		if ext == ".js" {
			route = strings.TrimSuffix(strings.TrimPrefix(filename, root), ".js")
			log.Printf("Reading %s\n", filename)
			source, err = ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			vm = otto.New()
			fullPath, err = filepath.Abs(filename)
			if err != nil {
				return nil, err
			}

			// Attempt to compile source and abort is there is a problem
			script, err = vm.Compile(fullPath, source)
			if err != nil {
				return nil, fmt.Errorf("File: %s, %s\n", fullPath, err)
			}
		}
	}
	return &Program{Route: route, Filename: filename, Source: source, VM: vm, Script: script}, nil
}

// Load evals a string and returns a Program Struct. This is an alternative to reading
// a file and eval with LoadFile().
func Load(root string) ([]Program, error) {
	var programs []Program

	err := filepath.Walk(root, func(filename string, fileInfo os.FileInfo, err error) error {
		ext := path.Ext(filename)
		if fileInfo != nil && fileInfo.IsDir() != true && ext == ".js" {
			p, err := LoadFile(root, filename, fileInfo, err)
			if err != nil {
				return err
			}
			programs = append(programs, Program{Route: p.Route, Filename: p.Filename, Source: p.Source, VM: p.VM, Script: p.Script})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func jsMethodHandler(method string, data string) string {
	return fmt.Sprintf(`function (debug) {
        var data_as_string = %q.trim(),
            key_value_as_string = [],
            data = {},
            space_re = /\+/g;

        if (debug === undefined) {
            debug = false;
        }
        if (this.Method === "%s") {
            // Is it a JSON post?  
            if ((data_as_string.substr(0,1) === "{" && 
                    data_as_string.substr(-1, 1) === "}") ||
                    (data_as_string.substr(0,1) === "[" && 
                    data_as_string.substr(-1, 1) === "]")) {
                 data = JSON.parse(data_as_string) 
            } else {
                // It's a normal URL encoded post
                key_value_as_string = data_as_string.split("&");
                if (key_value_as_string.length > 0) {
                    key_value_as_string.forEach(function (item) {
                        var parts = item.split("=",2);
                
                        if (parts.length === 2) {
                            key = decodeURIComponent(parts[0].replace(space_re, ' '));
                            value = decodeURIComponent(parts[1].replace(space_re, ' '));
                            data[key] = value;
                        }
                    });
                }
            }
        }
        if (debug === true) {
            console.log("%s REQUEST: ", JSON.stringify(data));
        }
        return data;
    }`, data, method, method)
}

func jsInjectMethod(name string, src string, buf []byte) string {
	end := bytes.LastIndex(buf, []byte("}"))
	if end > -1 {
		// Insert our literal function def for GET
		src = fmt.Sprintf("%s,%s:%s%s", string(buf[0:end]), name, src, string(buf[end]))
	} else {
		src = string(buf)
	}
	return src
}

func requestAsJSON(r *http.Request) (buf []byte, err error) {
	jsonURL, err := json.Marshal(r.URL)
	return []byte(fmt.Sprintf("{\"Method\":\"%s\",\"URL\":%s}", r.Method, jsonURL)), err
}

func createRequestLiteral(r *http.Request) string {
	buf, err := requestAsJSON(r)
	if err != nil {
		return "{}"
	}
	switch r.Method {
	case "POST":
		//FIXME: need to handle multi-part requests (E.g. uploading a file)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("POST read eror: %s", err)
			return string(buf)
		}
		return jsInjectMethod("POST", jsMethodHandler(r.Method, string(body)), buf)
	}
	return jsInjectMethod(r.Method, jsMethodHandler(r.Method, r.URL.RawQuery), buf)
}

// Response sets up the structure for forming HTTP headers for
// OttoEngine scripts.
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

// IsJSON is a predicate that returns true if the content is JSON otherwise false.
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

// IsHTML is a predicate the returns true if the content is HTML otherwise false.
func IsHTML(value otto.Value) bool {
	blob, _ := value.ToString()
	return strings.HasPrefix(blob, "<!DOCTYPE html>")
}

func numberLines(src string) string {
	var (
		results []string
	)

	for i, line := range strings.Split(src, "\n") {
		results = append(results, fmt.Sprintf("%d:\t%s\n", i, line))
	}

	return strings.Join(results, "")
}

// Engine provides the basic service for ws.go to handle JavaScript interactions server side.
func Engine(program Program) {
	http.HandleFunc(program.Route, func(w http.ResponseWriter, r *http.Request) {
		var (
			vm              *otto.Otto
			script          string
			closureScript   string
			runScript       string
			requestLiteral  string
			responseLiteral string
			goResponse      Response
		)

		// 1 Create fresh Request object literal.
		requestLiteral = createRequestLiteral(r)

		// 2. Create a fresh Response object literal.
		responseLiteral = createResponseLiteral()

		// 3. Setup the VM for the Route with our closure
		vm = program.VM
		script = string(program.Source)
		closureScript = `JSON.stringify((function(Request,Response){var value = %s;if (value) { Response.setContent(value); };return Response;}(%s,%s)));`
		runScript = fmt.Sprintf(closureScript, script, requestLiteral, responseLiteral)

		// 4. Load built in commands (e.g. Getenv, HttpGet, HttpPost)
		vm.Set("Getenv", func(call otto.FunctionCall) otto.Value {
			envvar := call.Argument(0).String()
			result, err := vm.ToValue(os.Getenv(envvar))
			if err != nil {
				log.Fatalf("Getenv(%q) error, %s", envvar, err)
			}
			return result
		})
		vm.Set("HttpGet", func(call otto.FunctionCall) otto.Value {
			uri := call.Argument(0).String()
			resp, err := http.Get(uri)
			if err != nil {
				log.Fatalf("Can't connect to %s, %s", uri, err)
			}
			defer resp.Body.Close()
			content, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Can't read response %s, %s", uri, err)
			}
			result, err := vm.ToValue(fmt.Sprintf("%s", content))
			if err != nil {
				log.Fatalf("HttpGet(%q) error, %s", uri, err)
			}
			return result
		})
		vm.Set("HttpPost", func(call otto.FunctionCall) otto.Value {
			uri := call.Argument(0).String()
			mimeType := call.Argument(1).String()
			payload := call.Argument(2).String()
			buf := strings.NewReader(payload)

			resp, err := http.Post(uri, mimeType, buf)
			if err != nil {
				log.Fatalf("Can't connect to %s, %s", uri, err)
			}
			defer resp.Body.Close()
			content, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Can't read response %s, %s", uri, err)
			}
			result, err := vm.ToValue(fmt.Sprintf("%s", content))
			if err != nil {
				log.Fatalf("HttpGet(%q) error, %s", uri, err)
			}
			return result
		})

		// 5. Run the VM wrapped with a closure containing`Request, Response
		output, err := vm.Run(runScript)
		if err != nil {
			msg := fmt.Sprintf("JavaScript Error: %s", err)
			wslog.LogResponse(500, "Internal Server Error",
				r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
			log.Println(numberLines(runScript))
			http.Error(w, "Internal Server Error", 500)
			return
		}
		jsonSrc, err := output.ToString()
		if err != nil {
			msg := fmt.Sprintf("Conversion to JSON: %s", err)
			wslog.LogResponse(500, "Internal Server Error",
				r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		// See if we're rendering from a returned text string or JSON
		// via the Response object.
		jsonErr := json.Unmarshal([]byte(jsonSrc), &goResponse)
		if jsonErr != nil {
			// We're rendering from a text string, try to calc the content type.
			// 5. Calc headers
			contentType := "text/plain; charset=utf-8"
			if IsJSON(output) {
				contentType = "application/json; charset=utf-8"
			} else if IsHTML(output) {
				contentType = "text/html; charset=utf-8"
			}
			w.Header().Set("Content-Type", contentType)

			// 6. send the output to the browser.
			fmt.Fprintf(w, "%s", output)
			wslog.LogResponse(200, "OK", r.Method, r.URL, r.RemoteAddr, program.Filename, contentType)
			return
		}

		// We're rendering completely from the response object.
		// 5. update headers from responseObject
		contentType := "text/plain; charset=utf-8"
		for key, value := range goResponse.Headers {
			if key == "content-type" {
				contentType = value
			}
			w.Header().Set(key, value)
		}
		fmt.Fprintf(w, "%s", goResponse.Content)
		wslog.LogResponse(goResponse.Code, goResponse.Status, r.Method, r.URL, r.RemoteAddr, program.Filename, contentType)
	})
}

// AddRoutes allows a Program structure to fine an approach route handler.
func AddRoutes(programs []Program) {
	for i := range programs {
		log.Printf("Adding route (%d) %s from %s\n", i, programs[i].Route, programs[i].Filename)
		Engine(programs[i])
	}
}
