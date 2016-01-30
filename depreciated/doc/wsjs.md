
# wsjs

_wsjs_ is a simple web server with added support for dynamic content generated
from JavaScript files. This JavaScript environment is largely restricted
(see [otto](https://github.com/robertkrimen/otto)'s ducmentation) with notable
differences around setInterval() and setTimeout() which are missing and
nuances in [regular expression](https://golang.org/pkg/regexp/) handling.

If you need a fully featured server side JavaScript environment then I would
suggest [NodeJS](https://github.com/nodejs/node) but if all you need is some
base data manipulation or perhaps need to make calls to other web APIs then
_wsjs_ might do.

## built in libraries

_wsjs_ provides three functions in addition to the JavaScript engine implemented
with [otto](https://github.com/robertkrimen/otto). They are

+ HttpGet(url, headers) performs a http GET request, it is blocking
  + url is a string and should be a fully formed URL along with any parameters appropriately url encoded
    + url = "http://example.org/search?q=Hello%20World"
  + headers is an array of header objects where they property is the header type and value is the value to be passed
    + headers = [{"Content-Type": "text/plain"}]
  + returns the content requested or is empty
    + contents = HttpGet(url, headers)
+ HttpPost(url, headers, payload) performs a http POST request, it is blocking
  + url is a string and should be a fully formed URL along with any parameters appropriately url encoded
    + url = "http://example.org/search?q=Hello%20World"
  + headers is an array of header objects where they property is the header type and value is the value to be passed
    + headers = [{"Content-Type": "application/json"}]
  + payload is a properly encoded (based on your headers) POST response. Normally this is url encoded a JSON encoded
    + payload = JSON.stringify(data)
  + returns the content requested or is empty
    + contents = HttpPost(url, headers, payload)
+ Getenv(varname)
  + varname is the name you wish to get a value for, it is an empty string if varname is not found.
    + apiURL = Getenv("API_URL")
    
