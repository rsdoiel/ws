// Demo the availability of Request and Response objects.
(function () {
    Response.ContentType = "text/html";
    return ["<h1>Hello World!</h1>",
        "<h2>Request object as JSON</h2>",
        "<pre>",
        JSON.stringify(Request),
        "</pre>",
        "<h2>Response object as JSON</h2>",
        "<pre>",
        JSON.stringify(Response),
        "</pre>"].join("");
}());
