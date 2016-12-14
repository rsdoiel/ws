ws
==

    A nimble webserver with friends for prototyping. 

# What is _ws_?

The project started in 2014 after having setup another instance Apache just to show something to a colleague.  During the same period I was busy developing web services with NodeJS and Golang. After a time Apache just seemed like overkill. I wanted to type "webserver" and have an ephemerial webserver instance running rather than creating another virtualhost.  That is the itch that _ws_ tries to scratch.  It is not intended to be a full featured webserver. It is designed to be simple to start from the command line, configurable via the environment (inspired by 12 factor apps) and have the minimum of functionality to do a quick prototype, serve a static site or API mockup.

## Use cases

### a basic httpd

Make sure _ws_ and friends are in your path. To run for basic _httpd_ service change to the directory you wish to serve and type _ws_ at the command prompt. Example -

```shell
    cd public_html
    ws
```

When _ws_ starts up you'll some configuration information and the URL that it is listening for. Notice the hostname is *localhost* and port is *8000*.  Either can be configured either via command line options (e.g. -H and -p) or through environment variables (e.g. WS_HOST and WS_PORT). By default the root document directory will be your current work directory. Again this can be configure via the command line or environment variables (e.g. -d and WS_DOCROOT). Log messages are display to the console and to stop the webserver you can press Control-C or use the Unix *kill* command and the process id.

```shell
    ws --help
```


Here is an example of using _ws_ to server the document root of */www*.

```shell
    ws -docs=/www
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
    export WS_URL="http://me.example.org:80"
```

This would have _ws_ listen for http://me.example.org request on the default http port of 80. Note that on most system you'll your account will need special privilleges to run on port 80.

Here is the equivallent using only the command line.

```shell
    ws -d /www -u http://me.example.org:80
```

The long option name version.


```shell
    ws -docs=/www -url=http://me.example.org:80
```

The environment variables for _http_ service are

+ WS_DOCROOT
+ WS_URL


## https support

If you want to run with _https_ support it works on the same principles as _http_ support. _ws_ requires three additional pieces of information. 

1. It needs to knows where to find your *cert.pem*
2. It needs to know where to find your  *key.pem*
3. You use the "https:" instead of "http:" in WS_URL or -url options

By _ws_ will look environment variables for *cert.pem* and *key.pem* in the values of _WS_SSL_KEY_ and _WS_SSL_CERT_. You can also specify where to find them via the _-key_ and _-cert_ command line option.


### Command line example

```bash
    ws -cert=etc/ssl/my-cert.pem -key=etc/ssl/my-key.pem -url=https://me.example.org:443 -htdocs=/www
```

### The environment version

```bash
    #!/bin/bash
    export WS_HTDOCS=/www
    export WS_URL="https://me.example.org:443"
    export WS_CERT=/etc/ssl/cert.pem
    export WS_KEY=/etc/ssl/key.pem
```

Like the *http* example running on port 443 will likely require a privilleged account.

