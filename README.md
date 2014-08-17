ws
==

    A nimble webserver for prototyping. 


# What is _ws_?

*ws* is a simple webserver with optional support for JavaScript defined route handlers.  It supports basic _http_ and _https_ (SSL via TLS) protocols.  It does not support _SPDY_ or _http2_ protocols. It does not support PHP, Perl, Python, etc.

+ It is great for server static web pages!
    - Built on Golangs native http/https modules
    - Implements a restricted http.FileServer based Golang's builtin http.FileServer
+ It includes OttoEngine enabling JavaScript defined route handling
    - built on Robert Krimen's excellent [otto](https://github.com/robertkrimen/otto) JavaScript VM

Got an idea for a new project? Want to prototype it quickly? 

1. run "ws -init" to set things up
2. run ". etc/config.sh" seed your environment
3. run "ws" and start working!

_ws_ feature set has been kept minimal. Only what you need when you turn it on.

+ Restricted file service, only from the docroot and no "dot files" are served
+ No dynamic content support unless you turn on OttoEngine for JavaScript defined routes (great for creating JSON blobs used by a client side demo)
+ Quick startup, everything logged to console for easy debugging or piping to a log processor


## USAGE 

```
    ws [options]
```

## OPTIONS

	-D	(defaults to ) This is your document root for static files.
	-H	(defaults to localhost) Set this hostname for webserver.
	-O	(defaults to ) Turns on otto engine using the path for route JavaScript route handlers
	-P	(defaults to 8000) Set the port number to listen on.
	-cert	(defaults to ) path to your SSL cert pem file.
	-docroot	(defaults to ) This is your document root for static files.
	-h	(defaults to false) This help document.
	-help	(defaults to false) This help document.
	-host	(defaults to localhost) Set this hostname for webserver.
	-init	(defaults to false) Creates a basic project structure in the current working directory
	-key	(defaults to ) Path to your SSL key pem file.
	-keygen	(defaults to false) Interactive tool to generate TLS certificates and keys
	-o	(defaults to false) When true this option turns on ottoengine. Uses the path defined by WS_OTTO_PATH environment variable or one provided by -O option.
	-otto	(defaults to false) When true this option turns on ottoengine. Uses the path defined by WS_OTTO_PATH environment variable or one provided by -O option.
	-otto-path	(defaults to ) Turns on otto engine using the path for route JavaScript route handlers
	-port	(defaults to 8000) Set the port number to listen on.
	-tls	(defaults to false) When true this turns on TLS (https) support.
	-v	(defaults to false) Display the version number of ws command.
	-version	(defaults to false) Display the version number of ws command.


# Tour

## http support

Make sure _ws_ is your path. To run for basic _http_ service change to 
the directory you wish to serve and type _ws_ at the command prompt. Example -

```shell
    cd public_html
    ws
```

When _ws_ starts up you'll some configuration information and the URL that it is listening for. Notice the default port is 8000 so you need to include that part in your URL too. If your machine was named _example.local_ then you the URL might look like "http://localhost:8000". Point your web browser at the URL you see for your system.  When the web browser connections you should see a stream of log information. The stream of text will continue as long as you continue to have requests to the server until you shutdown _ws_. To shutdown _ws_ you can press the "ctrl" key and letter "c". This will kill the process and shutdown your web server.

You don't have to run _ws_ with the defaults.  You can specity a different document root with the _-docroot_ option. Here is an example of telling _ws_ to use the */www* directory for the document root.

```shell
    ws -docroot=/www
```

You can also pass these setting via your operating system's environment. Here is an example of configuration the above setting in a Bash script.


```shell
    WS_DOCROOT=/www
    ws
```

More typically you'll create a configuration file, source it (E.g. in your .profile or .bashrc setup files) and _ws_ will pickup the settings that way.

```bash
    #!/bin/bash
    export WS_DOCROOT=/www
```

If that was sourced in our login scripts then typing _ws_ on the command line would server the contents of */www* by default. You can override the defaults with the command line option _-docroot_.

In this way you can configure hostname, port.  In the following example
the port will be set to 8081 and the hostname will be "localhost".


### Command line version

```shell
    ws -docroot=/www -host=localhost -port=8081
```

### The envirnonment version

```bash
    #!/bin/bash
    export WS_DOCROOT=/www
    export WS_HOST=localhost
    export WS_PORT=8081
```

Source the script above then run _ws_.

```shell
    ws
```

If you have the environment variables set and use a command line option
the command line option will override the event variable setting. In the
example _ws_ will listen on port 8007.

```bash
    #!/bin/bash
    WS_DOCROOT=/www
    WS_HOST=localhost
    WS_PORT=8080
```

Now run _ws_

```
    ws -port=8007
```

For a full list of command line options run _ws_ with the _-help_ option.

The environment variables for _http_ service are

+ WS_DOCROOT
+ WS_HOST
+ WS_PORT


## https support

If you want to run with _https_ support it works on the same principles as _http_ support. It requires three additional pieces of information. 

1. It needs to knows where to find your *cert.pem*
2. It needs to know where to find your  *key.pem*
3. It needs to know to use SSL/TLS support.

By default _ws_ will look for *cert.pem* and *key.pem* in your *$HOME/etc/ws* directory. You can specify alternate locations with the _-cert_ and _-key_ command line options or the _WS\_CERT_ and _WS\_KEY_ environment variables.  To turn _https_ support on you need the option _-tls=true_ or the environment variable _WS\_TLS_ set to "true".


### Command line example

```bash
    ws -tls=true -cert=my-cert.pem -key=my-key.pem
```


### The environment version

```bash
    #!/bin/bash
    export WS_CERT=/etc/ws/cert.pem
    export WS_KEY=/etc/ws/key.pem
    export WS_TLS=true
```

If this was sourced in your login scripts then by default _ws_ will run as a 
_https_ server with the document root set to your current working directory for your current hostname on port 8443.


### Generating TLS certificates and keys

_ws_ comes with a *-keygen* option for generating self-signed certificates and keys.

```SHELL
    ws -keygen
```

This was create a *cert.pen* and *key.pem* files in *$HOME/etc/ws* directory.


## Otto

[otto](https://github.com/robertkrimen/otto) is a JavaScript virtual machine written by Robert Krimen.  The _ottoengine_ allows easy route oriented API prototyping.  Each JavaScript file rendered in the Otto virtual machine becomes a route.  E.g. *example-1.js* becomes the route */example-1*. *example-1* should contain a closure which can recieve a "Request" and "Response" object as parameters. The "Response" object is used to tell the web server what to send back to the browser requesting the route.

```JavaScript
    /* example-1.js - a simple example of Request and Response objects */
    (function (req, res) {
        var header = req.Header;

        res.setHeader("content-type", "text/html");
        res.setContent("<p>Here is the Header array received by this request</p>" +
            "<pre>" + JSON.stringify(header) + "</pre>");
    }(Request, Response));
```

Assuming _ottoengine_ is turned on then the page rendered should have a content type of "text/html" with the body shoulding the paragraph about exposing the request headers as a JSON blob.  Two command line options or environment variables turn _ottoengine_ on.

+ -otto, WS\_OTTO - values true/false, defaults to false. True turns on _ottoengine_
+ -otto-path, WS\_OTTO\_PATH - sets the path to the scripts used to defined the routes being handled. Each file found in the path becomes a route.

## LICENSE

copyright (c) 2014 All rights reserved.
Released under the [Simplified BSD License](http://opensource.org/licenses/bsd-license.php)

