/**
 * ottoengine.go - ottoengine module provides a way to define route processing using
 * the Otto JavaScript virutal machine.
 * Otto is written by Robert Krimen, see https://github.com/robertkrimen/otto
 */
package ottoengine

import (
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

func Engine(program Program) {
	http.HandleFunc(program.Route, func(w http.ResponseWriter, r *http.Request) {
		// FIXME:
		// 1. Create fresh Response and Request objects.
		// 2. Run the VM passing in Response, Request objects via this.Response and this.Request.
		output, err := program.VM.Run(program.Script)
		if err != nil {
            log.Printf("{\"status\": 500, \"filename\": %q, \"error\": %q, %q: %q, \"protocol\": %q, \"referrer\": %q, \"user-agent\": %q}\n",
                program.Filename,
                err,
                r.Method,
                r.URL,
                r.Proto,
                r.Referer(),
                r.UserAgent())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		// 3. based on state of Response object
		//    a. update headers in ResponseWriter
		//    b. take care of any encoding issues and send back the contents of output
        log.Printf("{\"status\": 200, %q: %q, \"protocol\": %q, \"referrer\": %q, \"user-agent\": %q}\n",
            r.Method,
            r.URL,
            r.Proto,
            r.Referer(),
            r.UserAgent())
		fmt.Fprintf(w, "%s\n", output)
	})
}

func AddRoutes(programs []Program) {
	for i := range programs {
		log.Printf("Adding route (%d) %s from %s\n", i, programs[i].Route, programs[i].Filename)
		Engine(programs[i])
	}
}
