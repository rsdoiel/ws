//
// This is an example Hello World serve side JavaScript demo.
//
(function (req, res) {
    "use strict";
    res.setHeader("content-type", "text/plain");
    res.setContent("Hello world, I'm alive!");
}(Request, Response))

