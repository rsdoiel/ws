ws
==

A light weight webserver suitable for static web content development. 

*ws* supports basic _http_ and _https_ (SSL via TLS) protocols.  It does not 
support _SPDY_ or _http2_ protocols. It does not support PHP, Perl, Python, 
etc. It is a static content webserver.

It has experimental support for dynamic content via route handlers written in JavaScript
and execute by the [otto](https://github.com/robertkermin/otto) JavaScript
virtual machine.


# USAGE

## http support

Make sure _ws_ is your path. To run for basic _http_ service change to 
the directory you wish to serve and type _ws_ at the command prompt. Example -

```shell
    cd public_html
    ws
```

When _ws_ starts up you'll some configuration information and the URL that it 
is listening for. Notice the default port is 8000 so you need to include that 
part in your URL too. If your machine was named _example.local_ then you the URL 
might look like "http://example.local:8000". Point your web browser at the URL 
you see for your system.  When the web browser connections you should see a 
stream of log information. The stream of text will continue as long as you 
continue to have requests to the server until you shutdown _ws_. To shutdown
_ws_ you can press the "ctrl" key and letter "c". This will kill the process
and shutdown your web server.

You don't have to run _ws_ with the defaults.  You can specity a different 
document root with the _-docroot_ option. Here is an example of telling _ws_ to 
use the */www* directory for the document
root.

```shell
    ws -docroot=/www
```

You can also pass these setting via your operating system's environment. Here is 
an example of configuration the above setting in a Bash script.


```shell
    WS_DOCROOT=/www
    ws
```

More typically you'll create a configuration file, source it (E.g. in your
.profile or .bashrc setup files) and _ws_ will pickup the settings that
way.

```bash
    #!/bin/bash
    export WS_DOCROOT=/www
```

If that was sourced in our login scripts then typing _ws_ on the command
line would server the contents of */www* by default. You can override the
defaults with the command line option _-docroot_.



In this way you can configure hostname, port.  In the following example
the port will be set to 8081 and the hostname will be "example.com".

### Command line version

```shell
    ws -docroot=/www -host=example.com -port=8081
```

### The envirnonment version

```bash
    #!/bin/bash
    export WS_DOCROOT=/www
    export WS_HOST=example.com
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
    WS_HOST=example.com
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

If you want to run with _https_ support it works on the same 
principles as _http_ support. It requires three additional pieces
of information. 

1. It needs to knows where to find your *cert.pem*
2. It needs to know where to find your  *key.pem*
3. It needs to know to use SSL/TLS support.

By default _ws_ will look for *cert.pem* and *key.pem* in your *$HOME/etc/ws* 
directory. You can specify alternate locations with the _-cert_ and _-key_ 
command line options or the _WS\_CERT_ and _WS\_KEY_ environment variables.

To turn _https_ support on you need the option _-tls=true_ or the environment 
variable _WS\_TLS_ set to "true".

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
_https_ server with the document root set to your current working directory 
for your current hostname on port 8443.

### Generating TLS certificates and keys

_ws_ comes with a *-keygen* option for generating self-signed certificates
and keys.

```SHELL
    ws -keygen
```

This was create a *cert.pen* and *key.pem* files in *$HOME/etc/ws* directory.

## Otto

[otto](https://github.com/robertkrimen/otto) is a JavaScript virtual machine written by Robert Krimen.
The _ottoengine_ is an experimental route handler that uses _otto_ to render route content dynamically.
The goal of _ottoengine_ is to provide a platform for prototyping content APIs consumed browser side
by the static pages served by _ws_.

Enabled the _ottoengine_ involes setting to command line options or settin the equivalant environment
variables

+ -otto, WS_OTTO - values true/false, defaults to false.
+ -otto-path, WS_OTTO_PATH - sets the path to the scripts used to defined the routes being handled.


