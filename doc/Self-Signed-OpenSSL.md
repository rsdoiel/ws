
# Self Signed SSL Certificates

You can easily make you own self signed certificates for development purposes using the _openssl_ command found on most Linux installations.

```shell
    mkdir -p etc/ssl
    cd etc/ssl
	openssl genrsa -out key.pem 2048
	openssl req -new -key key.pem -out cert.csr
	chmod 700 key.pem
```

Make sure your key.pem file is protected. For _ws_ and _wsjs_ use command line options or environment variables to specify the certs.

You can also use the _ws -init_ command, to generate a directory structure as well as SSL certs if the WS_URL environment variable starts with an https protocol.

```shell
   export WS_URL=https://mysite.example.org:443
   ws -init
```

Same is try of the _ws -init_ command.

