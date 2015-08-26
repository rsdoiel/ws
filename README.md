ws
==

    A nimble webserver with friends

_ws_ started as a nimble static content webserver. It grew some friends along the way.  Now it is a small collection for command line utilities.

# _ws_ has friends

_ws_ is simple static content webserver with optional support for SSL. It is careful to avoid dot files and by default with server the content of the directory where it is launch. _ws_ has a companion webserver named _wsjs_. It is adds the ability to dynamic create content via JavaScript files.  Each JavaScript file becomes a route itself. _wsjs_ is handy when you're mocking up services or RESTful API.

_ws_ and _wsjs_ have friends. These are additional command line utilities useful for quickly prototyping websites. They include _wsinit_ which creates a project directory skeletons and configuration interactively, _wskeygen_ which generates self signed SSL certificates. Additional commands that are helpful in preprocessing text or augmenting shell scripts.

The complete list of utilities

+ [ws](doc/ws.md) a nimple webserver for static content
    + Built on Golangs native http/https modules
    + Restricted file service, only from the docroot and no "dot files" are served
    + No dynamic content support 
    + Quick startup, everything logged to console for easy debugging or piping to a log processor
+ [wsjs](doc/wsjs.md) a nimple webserver for static and JavaScript generated dynamic content
    + built on Robert Krimen's excellent [otto](https://github.com/robertkrimen/otto) JavaScript VM
    + Mockup your dynamic content via JavaScript defined routes (great for creating JSON blobs used by a browser side demo)
+ [wsinit](document) will generate a project structure as well as a shell script for configuring _ws_ and _wsjs_ in a 12 factor app manor.
+ [wsmarkdown](doc/wsmarkdown.md) renders markdown as HTML
    + this markdown is built on the [Blackfriday](https://github.com/russross/blackfriday) Markdown library
+ [shorthand](doc/shorthand.md) a domain specific language that expands shorthand labels into associated replacement values
+ [slugify](doc/slugify.md) turns text phrases into readibly URL friendly phrases (e.g. "Hello World" becomes "Hello_World")
+ [unslugify](doc/unslugify.md) returns slugs into simple phrase (e.g. "Hello_World" becomes "Hello World")

Got an idea for a new project? Want to prototype it quickly? 

1. run "wsinit" to set things up
2. run ". etc/config.sh" seed your environment
3. run "ws" and start working!
4. run "wsjs" to mockup a RESTful service

_ws_ and _wsjs_ feature sets have been kept minimal. Only what you need when you turn it on.

Want to expand out the site quickly, write the HTML skeleton with
markdown, sprinkle in some shorthand which can leverage some shell logic
and now you have HTML pages with common nav, headers, and footers.

## LICENSE

copyright (c) 2014 - 2015 All rights reserved.
Released under the [Simplified BSD License](http://opensource.org/licenses/bsd-license.php)

