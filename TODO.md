
# Todo items

## Finish v0.0.1 Alpha

+ Add expire, cache control and eTag headers to all requests (e.g. fsengine and ottoengine)
+ Confirm Server can support subsecond second responses for dynamic routes
+ Add a Persistence layer
    - Look at Tiedot, figure out if it can refactored to run on a Raspberry Pi
+ Finish ottoengine
    - Make it more like http.FileServe()
    - Add support for route files as path prefix
        + if /test is handled by test.js, it could also handle /test/something/else
+ Make sure docs cover all 
    - command line options
    - environment settings
    - static content settings
    - dynamic routes and how they work
    - a demo of a small interactive game show casing _ws_
+ Provide some test coverage
    - make sure -keygen is generating valid keys
    - make sure -init works in generating both default and customized setup
    - make sure turning on/off ottoengine and tls works
+ When ready tag as v0.0.1-beta


# Someday, maybe

+ Make a better -init and -keygen CUI using come sort of Curses like control
+ Add a page speed like module (e.g. automatic cache control headers, expire and etag; gzip content)
+ Explore a server side xhr object to let ottoengine function as middleware to other web services
+ Explore the Web Browser File API and see if that can be used as a basis for safe file system interaction
+ Integration OS Signals support (e.g. restart, reload)
+ Implement a file watcher to refresh in-memory content (e.g. JavaScript assets)
+ Implement some sort of web component server side pre-rendering
+ Add support for alternative JS engines (e.g. SpiderMonkey)
+ Add deb package support
+ Add Mac ports package support
+ Add support for Common Lisp

