
#  Getting started with _ws_ 

At its simplest _ws_ is a static content web server.  It makes it quick to prototype things
that run browser side.  After you have installed _ws_ server static content is as easy as changing
your directory to your document root and then starting _ws_. 


## Example 1

### Simple usage

You have create a directory called /Sites where you plan to develop your website.  To test
your site with _ws_ you need to--

1. Change to the /Sites directory
2. Start _ws_


```shell
    cd /Sites
    ws .
```

This should yield output like

```shell

                   TLS: false
                  Cert: 
                   Key: 
               Docroot: /Sites
                  Host: localhost
                  Port: 8000
                Run as: johndoe

         2014/07/15 17:52:27 Starting http://localhost:8000
```

You can now point your browser a [http://localhost:8000](http://localhost:8000) and see the contents
of the /Sites directory.

You can press ctrl-C (while holding the key marked "Ctr" or "Ctrl" press the "c" key).  The websere should now stop.


## Example 2

### Organizing and doing more

More typically if you are prototyping a website you will organize your code into different folders
based on your build process or tool set.  _wsinit_ can help here.  The *wsinit* will configure
a folder with a simple structure for further develop.  It also will setup up things for more complex
usage of _ws -init_ like accessing Otto Engine or running under SSL.

Here are the normal four steps you take to set things up. We will do the first two, stop, look around
then proceed to steps 3 and 4 to test it out.

1. Change to the directory that will hold your project (e.g. /Sites)
2. Run _ws -init_. For not accept the defaults by pressing enter you "y" and enter when prompted.
3. Source your _etc/config.sh_ file
4. Start _ws_ webserver and test with your browser.

Steps one and two.

```shell
    cd /Sites
    ws -init
```

Take a look at the directories and files created. By default your static content is configured to run
from the _static_ directory. You will find a new _index.html_ created there for you to modify. There
is also a _dynamic_ directory.  You will find a single file, _test.js_, there. Any Javascript file
in this directory will server is a dynamic route handler by _ws_.  This makes it easy to mockup 
web forms and responses or even mockup simple API services.

There is a directory called _etc_ too.  This is where you will find your configuration file
_config.sh_ as well as another sub-directory, _etc/ssl_, holding your SSL certificate and key files.

_ws_ draws configuration from either command line options or your shells environment variables.
If you source _etc/config.sh_ you will see anumber of environment vararibles set so _ws_ know
where to find your static content, dynamic route handlers as well as SSL certificates as needed.
Nows a good time to 

You should seem startup information simular to example 1. This time though the index.html file
delivered to your browser was be the in side this *static* sub-directory.

Now you are ready for steps three and four.

```shell
    . etc/config.sh
    ws
```

Like in your first example you will see some start up configuration. Notice that the Otto Engine
is turned on.  SSL should still be turned off. We can turn it on with a command line switch after
our first test.

Point your browser at the site [homepage](http://localhost:8000). Then point your browser at [/test](http://localhost:8000)

and _Otto Engine_

Otto engine is an experimental JavaScript engine implemented using [Otto](https://github.com/robertkrimen/otto) JavaScript virtual machine by Robert Krimen. The intended use to to allow the definition of route handlers in JavaScript and thus allow extensibility to _wsjs_ written in Golang.

Currently the Otto Engine is quite simple. All scripts found in the *WS_OTTO_PATH* are read in by _wsjs_ at start up. Once read they are each compiled in their own JavaScript virtual machine and their names used to define a URL path that they handle. If you define the directory "dynamic" as your *WS_OTTO_PATH* then all JavaScript files (files ending in ".js") will be evaluated with each become a URL end point.  If you have a file called "helloworld.js" inside the directory "dynamic" the Otto engine would use it to define what happens when the path "/helloworld" is request from _ws_.

Each Otto Engine defined route recieves a "Request" and "Response" object they can use to process the route (e.g. html forms submittions). The *Request* object provides full access it the request headers including two functions "GET()" and "POST()" for retrieving the values of form submissions. The *Response* object setups what gets sent to the browser by the web server. It has two methods if interest - "setHeader()" and "setContent()". Between these two objects and their methods a high level of useful API prototyping and be explored.

Some "features" in a dynamic webserver are not implemented from in _wsjs_ with Otto Engine. These
include--

+ Data Persistenence (e.g. access a database like MySQL on MongoDB)
+ Remote system access (e.g. CURL support)
+ File system access (reading or writing)

Forms of these may be introducted in the future.
