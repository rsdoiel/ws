/**
 * ottoengine.go - ottoengine module provides a way to define route processing using
 * the Otto JavaScript virutal machine.
 * Otto is written by Robert Krimen, see https://github.com/robertkrimen/otto
 */
package ottoengine

import (
	"fmt"
	"log"
	"os"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type Program struct {
	Route  string
	Path   string
	Source []byte
	VM     *otto.Otto
    Script *otto.Otto
}

func Load(root string) (map[string]Program, error) {
	programs := make(map[string]Program)
	err := filepath.Walk(root, func(p string, file_info os.FileInfo, err error) error {
		// Trim the leading path from the path string Trim ext from path string, save this as route.
		ext := path.Ext(p)
		if file_info.IsDir() != true && ext == ".js" {
			if ext == ".js" {
				route := strings.TrimSuffix(strings.TrimPrefix(p, root), ".js")
				source, err := ioutil.ReadFile(p)
				if err != nil {
					log.Println(err)
					return err
				}
				vm := otto.New()
                full_path, err := filepath.Abs(p)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }

                script, err := vm.Compiled(full_path, source)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
				programs[route] = Program{Route: route, Path: p, Source: source, VM: vm, Script: script}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func Engine(w http.ResponseWriter, r *http.Request, program Program) {
    //FIXME: Need to handle setting appropriate http headers from Otto VM.
	output, err := program.VM.Run(program.Script)
	if err != nil {
        log.Printf("file %s: \nerror: %s\n", program.Path, err)
        http.Error(w, "Internal Server Error", 500)
		return
	}
	// This write the body, should really write headers and render into body, etc.
	fmt.Fprintf(w, "%s\n", output)
}

func AddRoutes(programs map[string]Program) {
	for route, program := range programs {
		fmt.Printf("Creating route %s from %s\n", route, program.Path)
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			Engine(w, r, program)
		})
	}
}

