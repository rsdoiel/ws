
# TODO fixes

+ Finish integrating utilities pulled in from stngo

# Todo ideas


## Someday Maybe

+ add "debug" option to ws, wsjs which logs full get/post/delete/put requests
+ add support to wsjs to handle JSON post as well as url encoded posts
+ add session and auth support
+ wsjs improvements 
    + add a CURL like object to the wsjs
    + add a built in DB object like SQLite or a JSON store for persistence between restarts
+ create wsedit for remote editing content of specific files over https connections.
    + look at CodeMirror and AceEditor as candidates for web browser editing
    + research best approach to embedding the editor in the go compiled binary
    + review scripted-editor for general apprach to problem
    + decide how to handle TLS key generation seemlessly
        + use existing certs for host machine
        + create one time self signed certs with signatures in browser display along with onetime URL
+ create a nav generator based on file system organization
    + autogenerate sitemap.xml and siteindex.html for improved search engine ingest
+ develop a generator and JS for browser side site specific search
    + create an inverted word list as JSON file
    + create a sitemap JSON file
+ explore cli tools as CMS to produce static websites
    + a markdown processor for generating static site content
    + shorthand for embedding content or as a alternate to Golang template pages
    + slugify/unslugify to generate appropriate page titles based on titles in text files.
+ explore interfacing with Solr
+ explore adding Lisp support
    + look at LispEx (a Schema?)
    + look at glisp (a lisp/schema dialect?)
    + Find a CL, FranzLisp (e.g. PC-LISP) or XLisp port to Go.
+ explore a wsphp based on PHP parse in https://github.com/stephens2424/php
+ shorthand evolutions
    + " :{ " should function as a closure
    + " {: " a verb who's object is read in a the global space.
    + add write of current Assignments with :}
+ Try to limit the glyphs (verbs) to 7
    + assign value to symbol
    + assign value of an included file to a symbol
    + assign value of a eval'd expression to a symbol
    + assign value of a eval'd file to a symbol (with closure)
    + assign value of a eval'd file to a symbol (without a closure)
    + assign value of from a shell expression to a symbol
    + write out the current assignments to a file
+ create a repl out of shorthand by adding a prompt option and exit glyph.

