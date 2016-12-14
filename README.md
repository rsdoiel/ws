
[![Go Report Card](http://goreportcard.com/badge/rsdoiel/ws)](http://goreportcard.com/report/rsdoiel/ws)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

# ws

## A nimble web server

_ws_ is a nimble static content web server.  That's it, no bells,
no whistles. It can serve http and https version 1 and 2 websites
as it is a minimal solution built on Go's http/https packages.
By default it serves the content out your current working directory
as http://localhost:8000.  It will NOT serve out content in directories
beginning with a period (e.g. .ssh, .git).

## Requirements

+ [Golang](http://golang.org) version 1.7.3 or better

## Setup

_ws_ is "go gettable"

```shell
    go get github.com/rsdoiel/ws/...
```

### compiling from source

Here's my basic approach to get things setup. I assume you've got *Golang* already installed.

+ [ws](README.md) a nimple webserver for static content
    + Built on Golangs native http/https modules
    + Restricted file service, only from the docroot and no "dot files" are served
    + No dynamic content support 
    + Quick startup, everything logged to console for easy debugging or piping to a log processor

```
  git clone https://github.com/rsdoiel/ws $HOME/src/github.com/rsdoiel/ws
  cd $HOME/src/github.com/rsdoiel/ws
  go test
  go build
  go build cmds/ws/ws.go
```

If everything compiles fine then I do something like this--

```
  go install cmds/ws/ws.go
```

### Compiled versions

You can also grab one of the precompiled version at https://github.com/rsdoiel/ws/releases/latest


## features

+ http/https protocols
+ static file server

## Configuration

Configuration for _ws_  can be passed directly from environment
variables. That makes them container friendly.  The environment can be
overwritten by command line options.

The standard set of environment variables are

+ WS_URL the URL to listen for by _ws_
  + default is http://localhost:8000
+ WS_DOCROOT the directory of your static content you need to serve
  + the default is . (your current directory)
+ WS_SSL_KEY the path the the SSL key file (e.g. etc/ssl/site.key)
  + default is empty, only checked if your WS_URL is starts with https://
+ WS_SSL_CERT the path the the SSL cert file (e.g. etc/ssl/site.crt)
  + default is empty, only checked if your WS_URL is starts with https://

### command line options

+ -u, -url overrides WS_URL
+ -d, -docs overrides WS_DOCROOT
+ -k, -key overrides WS_SSL_KEY
+ -c, -cert overrides WS_SSL_CERT
+ -h, -help displays the help documentation
+ -v, -version display the version of the command
+ -l, -license display software license


