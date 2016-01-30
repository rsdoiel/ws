//
// Package prompt is a simple message and response functions. It is used by
// ws.go when the generate config option is chosen.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the Simplified BSD License.
// See: http://opensource.org/licenses/BSD-2-Clause
//
package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// YesNo displays a message and return true for y, Y and return false for n, N.
func YesNo(msg string) bool {
	buf := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n) ", msg)
	yesOrNo, _ := buf.ReadString('\n')
	yesOrNo = strings.ToLower(yesOrNo[0:1])
	if yesOrNo == "y" {
		return true
	}
	return false
}

// Question displays a message and prompt for a String response.
func Question(msg string, defaultResponse string) string {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	response, _ := buf.ReadString('\n')
	response = strings.TrimSpace(response)
	if response == "" {
		return defaultResponse
	}
	return response

}
