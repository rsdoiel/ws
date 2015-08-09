# ok

A small collection of assertion like functions for use with Go's testing package. 
E.g. they use testing.T rather than explicit calls to log.Fatal() or os.Exit().

## example usage

Pseudo go code to illustrate how you might use ok.Ok() and ok.NotOk().

```go
    package hello

    imports (
        "testing"
        "github.com/rsdoiel/ok"
    )
    
    // your test code...
    
    func TestHello(t *testing.T) {
        ok.Ok(t, Hello("George") == "Hello George", "Should say Hello George.")
        ok.NoOk(t, Hello("Fred") == "Hello George", "Should not say Hello George for Fred.")
    }
```

## testing ok

If all goes well you should see output similar to

```shell
    localhost:ok you$ go test
    --- FAIL: TestOk (0.00s)
        ok_test.go:19: Ok true is OK!!
        ok.go:21: Failed (expected true): [This should fail.]
        ok_test.go:24: Ok false should fail, we're actually OK here.
    --- FAIL: TestNotOk (0.00s)
        ok_test.go:32: NoOk false is OK!!
        ok.go:29: Failed (expected false): [This should fail.]
        ok_test.go:37: NotOk true should fail, we're actually OK here.
    FAIL
    exit status 1
    FAIL    _/home/you/git-repos/ok    0.011s
    localhost:ok you$ vi README.md 
```

