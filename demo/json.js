/**
 * json-demo.js - return a blob of JSON content.
 */
/*jslint browser: false, indent: 4 */
(function (Request, Response) {
    var blob = {
        Message: "Hello World",
        i: 1,
        today: new Date()
    };
    Response.setHeader('Content-Type', 'application/json');
    return JSON.stringify(blob);
}(Request, Response));
