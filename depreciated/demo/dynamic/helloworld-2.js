/**
 * helloworld-2.js - This a more typical Hello World program.  We take
 * advantage of the Request and Response objects passed into our closure.
 */
(function (req, res) {
	res.setHeader('Content-Type', 'text/html');
	res.setContent("<h1>Hello World!</h1><p>this is HTML</p>");
}(Request, Response));
