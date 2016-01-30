/**
 * helloworld-4.js - This time we're expanding on the previous example and using a 
 * "POST" method in our HTML form. As a result we use the req.POST() function to
 * get back the data.
 */
(function (req, res) {
    var raw_params = [],
        POST = {},
        output = "";

    // 1. All our transactions will be of text/html type so go ahead and set that now.
	res.setHeader('Content-Type', 'text/html');

    // 2. First see if we have a query string, and proces it
    if (req.Method === "POST") {
        POST = req.POST();
    }
    // 3. If GET.name then display greeting
    if (POST.name) {
        output = "<h1>Hello " + POST.name + "</h1><p>A few moments latter &hellip;</p>";
    }    
    // 4. Otherwise display webform
    output += '<p><form method="POST" action="/helloworld-4">Forgetful Harry asks, &ldquo;Excuse me, what is your name?&rdquo; <input type="text" name="name"> <button type="submit">Answer</button></form>';

    res.setContent(output);
}(Request, Response));
