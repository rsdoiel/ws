
[![Go Report Card](http://goreportcard.com/badge/rsdoiel/ws)](http://goreportcard.com/report/rsdoiel/ws)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

# ws

## A nimble web server

_ws_ is a prototyping platform for web based services and websites.

_ws_ started as a nimble static content web server.  It now includes
support for server side JavaScript and can initialize a basic project
layout. The *init* option generates a project directory
structure, creates self signed SSL keys and displays suggested environment
variables for use with _ws_.  If the directory structure already exists it
will display the suggest setup.

## Requirements

+ [Golang](http://golang.org) version 1.5.3 or better
+ A 3rd Party Go package
  + [Otto](https://github.com/robertkrimen/otto) by Robert Krimen, MIT license
    + a JavaScript interpreter written in Go

## Compile

Here's my basic approach to get things setup. I assume you've got *Golang* already installed.

+ [ws](README.md) a nimple webserver for static content
    + Built on Golangs native http/https modules
    + Restricted file service, only from the docroot and no "dot files" are served
    + No dynamic content support 
    + Quick startup, everything logged to console for easy debugging or piping to a log processor
+ [ws js support](README.md) a nimple webserver for static and JavaScript generated dynamic content
    + built on Robert Krimen's excellent [otto](https://github.com/robertkrimen/otto) JavaScript VM
    + Mockup your dynamic content via JavaScript defined routes (great for creating JSON blobs used by a browser side demo)

```
  git clone https://github.com/rsdoiel/ws
  cd ws
  go get -u github.com/robertkrimen/otto
  go test
  go build
  go build cmds/ws/ws.go
```

If everything compiles fine then I do something like this--

```
  go install cmds/ws/ws.go
```


### _ws_ features

+ http/https protocols
+ static file server
+ a simplified server side JavaScript platform
  + if you need more, check out [NodeJS](http://nodejs.org)
+ a project setup option called *init*

*_ws_ init* takes three actions

+ create a basic site directory structure (e.g. htdocs, jsdocs, etc) if needed
+ creates self signed SSL certificates (e.g. etc/site.key, etc/site.pem) if appropriate
+ suggests environment variable settings (e.g like you might put in etc/setup.conf)


## Configuration

Configuration for _ws_  can be passed directly from environment
variables. That makes them container friendly.  The environment can be
overwritten by command line options.

The standard set of environment variables are

+ WS_URL the URL to listen for by _ws_
  + default is http://localhost:8000
+ WS_HTDOCS the directory of your static content you need to serve
  + the default is ./htdocs
+ WS_JSDOCS the directory for any server side JavaScript processing
  + the default is ./jsdocs (if not found then server side JavaScript is turned off)
+ WS_SSL_KEY the path the the SSL key file (e.g. etc/ssl/site.key)
  + default is empty, only checked if your WS_URL is starts with https://
+ WS_SSL_CERT the path the the SSL cert file (e.g. etc/ssl/site.crt)
  + default is empty, only checked if your WS_URL is starts with https://

### command line options

+ -url overrides WS_URL
+ -htdocs overrides WS_HTDOCS
+ -jsdocs overrides WS_JSDOCS
+ -ssl-key overrides WS_SSL_KEY
+ -ssl-pem overrides WS_SSL_PEM
+ -init triggers the initialization process
+ -h, --help displays the help documentation
+ -v, --version display the version of the command


## A word about the Server Side JavaScript implementation

[otto](https://github.com/robertkrimen/otto) is a JavaScript virtual machine
written by Robert Krimen. It is the engine that powers _ws_ server side
JavaScript capabilities. Each JavaScript file in the *jsdocs* directory tree
becomes a URL end point or route. E.g. *jsdocs/example-1.js* becomes the
route */example-1*. *example-1*. Each of the server side JavaScript files
should contain a closure accepting a "Request" and "Response" object as
parameters.  E.g.

```JavaScript
    /* example-1.js - a simple example of Request and Response objects */
    (function (req, res) {
        var header = req.Header;

        res.setHeader("content-type", "text/html");
        res.setContent(
          "<p>Here is the Header array received by this request</p>" +
          "<pre>" + JSON.stringify(header) + "</pre>");
    }(Request, Response));
```

Assuming server side JavaScript is enabled then the page rendered should have a
content type of "text/html" with the body should be holding the paragraph and
pre element.

Some additional functions are provided to facilitate server side
JavaScript development--

+ http related
  + WS.httpGet(url, array_of_headers) which performs a HTTP GET
  + WS.httpPost(url, array_of_headers, payload) which performs an HTTP POST
+ os related
  + WS.getEnv(varname) which will read an environment variable

Want to expand out the site quickly, write the HTML skeleton with markdown, sprinkle in some [shorthand](http://rsdoiel.github.io/shorthand) which can leverage some shell logic and now you have HTML pages with common nav, headers, and footers.


## Installation

_ws_ can be installed with the *go get* command.

```
    go get github.com/rsdoiel/ws/...
```

