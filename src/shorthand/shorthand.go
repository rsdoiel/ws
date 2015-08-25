/**
 * shorthand.go - A simple definition and expansion notation to use
 * as shorthand when a template language is too much.
 *
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 * copyright (c) 2015 all rights reserved.
 * Released under the BSD 2-Clause license
 * See: http://opensource.org/licenses/BSD-2-Clause
 */

// Package shorthand provides shorthand definition and expansion for ws and stngo projects.
package shorthand

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

// Abbrevations holds the shorthand and translation
var Abbreviations = make(map[string]string)

// IsAssignment checks to see if a string contains an assignment (i.e. has a ' := ' in the string.)
func IsAssignment(text string) bool {
	if strings.Index(text, " := ") == -1 &&
		strings.Index(text, " :< ") == -1 &&
		strings.Index(text, " :! ") == -1 &&
		strings.Index(text, " :{ ") == -1 {
		return false
	}
	return true
}

// HasAssignment checks to see if a shortcut has already been assigned.
func HasAssignment(key string) bool {
	_, ok := Abbreviations[key]
	return ok
}

// Assign stores a shorthand and its expansion
func Assign(s string) bool {
	var parts []string
	if strings.Index(s, " :{ ") != -1 {
		parts = strings.SplitN(strings.TrimSpace(s), " :{ ", 2)
		parts[1] = Expand(parts[1])
	} else if strings.Index(s, " := ") != -1 {
		parts = strings.SplitN(strings.TrimSpace(s), " := ", 2)
	} else if strings.Index(s, " :< ") != -1 {
		parts = strings.SplitN(strings.TrimSpace(s), " :< ", 2)
		buf, err := ioutil.ReadFile(parts[1])
		if err != nil {
			log.Fatalf("Cannot read %s: %v\n", parts[1], err)
		}
		parts[1] = string(buf)
	} else if strings.Index(s, " :! ") != -1 {
		parts = strings.SplitN(strings.TrimSpace(s), " :! ", 2)
		buf, err := exec.Command("bash", "-c", parts[1]).Output()
		if err != nil {
			log.Fatal(err)
		}
		parts[1] = string(buf)
	} else {
		log.Fatalf("[%s] is an invalid assignment.\n", s)
	}
	key, value := parts[0], parts[1]
	if key == "" || value == "" {
		return false
	}
	Abbreviations[key] = value
	_, ok := Abbreviations[key]
	return ok
}

// Expand takes a text and expands all shorthands
func Expand(text string) string {
	// Iterate through the list of key/values in abbreviations
	for key, value := range Abbreviations {
		text = strings.Replace(text, key, value, -1)
	}
	return text
}

// Clear remove all the elements of a map.
func Clear() {
	for key := range Abbreviations {
		delete(Abbreviations, key)
	}
}
