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
    "net/url"
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

func createRequestSource(r *http.Request) (string, error) {
    buf, err := json.Marshal(r.Header)
    src := fmt.Sprintf(`Request = {
            Headers: %s,
            Method: %q,
            URL: %q,
            Proto: %q,
            Referrer: %q,
            UserAgent: %q
    };`, string(buf), r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())

    //FIXME: need to handle GET, POST, PUT, DELETE methods
    return src, err
}

func createResponseSource() (string, error) {
    src := `Response = {
        ContentType: "",
        Location: ""
    };`
    return src, nil
}

func log_response(code int, msg string, filename_or_src string, method string, url *url.URL, proto string, referrer string, user_agent string) {
	log.Printf("{\"response\": %d, \"status\": %q, \"filename\": %q, %q: %q, \"protocol\": %q, \"referrer\": %q, \"user-agent\": %q}\n",
        code,
        msg,
        filename_or_src, 
        method,
        url,
        proto,
        referrer,
        user_agent)
}

func Engine(program Program) {
	http.HandleFunc(program.Route, func(w http.ResponseWriter, r *http.Request) {
		// FIXME:
		// 1. Create fresh Response and Request objects.
        request, err := createRequestSource(r)
        if err != nil {
            msg := fmt.Sprintf("%s", err)
            log_response(500, msg, request, r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
        }
        
        response, err := createResponseSource()
        if err != nil {
            msg := fmt.Sprintf("%s", err)
            log_response(500, msg, response, r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
        }

		// 2. Run the VM passing with Request, Response objects already created
        combined_src := fmt.Sprintf("%s\n%s\n\n%s\n", request, response, program.Source)
		output, err := program.VM.Run(combined_src)
		if err != nil {
            msg := fmt.Sprintf("%s", err)
            log_response(500, msg, program.Filename, r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		// 3. based on state of Response object
		//    a. update headers in ResponseWriter
		//    b. take care of any encoding issues 
        //    c. send back the contents of output
        //value, _ := program.VM.Get("Response.ContentType");
        //src, _ := value.ToString()
        //fmt.Printf("DEBUG, content types? %v %s\n", value, src);

        log_response(200, "OK", program.Filename, r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
		fmt.Fprintf(w, "%s\n", output)
	})
}

func AddRoutes(programs []Program) {
	for i := range programs {
		log.Printf("Adding route (%d) %s from %s\n", i, programs[i].Route, programs[i].Filename)
		Engine(programs[i])
	}
}
