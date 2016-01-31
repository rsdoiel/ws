//
// ws_test.go test routines for ws.go
//
// Copyright (c) 2014 - 2016, R. S. Doiel
// All rights not granted herein are expressly reserved by R. S. Doiel.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package ws

import (
	"os"
	"testing"
)

func TestReadJSFiles(t *testing.T) {
	jsSources, err := ReadJSFiles("jsdocs")
	if err != nil {
		t.Errorf("reading jsdocs %s", err)
	}
	if len(jsSources) < 1 {
		t.Error("should find more than one JavaScript file in jsdocs")
	}
	jsFilename := "jsdocs/helloworld.js"
	jsSrc, ok := jsSources[jsFilename]
	if ok == false {
		t.Errorf("Should find %s in jsSources", jsFilename)
	}
	if jsSrc == nil {
		t.Errorf("Should find source code for %s", jsFilename)
	}
}

func TestJSEngine(t *testing.T) {
	vm := NewJSEngine()
	if vm == nil {
		t.Errorf("should have created a new JavaScript VM")
	}

	val, err := vm.Eval(`
    Getenv != undefined
  `)
	b, err := val.ToBoolean()
	if err != nil {
		t.Errorf("Error from Getenv != undefined, %s", err)
	}
	if b == false {
		t.Errorf("Expected Getenv != undefined to return true, %b", b)
	}
	val, err = vm.Eval(`
    HttpGet != undefined
  `)
	b, err = val.ToBoolean()
	if err != nil {
		t.Errorf("Error from HttpGet != undefined, %s", err)
	}
	if b == false {
		t.Errorf("Expected HttpGet != undefined to return true, %b", b)
	}
	val, err = vm.Eval(`
    HttpPost != undefined
  `)
	b, err = val.ToBoolean()
	if err != nil {
		t.Errorf("Error from HttpPost != undefined, %s", err)
	}
	if b == false {
		t.Errorf("Expected HttpPost != undefined to return true, %b", b)
	}

	err = os.Setenv("WS_TEST", "hello world")
	val, err = vm.Eval(`
    s = Getenv("WS_TEST");
    s;
  `)
	if err != nil {
		t.Errorf(`vm.Eval() error, %s`, err)
	}
	s, err := val.Export()
	if err != nil {
		t.Errorf("Can't export 's' from JS eval, %s", err)
	}
	if s != "hello world" {
		t.Errorf("Expected JS to return s of 'hello world', got %s", s)
	}
}
