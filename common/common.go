/**
 * common.go - Process default environment values.
 */
package common

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var Usage = func(exit_code int, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	fmt.Fprintf(fh, `%s
 USAGE %s [options]

 EXAMPLES
      
 OPTIONS

`, msg, os.Args[0])

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(fh, "\t-%s\t(defaults to %s) %s\n", f.Name, f.DefValue, f.Usage)
	})

	fmt.Fprintf(fh, `

 copyright (c) 2014 all rights reserved.
 Released under the Simplified BSD License
 See: http://opensource.org/licenses/bsd-license.php

`)
	os.Exit(exit_code)
}

func RequestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wslog.LogRequest(r.Method, r.URL, r.RemoteAddr, r.Proto, r.Referer(), r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

func DefaultEnvBool(environmentVar string, defaultValue bool) bool {
	tmp := strings.ToLower(os.Getenv(environmentVar))
	if tmp == "true" {
		return true
	}
	if tmp == "false" {
		return false
	}
	return defaultValue
}

func DefaultEnvString(environmentVar string, defaultValue string) string {
	tmp := os.Getenv(environmentVar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

func DefaultEnvInt(environmentVar string, defaultValue int) int {
	tmp := os.Getenv(environmentVar)
	if tmp != "" {
		i, err := strconv.Atoi(tmp)
		if err != nil {
			usage(1, environmentVar+" must be an integer.")
		}
		return i
	}
	return defaultValue
}
