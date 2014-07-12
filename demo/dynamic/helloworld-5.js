/**
 * helloworld-5.js - This time we're expanding on the previous example and 
 * using a "PUT" method in our HTML form. As a result we use the req.PUT() 
 * function to get back the data.
 */
(function (req, res) {
    var raw_params = [],
        PUT = {},
        output = "";

    // 1. All our transactions will be of text/html type so go ahead and set that now.
	res.setHeader('Content-Type', 'text/html');
    // 2. First see if we have a query string, and proces it
    if (req.Method === "PUT") {
        PUT = req.PUT();
    }
    // 3. If GET.name then display greeting
    if (PUT.name) {
        output = "<h1>Hello " + PUT.name + "</h1><p>I will try to remember your name next time.</p>";
    } else {
        output = "<p>I'm confused, did you say your name?</p>";
    }
    // 4. Otherwise display webform
    res.setContent(output);
}(Request, Response));
