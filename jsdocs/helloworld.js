//
// This is an example Hello World serve side JavaScript demo.
//
(function (req, res) {
    console.log("Incoming request", JSON.stringify(req));
    res.code = 200;
    res.headers = [{"Content-Type": "text/html"}];
    res.content = "<DOCTYPE html>\n<html><head><title>Testing Server Side JavaScript</title></head><body><h1>Testing Server Side JavaScript</h1><p>Hello World!</p><a href=\"/\">homepage</a></p></body></html>";
    console.log("Processed response", JSON.stringify(res));
}(Request, Response));
