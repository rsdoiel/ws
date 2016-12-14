
# Self Signed SSL Certificates

You can easily make you own self signed certificates for development purposes using the _openssl_ command found on most Linux installations.

```shell
    mkdir -p etc/ssl
    cd etc/ssl
	openssl genrsa -out key.pem 2048
	openssl req -new -key key.pem -out cert.csr
	chmod 700 key.pem
```

Make sure your key.pem file is protected. For _ws_ command line options or environment variables to specify the key and cert pair.

```shell
   export WS_DOCROOT=Site
   export WS_URL=https://mysite.example.org:443
   export WS_SSL_KEY=etc/ws/key.pem
   export WS_SSL_CERT=etc/ws/cert.pem
   ws 
```

