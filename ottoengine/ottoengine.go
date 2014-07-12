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
	//"./reload"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
    "bytes"
	"strings"
	"errors"
)

type Program struct {
	Route    string
	Filename string
	Source   []byte
	VM       *otto.Otto
	Script   *otto.Script
}

func LoadFile(root string, filename string, file_info os.FileInfo, err error) (*Program, error) {
	var (
		ext string
		full_path string
		route string
		source []byte
		vm *otto.Otto
		script *otto.Script
	)

	// Trim the leading path from the path string Trim ext from path string, save this as route.
	ext = path.Ext(filename)
	if file_info != nil && file_info.IsDir() != true && ext == ".js" {
		if ext == ".js" {
			route = strings.TrimSuffix(strings.TrimPrefix(filename, root), ".js")
			log.Printf("Reading %s\n", filename)
			source, err = ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			vm = otto.New()
			full_path, err = filepath.Abs(filename)
			if err != nil {
				return nil, err
			}

			// Attempt to compile source and abort is there is a problem
			script, err = vm.Compile(full_path, source)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("File: %s, %s\n", full_path, err))
			}
		}
	}
	return &Program{Route: route, Filename: filename, Source: source, VM: vm, Script: script}, nil
}

func Load(root string) ([]Program, error) {
	var programs []Program

	err := filepath.Walk(root, func(filename string, file_info os.FileInfo, err error) error {
	    ext := path.Ext(filename)
	    if file_info != nil && file_info.IsDir() != true && ext == ".js" {
		    p, err := LoadFile(root, filename, file_info, err)
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

func jsGET() string {
    return `function () {
        var raw_params = [],
            getargs = {},
            space_re = /\+/g;
        if (this.Method === "GET") {
            raw_params = this.URL.RawQuery.split("&");
            if (raw_params.length > 0) {
                raw_params.forEach(function (item) {
                    var parts = item.split("=",2);
                
                    if (parts.length === 2) {
                        key = decodeURIComponent(parts[0].replace(space_re, ' '));
                        value = decodeURIComponent(parts[1].replace(space_re, ' '));
                        getargs[key] = value;
                    }
                });
            }
        }
        return getargs;
    }`
}

func jsPOST(post string) string {
    src := `function () {
        var post_string = %q,
            raw_params = [],
            postargs = {},
            space_re = /\+/g;

        if (this.Method === "POST") {
            raw_params = post_string.split("&");
            if (raw_params.length > 0) {
                raw_params.forEach(function (item) {
                    var parts = item.split("=",2);
                
                    if (parts.length === 2) {
                        key = decodeURIComponent(parts[0].replace(space_re, ' '));
                        value = decodeURIComponent(parts[1].replace(space_re, ' '));
                        postargs[key] = value;
                    }
                });
            }
        }
        return postargs;
    }`
    return fmt.Sprintf(src, post)
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

func createRequestLiteral(r *http.Request) string {
	var src string

	buf, err := json.Marshal(r)
	if err != nil {
		src = "{}"
	} else {
        switch r.Method {
            case "GET":
                return jsInjectMethod("GET", jsGET(), buf) 
            case "POST":
                body, err := ioutil.ReadAll(r.Body)
                if err != nil {
                    log.Printf("POST read eror: %s", err)
                    return string(buf)
                }
                return jsInjectMethod("POST", jsPOST(string(body)), buf)
            default:
                src = string(buf)
        }
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

		// 1 Create fresh Request object literal.
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
        // See if we're rendering from a returned text string or JSON
        // via the Response object.
		json_err := json.Unmarshal([]byte(json_src), &go_response)
		if json_err != nil {
            // We're rendering from a text string, try to calc the content type.
			// 5. Calc headers
			content_type := "text/plain; charset=utf-8"
			if IsJSON(output) {
				content_type = "application/json; charset=utf-8"
			} else if IsHTML(output) {
				content_type = "text/html; charset=utf-8"
			}
			w.Header().Set("Content-Type", content_type)

			// 6. send the output to the browser.
			fmt.Fprintf(w, "%s", output)
			wslog.LogResponse(200, "OK", r.Method, r.URL, r.RemoteAddr, program.Filename, content_type)
			return
		}

        // We're rendering completely from the response object.
		// 5. update headers from responseObject
		content_type := "text/plain; charset=utf-8"
		for key, value := range go_response.Headers {
            if key == "content-type" {
                content_type = value
            }
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
