
# Ideas to explore


## Next

+ doc path should be a non-flag reference the command line `ws mySite`
+ -htdocs should be -d, -docs and avoid confusion -h, -help in command line parsing
+ Bring command line option into align with other web tools
    + -l license
    + -v version
    + -h help
+ Drop support for embedded JS engine 
    + Too hard to explain in relation to NodeJS
    + JS is moving beyond ES5 to ES6 in the browser, otto is ES5
    + [mkpage](https://rsdoiel.github.io/mkpage) can already include external JSON resources
    + I never used it beyond demos
+ Add option -n, -no-cache header for web development mode

## Someday, maybe

+ Add support for ACME SSL/TLS certs (SSL Everywhere)
+ Add support for on the fly compression (gzip) of text/* content types
+ Intergrate [mkpage](https://rsdoiel.github.io/mkpage) single page mode 
+ Look at [Echo](http://echo.labstack.com/) router and see if this is a useful way to support FastCGI
+ Look at [libsecurity](https://developer.ibm.com/open/libsecurity/) and see how it might help *ws* stay safe.


