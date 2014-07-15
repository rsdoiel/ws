
#  Getting started with _ws_ and _Otto Engine_

Otto engine is an experimental JavaScript engine implemented using [Otto](https://github.com/robertkrimen/otto) JavaScript virtual machine by Robert Krimen. The intended use to to allow the definition of route handlers in JavaScript and thus allow extensibility to _ws_ written in Golang.

Currently the Otto Engine is quite simple. All scripts found in the *WS_OTTO_PATH* are read in by _ws_ at start up. Once read they are each compiled in their own JavaScript virtual machine and their names used to define a URL path that they handle. If you define the directory "dynamic" as your *WS_OTTO_PATH* then all JavaScript files (files ending in ".js") will be evaluated with each become a URL end point.  If you have a file called "helloworld.js" inside the directory "dynamic" the Otto engine would use it to define what happens when the path "/helloworld" is request from _ws_.

Each Otto Engine defined route recieves a "Request" and "Response" object they can use to process the route (e.g. html forms submittions). The *Request* object provides full access it the request headers including two functions "GET()" and "POST()" for retrieving the values of form submissions. The *Reponse* object setups what gets sent to the browser by the web server. It has two methods if interest - "setHeader()" and "setContent()". Between these two objects and their methods a high level of useful API prototyping and be explored.

The following "features" are currently remain unimplemented--

+ File system access (reading or writing)
+ Remote system access (e.g. CURL support)
+ Data Persistenence (e.g. access a database like MySQL on MongoDB)

These features may be introduced in the future but at present are not
required in the _ws_ role of prototyping tool.


## Tutorials

### Setting up

Here are the steps I use to build a typical
prototype. My prototype will be build in a directory called "demo"

1. Create a directory from my prototype.
    - *mkdir demo; cd demo*
2. Create a directory for my static content
    - *mkdir static*
3. In the "demo" directory create a Bash shell script named *ws-setup.sh* with the following exported environment variables--

```shell
    #!/bin/bash
    # I usually call this file

    # Example prototype configuration
    export WS_HOST=localhost
    export WS_PORT=8000
    export WS_DOCROOT=$(pwd)/static
```

These can then be made availab to _ws_ by "sourcing" the Bash script.

```shell
    . ws-setup.sh
```

You can confirm that variables *WS_OTTO*, *WS_OTTO_PATH* and *WS_DOCROOT* are set appropriately by "echoing" them to the console--

```shell
    echo "Docroot $WS_DOCROOT, Otto $WS_OTTO and $WS_OTTO_PATH"
```

If you have made it this far you are ready to begin prototyping but before we start _ws_ we need to at least have some static content. Create *index.html* in the static directory.

```HTML
    <!DOCTYPE html>
    <html>
        <head><title>Hello from ws</title></head>
        <body>Hello from ws</body>
    </html>
```

Now you are ready to start _ws_ webserver.  In your shell run the *ws* command--

```shell
    ws
```

As _ws_ starts you should see the configuration logged to the script. When _ws_ is ready it will display a startup log message indicating the URL it is listening for (e.g. http://localhost:8000). Point your web browser at it and you should see the HTML page we created.

To stop _ws_ switch back to your console window and press "ctr-c" (holding down the control key you press the letter c). _ws_ should quick.

### Hello World

In the first tutorial we setup _ws_ to support development of static content.  This is useful. You can test your CSS, HTML and browser side JavaScript with this basic configuration.  At some point you may need to prototype a more dynamic (server side) solution. That is where Otto Engine comes in. Otto Engine supports a restricted JavaScript virtual machine where you can simulte things like API responses (e.g. Send back JSON information based on a form submission).  To add this to our protype we will do three things

1. create a subdirectory for our dynamic content
2. Update our *ws-setup.sh* to include the Otto Engine settings
3. create a JavaScript file that returns a JSON object "hello world" message




