
    A nimble webserver for prototyping. 


# What is _ws_?

The project started in 2014 after having setup another instance Apache just to show something to a colleague.
Have been playing around with NodeJS and Golang for building RESTful services it just seemed that Apache was
overkill. I just wanted to type "webserver" and have an ephemerial webserver instance running rather than
creating another virtualhost.  That is the itch that _ws_ tries to skratch.  It is not intended to be a full
featured webserver. It is designed to be simple to start from the command line, configurable via the environment
(inspired by 12 factor apps) and have the minimum of functionality to do a quick prototype, static site or 
API mockup.

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

_ws_ comes with *wskeygen* for generating self-signed certificates and keys.

```SHELL
    wskeygen
```

This was create a *cert.pem* and *key.pem* files in *$HOME/etc/ws* directory.


### Generate a project folder and certificates

_ws_ comes with _wsinit_ for interactively generating a project tree and certificates.

```SHELL
    wsinit
```


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
