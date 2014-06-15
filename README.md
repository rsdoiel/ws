ws
==

A light weight TLS webserver suitable for static web development

# USAGE

Make sure _ws_ is your path. Make sure your know where your *cert.pem*
and *key.pem* files are located. By default _ws_ looks in *$HOME/etc/ws*.
Pick an open port you have access to (by default the port is 8443).
Invoke the command _ws_ in a shell.

Here is an example of specifying all the options on the command line.
It assume your TLS cert/key are called */etc/my-cert.pem* and */etc/my-key.pem* 
respectively, your docroot will be */www* and you'll run TLS server
on port 8007.

```shell
    ws -cert=/etc/my-cert.pem -key=/etc/my-key.pem -docroot=/www -port=8007
```

Now point your web browser at [https://localhost:8007](https://localhost:8007). This is a TLS based webserver. You need to use the **https** protocol in your URL.


