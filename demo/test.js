// Demo the availability of Request and Response objects.
(function () {
    return ["<!DOCTYPE html>",
        "<html>",
        "<body>",
        "<h1>Test</h1>",
        "<h2>Request object as JSON</h2>",
        "<pre>",
        JSON.stringify(Request, null, true),
        "</pre>",
        "</body>",
        "</html>"].join("\n");
}());
