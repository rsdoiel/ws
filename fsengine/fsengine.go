//
// Package fsengine is a simple restricted FileServer Handler for
// ws.go. The primary difference between the stock http file server
// and fsengine is fsengine avoids exposing dot files via http requests.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the Simplified BSD License.
// See: http://opensource.org/licenses/BSD-2-Clause
//
package fsengine

import (
	"../cfg"
	"../wslog"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

// Engine is a Restricted FileService excluding dot files and directories.
func Engine(config *cfg.Cfg, w http.ResponseWriter, r *http.Request) {
	var (
		hasDotPath = regexp.MustCompile(`\/\.`)
		docroot    = config.Docroot
	)

	uncleanPath := r.URL.Path
	if !strings.HasPrefix(uncleanPath, "/") {
		uncleanPath = "/" + uncleanPath
	}
	cleanPath := path.Clean(uncleanPath)
	r.URL.Path = cleanPath
	resolvedPath := path.Clean(path.Join(docroot, cleanPath))
	_, err := os.Stat(resolvedPath)
	if hasDotPath.MatchString(cleanPath) == true ||
		strings.HasPrefix(resolvedPath, docroot) == false ||
		os.IsPermission(err) == true {
		wslog.LogResponse(401, "Not Authorized", r.Method, r.URL, r.RemoteAddr, resolvedPath, "")
		http.Error(w, "Not Authorized", 401)
	} else if os.IsNotExist(err) == true {
		wslog.LogResponse(404, "Not Found", r.Method, r.URL, r.RemoteAddr, resolvedPath, "")
		http.NotFound(w, r)
	} else if err == nil {
		wslog.LogResponse(200, "OK", r.Method, r.URL, r.RemoteAddr, resolvedPath, "")
		http.ServeFile(w, r, resolvedPath)
	} else {
		// Easter egg
		wslog.LogResponse(418, "I'm a teapot", r.Method, r.URL, r.RemoteAddr, resolvedPath, "")
		http.Error(w, "I'm a teapot", 418)
	}
}
