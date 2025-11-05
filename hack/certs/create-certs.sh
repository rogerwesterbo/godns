#!/bin/bash

# check if openssl is installed
if ! command -v openssl &> /dev/null
then
    echo "openssl could not be found, please install it"
    exit
fi

function recreateAll {
  # create certs directory if it doesn't exist
  openssl genrsa -out ca.key 4096
  openssl req -x509 -new -key ca.key -days 3650 -sha256 -subj "/CN=GoDNS-Local-CA" -out ca.crt

  openssl genrsa -out server.key 4096
  openssl req -new -key server.key -subj "/CN=GoDNS-Server" -out server.csr

  openssl genrsa -out client.key 4096
  openssl req -new -key client.key -subj "/CN=GoDNS-Client" -out client.csr

  openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out server.crt -days 825 -sha256

  openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out client.crt -days 825 -sha256

  # create pem file for server
  cat server.crt server.key > server.pem
  cat client.crt client.key > client.pem
  cat ca.crt client.key > ca.pem    
}

# check all certs if one is expired
for cert in *.crt; do
    if openssl x509 -checkend 86400 -noout -in "$cert"; then
        echo "$cert is valid for more than 1 day"
    else
        echo "$cert is expired or will expire within 1 day"
        rm -f ca.crt ca.key ca.srl server.crt server.key server.csr client.crt client.key client.csr server.pem client.pem ca.pem
        break
    fi
done

# if any of the certs is missing, recreate all
if [ ! -f ca.crt ] || [ ! -f ca.key ] || [ ! -f server.crt ] || [ ! -f server.key ] || [ ! -f client.crt ] || [ ! -f client.key ] || [ ! -f server.pem ] || [ ! -f client.pem ] || [ ! -f ca.pem ]; then
    echo "One or more certificates are missing, recreating all"
    recreateAll
else
    echo "All certificates are present and valid"
fi