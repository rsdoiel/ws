ws
==

## A nimble web server

_ws_ is a prototyping platform for web based services and websites.

_ws_ started as a nimble static content web server.  It now includes two
friends - _wsinit_ and _wsindexer_.  The first setups a project
directory structures, creates self signed SSL keys or displays suggested
environment variables whilethe second creates and updates
[bleve](http://blevesearch.com) based indexes for use with _ws_.

## Requirements

+ [Golang](http://golang.org) version 1.5.3 or better
+ 3rd Party Go packages
  + [Otto](https://github.com/robertkrimen/otto) by Robert Krimen, MIT license
    + a JavaScript engine for Go
  + [Bleve](https://github.com/blevesearch/bleve) by [Blevesearch](http://blevesearch.com), Apache License, Version 2.0
    + think lucene-lite for Go

## Compile

Here's my basic approach to get things setup. I assume you've got *Golang* already installed.

```
  git clone https://github.com/rsdoiel/ws
  cd ws
  go get -u github.com/robertkrimen/otto
  go get -u github.com/blevesearch/bleve
  go test
  go build
  go build cmds/ws/ws.go
  go build cmds/wsinit/wsinit.go
  go build cmds/wsindexer/wsindexer.go
```

If everything compiles fine then I do something like this--

```
  go install cmds/ws/ws.go
  go install cmds/wsinit/wsinit.go
  go install cmds/wsindexer/wsindexer.go
```


### _ws_ features

+ http/https protocols
+ static file server
+ a simplified server side JavaScript platform
  + if you need more, I suggest [NodeJS](http://nodejs.org)
+ built-in engine via Bleve

The server side JavaScript environment also supports template rendering via Golang's html/template system.

### _wsinit_ features

_wsinit_ takes three actions

+ create a basic site directory structure (e.g. htdocs, jsdocs, etc)
+ create self signed SSL certificates (e.g. etc/site.key, etc/site.pem)
+ suggests environment variable settings (e.g like you might put in etc/setup.conf)

### _wsindexer_ features

_wsindexer_ does one of two things - create a Bleve index or update a Bleve index


## Configuration

All configuration for _ws_ and friends can be configured via environment
variables. That makes them container friendly.  The environment can be
overwritten by command line options.

The standard set of environment variables are

+ WS_URL the URL to listen for by _ws_
  + default is http://localhost:8000
+ WS_HTDOCS the directory of your static content you need to serve
  + the default is ./htdocs
+ WS_JSDOCS the directory for any server side JavaScript processing
  + the default is ./jsdocs (if not found then server side JavaScript is turned off)
+ WS_SSL_KEY the path the the SSL key file
  + default is empty, only checked if your WS_URL is starts with https://
+ WS_SSL_PEM the path the the SSL pem file
  + default is empty, only checked if your WS_URL is starts with https://
+ WS_BLEVE_INDEX sets the path for the BLEVE index and turns on the search features
  + default is empty, search is only enabled if this is defined

### command line options

+ -h, --help displays the help documentation
+ -url overrides WS_URL
+ -htdocs overrides WS_HTDOCS
+ -jsdocs overrides WS_JSDOCS
+ -ssl-key overrides WS_SSL_KEY
+ -ssl-pem overrides WS_SSL_PEM
+ -bleve-index overrides WS_BLEVE_INDEX
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

Some additional objects are provided to facilitate server side
JavaScript development--

+ http
  + http.Get(url, array_of_headers) which performs a HTTP GET
  + http.Post(url, array_of_headers, payload) which performs an HTTP POST
+ os
  + os.Getenv(varname) which will read an environment variable
+ site
  + site.Search(query, all, exact, excluded) a simple search service for your site


## LICENSE

copyright (c) 2014 - 2016 All rights reserved.
Released under the [Simplified BSD License](http://opensource.org/licenses/bsd-license.php)
