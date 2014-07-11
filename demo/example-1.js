/* example-1.js - a simple example of Request and Response objects */
(function (req, res) {
    var headers = req.Headers;

    res.setHeader("content-type", "text/html");
    res.setContent("<p>Here are the headers received by this request</p>" +
        "<pre>" + JSON.stringify(headers) + "</pre>");
}(Request, Response));
