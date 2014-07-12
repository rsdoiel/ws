/* This is a route handler for /hi */
(function (Request, Response){
    console.log("This is the hi route");
    Response.setHeader('Content-Type', 'text/html');
    return [
    "<!DOCTYPE html>",
    "<html>",
    "<body>",
    '<em>Hi There.</em> Go <a href="/">home</a>.',
    "</body>",
    "</html>"].join("\n");
}(Request, Response));
