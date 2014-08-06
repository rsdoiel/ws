/**
 * fsengine.go - a simple restricted FileServer Handler for
 * ws.go.  Avoids exposing dot files via http requests.
 */
package fsengine

import (
	"../app"
	"../wslog"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

// This is a Restricted FileService excluding dot files and directories.
func Engine(cfg *app.Cfg, w http.ResponseWriter, r *http.Request) {
	var (
		hasDotPath = regexp.MustCompile(`\/\.`)
		docroot    = cfg.Docroot
	)

	unclean_path := r.URL.Path
	if !strings.HasPrefix(unclean_path, "/") {
		unclean_path = "/" + unclean_path
	}
	clean_path := path.Clean(unclean_path)
	r.URL.Path = clean_path
	resolved_path := path.Clean(path.Join(docroot, clean_path))
	_, err := os.Stat(resolved_path)
	if hasDotPath.MatchString(clean_path) == true ||
		strings.HasPrefix(resolved_path, docroot) == false ||
		os.IsPermission(err) == true {
		wslog.LogResponse(401, "Not Authorized", r.Method, r.URL, r.RemoteAddr, resolved_path, "")
		http.Error(w, "Not Authorized", 401)
	} else if os.IsNotExist(err) == true {
		wslog.LogResponse(404, "Not Found", r.Method, r.URL, r.RemoteAddr, resolved_path, "")
		http.NotFound(w, r)
	} else if err == nil {
		wslog.LogResponse(200, "OK", r.Method, r.URL, r.RemoteAddr, resolved_path, "")
		http.ServeFile(w, r, resolved_path)
	} else {
		// Easter egg
		wslog.LogResponse(418, "I'm a teapot", r.Method, r.URL, r.RemoteAddr, resolved_path, "")
		http.Error(w, "I'm a teapot", 418)
	}
}
