// prompt.go - Simple message and response functions.
package prompt

import (
    "os"
    "bufio"
    "strings"
    "fmt"
)


// Display a message and return true for y, Y and return false for n, N.
func YesNo (msg string) bool {
    buf := bufio.NewReader(os.Stdin)
    fmt.Printf("%s (y/n) ", msg)
    yes_or_no, _ := buf.ReadString('\n')
    yes_or_no = strings.ToLower(yes_or_no[0:1])
    if yes_or_no  == "y" {
        return true
    }
    return false
}

// Display a message and prompt for a String response.
func PromptString (msg string, default_response string) string {
    buf := bufio.NewReader(os.Stdin)
    fmt.Print(msg)
    response, _ := buf.ReadString('\n')
    response = strings.TrimSpace(response)
    if response == "" {
        return default_response
    }
    return response
    
}

