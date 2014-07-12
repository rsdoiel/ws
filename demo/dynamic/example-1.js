/* example-1.js - a simple example of Request and Response objects */
(function (req, res) {
    var header = req.Header;

    res.setHeader("content-type", "text/html");
    res.setContent("<p>Here is the Header array received by this request</p>" +
        "<pre>" + JSON.stringify(header) + "</pre>");
}(Request, Response));
