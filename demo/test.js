// Demo the availability of Request and Response objects.
(function (req, res) {
    res.setHeader('Content-Type', 'text/html');
    res.setContent( ["<!DOCTYPE html>",
        "<html>",
        "<body>",
        "<h1>Test</h1>",
        "<h2>Request object as JSON</h2>",
        "<pre>",
        JSON.stringify(req, null, true),
        "</pre>",
        "</body>",
        "</html>"].join("\n"));
        return res;
}(Request, Response));
