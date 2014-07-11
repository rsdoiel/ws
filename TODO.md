
# Todo items

## v0.0.1 pre-release

+ Finish ottoengine prototype
    - Figure out how to support extended path E.g.
        + if /test is handled by test.js, it should also handle /test/something/else
+ Make sure docs cover all 
    - command line options
    - environment settings
    - static content 
    - dynamic routes defined in JavaScript via Otto Engine
    - have some tutorial examples of building a small interactive fiction game.
+ When ready tag as v0.0.1


# Someday, maybe

+ Some sort of ephmeral shared data (e.g. a map of maps, where the outer map is a collections of collection, and the inner maps are the individual collections shared between routes)
+ Built in page speed support (e.g. automatic cache control headers, expire and etag; gzip content)
+ Support a shared in-memory map between route requests (memcache light) to reduce remote content calls
+ Figure out how to pass in a DB connector object for persistence
+ Figure out how to function as middleware
    - integrate a server side xhr object for bridging to remote services (e.g. content storage systems like Dropbox)
+ Implement a file watcher to reload scripts on change rather then only at server startup.
+ Add deb package support
+ Add Mac ports package support
+ Add support for Lisp

