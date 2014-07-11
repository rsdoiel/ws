/**
 * helloworld-3.js - This a more complicate example. We are going to process
 * a little form so we can personalize the greeting.
 */
(function (req, res) {
    var raw_params = [],
        GET = {},
        output = "";

    // 1. All our transactions will be of text/html type so go ahead and set that now.
	res.setHeader('Content-Type', 'text/html');
    // 2. First see if we have a query string, and proces it
    if (req.Method === "GET") {
        GET = req.GET();
    }
    // 3. If GET.name then display greeting
    if (GET.name) {
        output = "<h1>Hello " + GET.name + "</h1>";
    }    
    // 4. Otherwise display webform
    output += '<p><form method="get" action="/helloworld-3">Excuse me, what is your name? <input type="text" name="name"> <button type="submit">Answer</button></form>';

    res.setContent(output);
}(Request, Response));
