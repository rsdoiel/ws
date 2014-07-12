/**
 * view-request.js - Display the contents of te Request object
 * as a JSON file.
 */
(function (req, res) {
    res.setHeader("Content-Type", "application/json");
    res.setContent(JSON.stringify(req));
}(Request, Response))
