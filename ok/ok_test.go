//
// Package ok provides a functions similar to those in NodeJS's assert module without catastrophic side effects.
//
// ok_test.go - tests for ok.go
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the Simplified BSD License.
// See: http://opensource.org/licenses/BSD-2-Clause
//
package ok

import (
	"testing"
)

func TestOk(t *testing.T) {
	if Ok(t, true, "Should should not prove fatal.") == true {
		t.Log("Ok true is OK!!")
	} else {
		t.Fatal("Ok true failed!")
	}
	if Ok(t, false, "This should fail.") == false {
		t.Log("Ok false should fail, we're actually OK here.")
	} else {
		t.Fatal("Ok false failed!!!! something went wrong!")
	}
}

func TestNotOk(t *testing.T) {
	if NotOk(t, false, "Should should not prove fatal.") == true {
		t.Log("NoOk false is OK!!")
	} else {
		t.Fatal("NotOk false failed!")
	}
	if NotOk(t, true, "This should fail.") == false {
		t.Log("NotOk true should fail, we're actually OK here.")
	} else {
		t.Fatal("Ok true failed!!!! something went wrong!")
	}
}
