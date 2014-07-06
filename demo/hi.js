/* This is a route handler for /hi */
(function (global){
    console.log("This is the hi route");
    return [
    "<!DOCTYPE html>",
    "<html>",
    "<body>",
    '<em>Hi There.</em> Go <a href="/">home</a>.',
    "</body>",
    "</html>"].join("\n");
}(this))
