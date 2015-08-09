package example

import (
    "testing"
    "../../ok"
)

func TestHello (t *testing.T) {
   ok.Ok(t, Hello("George") == "Hello George", "Should say, Hello George")
   ok.NotOk(t, Hello("Fred") == "Hello George", "Should not say, Hello George")
}
