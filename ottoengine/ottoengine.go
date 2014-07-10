/**
 * ottoengine.go - ottoengine module provides a way to define route processing using
 * the Otto JavaScript virutal machine.
 * Otto is written by Robert Krimen, see https://github.com/robertkrimen/otto
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

func createRequestObject(vm *otto.Otto, r *http.Request) (*otto.Object, error) {
	buf, err := json.Marshal(r.Header)
	if err != nil {
		return nil, err
	}

	//FIXME: need to handle GET, POST, PUT, DELETE methods
	src := fmt.Sprintf(`Request = {
            Headers: %s,
            Method: %q,
            URL: %q,
            Proto: %q,
            Referrer: %q,
            UserAgent: %q
    };`, string(buf), r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())

	obj, err := vm.Object(src)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func createResponseObject(vm *otto.Otto) (*otto.Object, error) {
	src := `Response = {
        code: 200,
        status: "OK",
        headers: {},
        header_keys: [],
        getHeader: function (key) {
            if (this.headers[key.toLowerCase()] !== undefined) {
                return this.headers[key.toLowerCase()];
            }
            return false;
        },
        setHeader: function (key, value) {
            this.headers[key.toLowerCase()] = value;
            return (this.headers[key.toLowerCase()] === value);
        },
        collectHeaderKeys: function () {
            return this.header_keys = Object.keys(this.headers);
	},
	popHeaderKey: function () {
	    return this.header_keys.pop();
	}
    }`

	obj, err := vm.Object(src)
	if err != nil {
		return nil, err
	}
	return obj, nil
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
		// 1. Create a fresh VM
		vm := otto.New()

		// 2. Create fresh Request object.
		createRequestObject(vm, r)
		// 3. Create a fresh Response object.
		createResponseObject(vm)
		// 4. Run the VM passing Request along with Script
		output, err := vm.Run(program.Script)
		if err != nil {
			msg := fmt.Sprintf("Script: %s", err)
			wslog.LogResponse(500, "Internal Server Error", r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		// 3. update headers from responseObject
		content_type := ""
		vm.Run(`Response.collectHeaderKeys()`)
		key_cnt_value, _ := vm.Run("Response.header_keys.length")
		key_cnt, _ := key_cnt_value.ToInteger()
		for i := int64(0); i < key_cnt; i++ {
			key_value, _ := vm.Run(`Response.popHeaderKey();`)
			key, _ := key_value.ToString()
			if key != "" {
				value_value, _ := vm.Run(fmt.Sprintf("Response.getHeader(%q);", key))
				value, _ := value_value.ToString()
				if value != "" {
					if key == "content-type" {
						content_type = value
					}
					w.Header().Set(key, value)
				}
			}
		}
		// 4. Calc fallback content types if needed.
		if content_type == "" && IsJSON(output) {
			w.Header().Set("Content-Type", "application/json")
		} else if content_type == "" && IsHTML(output) {
			w.Header().Set("Content-Type", "text/html")
		}
		// 5. send the output to the browser.

		fmt.Fprintf(w, "%s\n", output)
		wslog.LogResponse(200, "OK", r.Method, r.URL, r.RemoteAddr, program.Filename, "")
	})
}

func AddRoutes(programs []Program) {
	for i := range programs {
		log.Printf("Adding route (%d) %s from %s\n", i, programs[i].Route, programs[i].Filename)
		Engine(programs[i])
	}
}
