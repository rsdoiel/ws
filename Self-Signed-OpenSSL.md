
# Self Signed SSL Certificates

You can easily make you own self signed certificates for development purposes using the _openssl_ command found on most Linux installations.

```shell
	sudo openssl genrsa -out key.pem 2048
	sudo openssl req -new -key key.pem -out cert.csr
	sudo chmod 700 key.pem
```

Make sure your key.pem file is protected. I usually put these in $HOME/etc/ws/ for my _ws_ development server.

