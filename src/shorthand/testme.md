
Q: What is _ws_?

A: A nimble webserver and friends.

Q: What would JavaScipt route look like for _wsjs_?

A: Something like....

```JavaScript
    // Demo the availability of Request and Response objects.
    (function (req, res) {
        var content = ["<!DOCTYPE html>",
            "<html>",
            "<body>",
            "<h1>Test</h1>",
            "<h2>Request object as JSON</h2>",
            "<pre><code>",
            JSON.stringify(req, null, true),
            "</code></pre>",
            "<h2>Response object as JSON </h2>",
            "<p>(before running setContent())</p>",
            "<pre><code>",
            JSON.stringify(res, null, true),
            "</code></pre>",
            "</body>",
            "</html>"].join("\n");

        res.setHeader('Content-Type', 'text/html');
        res.setContent(content);
    }(Request, Response));
```
