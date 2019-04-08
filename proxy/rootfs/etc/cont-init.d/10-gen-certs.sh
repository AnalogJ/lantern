#!/usr/bin/with-contenv sh
set -ex

cd /srv/lantern/certs

# generate CA's  key
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:4096 -out ca.key.pem

openssl req -config /defaults/openssl.cnf -key ca.key.pem -new -x509 -days 7300 -sha256 -extensions v3_ca -out ca.crt

# generate mobileconfig file

lantern_mobileconfig