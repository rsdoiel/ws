//
// This is an example Hello World serve side JavaScript demo.
//
(function (req, res) {
    console.log("Incoming request", JSON.stringify(req));
    res.code = 200;
    res.headers = [{"Content-Type": "text/plain"}];
    res.content = "Hello World!";
    console.log("Processed response", JSON.stringify(res));
}(Request, Response));
