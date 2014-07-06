/**
 * helloworld.js - This is a simple Hello World program for
 * implementing a route end poit using the OttoEngine
 *
 * @param global - the global object contains two default objects
 * Request and Response. Request exposes the http/https request object that Golang
 * provides, Response allows manipulation of the resulting Golang Response
 * @return the route handle should return a string or nil. This will become of body
 * of the http/https response Golang uses.
 */
(function (global) {
    console.log("Hello from inside of the OttoEngine.");
    return 'OttoEngine says Hello World!';
}(this));

