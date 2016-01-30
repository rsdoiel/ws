/**
 * html.js - Output a simple HTML page.
 */
/*jslint browser: false, indent: 4 */
(function (Request, Response){
    "use strict";
    Response.setHeader('Content-Type', 'text/html');
    return [
        "<!DOCTYPE html>",
        "<html>",
        "\t<head><title>HTML Demo</title></head>",
        "\t<body>",
        "\t\t<h1>HTML Demo</h1>",
        "\t\t<p>This is a simple HTML page.</p>",
        "\t\t<ul>",
        "\t\t\t<li>One</li>",
        "\t\t\t<li>Two</li>",
        "\t\t\t<li>Three</li>",
        "\t\t</ul>",
        "\t</body>",
        "</html>"
    ].join("\n");
}(Request, Response));
