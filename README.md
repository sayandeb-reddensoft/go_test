# cdc-auth

# How to use :

1.  Clone this repo
2.  Run `go mod tidy` to install all the required dependencies
3.  Rename `.env.sample` to `.env` & replace values
4.  **Generate JWT keypairs :**<br/>
    - `mkdir keys`<br/>
    - `cd ./keys` & `mkdir temp`<br/>
    - Run `openssl genrsa -des3 -out ./temp/private.pem 2048` & enter the passphrase for it<br/>
    - Run `openssl rsa -in ./temp/private.pem -outform PEM -pubout -out public.pem` & enter the same passphrase for it<br/>
    - Run `openssl pkcs8 -in ./temp/private.pem -topk8 -nocrypt -out ./private.pem`& enter the same passphrase for it<br/>
    - Run `rm -r ./temp`
5.  `cd ..` & Run `go run main.go`
