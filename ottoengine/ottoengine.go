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
	Filename   string
	Source []byte
	VM     *otto.Otto
    Script *otto.Script
}

func Load(root string) ([]Program, error) {
	var programs []Program

	err := filepath.Walk(root, func(filename string, file_info os.FileInfo, err error) error {
		// Trim the leading path from the path string Trim ext from path string, save this as route.
		ext := path.Ext(filename)
		if file_info.IsDir() != true && ext == ".js" {
			if ext == ".js" {
				route := strings.TrimSuffix(strings.TrimPrefix(filename, root), ".js")
				source, err := ioutil.ReadFile(filename)
				if err != nil {
					log.Println(err)
					return err
				}
				vm := otto.New()
                full_path, err := filepath.Abs(filename)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }

                // Attempt to compile source and abort is there is a problem
                script, err := vm.Compile(full_path, source)
                if err != nil {
                    fmt.Printf("File: %s, %s\n", full_path, err)
                    os.Exit(1)
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
	    output, err := program.VM.Run(program.Script)
	    if err != nil {
            log.Printf("file %s: \nerror: %s\n", program.Filename, err)
            http.Error(w, "Internal Server Error", 500)
		    return
	    }
	    // This write the body, should really write headers and render into body, etc.
	    fmt.Fprintf(w, "%s\n", output)
	})
}

func AddRoutes(programs []Program) {
	for i := range programs {
		fmt.Printf("Adding route %d %s from %s\n", i, programs[i].Route, programs[i].Filename)
		Engine(programs[i])
	}
}

