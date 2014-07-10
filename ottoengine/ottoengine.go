/**
 * ottoengine.go - ottoengine module provides a way to define route processing using
 * the Otto JavaScript virutal machine.
 * Otto is written by Robert Krimen, see https://github.com/robertkrimen/otto
 */
package ottoengine

import (
    "../logger"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
    "encoding/json"
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
    if err != nil  {
        return nil, err
    }
    return obj, nil
}

func createResponseObject(vm *otto.Otto) (*otto.Object, error) {
    src := `Response = {
        headers: [],
        setHeader: function (key, value, replace) {
            var i = 0;
            if (replace === undefined) {
                replace = true;
            }
            if (replace === true) {
                for (i = 0; i < this.headers.length; i += 1) {
                    if (this.headers[i].key === key) {
                        this.headers[i].value = value;
                        return (this.headers[i].key === key && this.headers[i].value === value);
                    }
                }
            }
            return this.headers.push({key: key, value: value});
        }
    }`

    obj, err := vm.Object(src)
    if err != nil  {
        return nil, err
    }
    return obj, nil
}

func isJSON (value otto.Value) bool {
    blob, _ := value.ToString()
    if (strings.HasPrefix(blob, "[\"") == true &&
            strings.HasSuffix(blob, "\"]") == true) || 
            (strings.HasPrefix(blob, "{\"") == true && strings.HasSuffix(blob, "\"}") == true) {
        return true
    }
    return false
}

func isHTML (value otto.Value) bool {
    blob, _ := value.ToString()
    return strings.HasPrefix(blob, "<!DOCTYPE html>")
}

func Engine(program Program) {
	http.HandleFunc(program.Route, func(w http.ResponseWriter, r *http.Request) {
        // 1. Create a fresh VM
        vm := otto.New()

		// 2. Create fresh Request object.
        requestObject, err := createRequestObject(vm, r)
        if err != nil {
            msg := fmt.Sprintf("Request Object: %s", err)
            logger.LogResponse(500, "Internal Server Error", r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
            http.Error(w, "Internal Server Error", 500)
            return;
        }
        
        // 3. Create a fresh Response object.
        responseObject, err := createResponseObject(vm)
        if err != nil {
            msg := fmt.Sprintf("Response Object: %s", err)
            logger.LogResponse(500, "Internal Server Error", r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
            http.Error(w, "Internal Server Error", 500)
            return;
        }

		// 4. Run the VM passing Request along with Script
		output, err := vm.Run(program.Script)
		if err != nil {
            msg := fmt.Sprintf("Script: %s", err)
            logger.LogResponse(500, "Internal Server Error", r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		// 3. based on returned output
		//    a. update headers from responseObject
        value, err := responseObject.Get("headers.length")
        if err != nil {
            msg := fmt.Sprintf("Headers length: %s", err)
            logger.LogResponse(500, "Internal Server Error", r.Method, r.URL, r.RemoteAddr, program.Filename, msg)
			http.Error(w, "Internal Server Error", 500)
			return
        }

        header_length, err := value.ToInteger()
        if err != nil {
            header_length = 0
        }

        for i := 0; i < header_length; i++ {
            key_value, _ := responseObject.Get(fmt.Sprintf("header[%d].key", i))
            value_value, _ := responseObject.Get(fmt.Sprintf("header[%d].value", i))
            key, _ := key_value.ToString()
            value, _ := value_value.ToString()
            w.Header().Set(key, value)
        }

        /*
        //    b. send output
        if isJSON(output) {
            w.Header().Set("Content-Type", "application/json")
        } else if isHTML(output) {
            w.Header().Set("Content-Type", "text/html")
        }*/
         
		fmt.Fprintf(w, "%s\n", output)
        logger.LogResponse(200, "OK", r.Method, r.URL, r.RemoteAddr, program.Filename, "")
	})
}

func AddRoutes(programs []Program) {
	for i := range programs {
		log.Printf("Adding route (%d) %s from %s\n", i, programs[i].Route, programs[i].Filename)
		Engine(programs[i])
	}
}
