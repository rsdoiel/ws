/**
 * json-demo.js - return a blob of JSON content.
 */
/*jslint browser: false, indent: 4 */
(function (Request) {
    var blob = {
        Message: "Hello World",
        i: 1,
        today: new Date()
    };
    return JSON.stringify(blob);
}(Request));
