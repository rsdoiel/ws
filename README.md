ws
==

A light weight webserver suitable for static web content development. 

*ws* supports basic _http_ and _https_ (via TLS) protocols.  It does not support _SPDY_ or _http2_ protocols. It does not support PHP, Perl, Python, etc. It is a static content webserver.

# USAGE

## http support

Make sure _ws_ is your path. To run for basic _http_ service change to 
the directory you wish to serve and type _ws_ at the command prompt. 

```shell
    cd public_html
    ws
```

You see a stream of log information with the webserver starting and about
any requests received. The default port is 8000.

You can specity a different document root with the _-docroot_ option. Here
is an example of telling _ws_ to use the */www* directory for the document
root.

```shell
    ws -docroot=/www
```

You can also pass these setting via your operating system's environment. Here is an example of configuration the above setting in a Bash script.


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
3. It needs to know to use TLS support.

By default _ws_ will look for *cert.pem* and *key.pem* in your *$HOME/etc/ws* directory. You can specify alternate locations with the _-cert_ and _-key_ command line options or the _WS\_CERT_ and _WS\_KEY_ environment variables.

To turn _https_ support on you need the option _-tls=true_ or the environment variable _WS\_TLS_ set to "true".

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

If this was sourced in your login scripts then by default _ws_ will run as a _https_ server with the document root set to your
current working directory for your current hostname on port 8443.
