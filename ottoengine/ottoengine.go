/**
 * ottoengine.go - this module defines wraps the Otto JS VM and allows support of
 * defining route handles in JavaScript for ws.go
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
				programs[route] = Program{Route: route, Path: p, Source: source, VM: vm}
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
	output, err := program.VM.Run(program.Source)
	if err != nil {
		fmt.Fprintf(w, "file: %s\nerror: %s\n", program.Path, err)
		return
	}
	// This write the body, should really write headers and render into body, etc.
	fmt.Fprintf(w, "%s\n", output)
}
