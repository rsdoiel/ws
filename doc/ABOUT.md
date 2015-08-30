ws
==

    A nimble webserver with friends for prototyping. 

# What is _ws_?

The project started in 2014 after having setup another instance Apache just to show something to a colleague.
Have been playing around with NodeJS and Golang for building RESTful services it just seemed that Apache was
overkill. I just wanted to type "webserver" and have an ephemerial webserver instance running rather than
creating another virtualhost.  That is the itch that _ws_ tries to scratch.  It is not intended to be a full
featured webserver. It is designed to be simple to start from the command line, configurable via the environment
(inspired by 12 factor apps) and have the minimum of functionality to do a quick prototype, static site or 
API mockup.

## Use cases

### a basic httpd

Make sure _ws_ and friends are in your path. To run for basic _httpd_ service change to the directory you wish to serve and type _ws_ at the command prompt. Example -

```shell
    cd public_html
    ws
```

When _ws_ starts up you'll some configuration information and the URL that it is listening for. Notice the hostname is *localhost* and port is *8000*.  Either can be configured either via command line options (e.g. -H and -p) or through environment variables (e.g. WS_HOST and WS_PORT). By default the root document directory will be your current work directory. Again this can be configure via the command line or environment variables (e.g. -d and WS_DOCROOT). Log messages are display to the console and to stop the webserver you can press Control-C or use the Unix *kill* command and the process id.

Getting a list of all the option that _ws_ (or _wsjs_) supports use the command line options of "-h" or "--help". Most options have a long and short form.

```shell
    ws --help
```


Here is an example of using _ws_ to server the document root of */www*.

```shell
    ws -docroot=/www
```

Notice that we used the long form of the option in this case. It works the same way of the environment variable and "-d".  Here is an example of configuration the above setting in a Bash script.


```shell
    export WS_DOCROOT=/www
    ws
```

It is easy to use Bash files as configuration for _ws_. Just source your file with the settings then type _ws_ at the command line. Example Bash configuration


```bash
    #!/bin/bash
    export WS_DOCROOT=/www
    export WS_HOST="me.example.org"
    export WS_PORT="80"
```

This would have _ws_ listen for http://me.example.org request on the default http port of 80. Note that on most system you'll your account will need special
privilleges to run on port 80.

Here is the equivallent using only the command line.

```shell
    ws -d /www -H me.example.org -p 80
```

The long option name version.


```shell
    ws -docroot=/www -host=me.example.org -port=80
```

The environment variables for _http_ service are

+ WS_DOCROOT
+ WS_HOST
+ WS_PORT


## https support

If you want to run with _https_ support it works on the same principles as _http_ support. _ws_ requires three additional pieces of information. 

1. It needs to knows where to find your *cert.pem*
2. It needs to know where to find your  *key.pem*
3. It needs to know to use SSL/TLS support.

By default _ws_ will look for *cert.pem* and *key.pem* in your *$HOME/etc/ws* directory. You can specify alternate locations with the _-cert_ and _-key_ command line options or the _WS_CERT_ and _WS_KEY_ environment variables.  To turn _https_ support on you need the option _-tls=true_ or the environment variable _WS_TLS_ set to "true".


### Command line example

```bash
    ws -tls=true -cert=etc/ssl/my-cert.pem -key=etc/ssl/my-key.pem -host=me.example.org -port=443 -docroot=/www
```

### The environment version

```bash
    #!/bin/bash
    export WS_DOCROOT=/www
    export WS_HOST="me.example.org"
    export WS_PORT="443"
    export WS_CERT=/etc/ssl/cert.pem
    export WS_KEY=/etc/ssl/key.pem
    export WS_TLS=true
```

Like the *http* example running on port 443 will likely require a privilleged account.


## Generating TLS certificates and keys

_ws_ comes with _wskeygen_ for generating self-signed certificates and keys.

```shell
    wskeygen
```

This is an interactive problem. It will prompt for information about where to store the certificates. Alternatively if you want to also setup a basic directory structure you can use the _wsinit_ command which will include the option of generating (or replacing) the desired SSL certificates. Both are interactive.

```shell
    wsinit
```

### Generate a project folder and certificates

_ws_ comes with _wsinit_ for interactively generating a project tree and certificates.

```SHELL
    wsinit
```

## Dynamic content via server side JavaScript

_wsjs_ provides a mechanism to define routes using JavaScript files where the JavaScript files can process teh http/https request and populate the http/https response. It relies on Robert Krimen's otto for enterpretting JavaScript.

### Otto

[otto](https://github.com/robertkrimen/otto) is a JavaScript virtual machine written by Robert Krimen.  The _ottoengine_ allows easy route oriented API prototyping.  Each JavaScript file rendered in the Otto virtual machine becomes a route.  E.g. *example-1.js* becomes the route */example-1*. *example-1* should contain a closure which can recieve a "Request" and "Response" object as parameters. The "Response" object is used to tell the web server what to send back to the browser requesting the route.

### example JavaScript route definition

If this file is in the directory defined by the environment variable WS_OTTO_PATH and is named *example-1.js" then the URL requesting (e.g. GET, POST, PUT, DELETE) of /example-1 will run receive the response based on the following JavaScript.

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

+ -otto-path, WS_OTTO_PATH - sets the path to the scripts used to defined the routes being handled. This path turns on otto engine support. Each file found in the path becomes a route.

That's the basics of _ws_ and _wsjs_ nimble webservers.


## LICENSE

copyright (c) 2014 All rights reserved.
Released under the [Simplified BSD License](http://opensource.org/licenses/bsd-license.php)

