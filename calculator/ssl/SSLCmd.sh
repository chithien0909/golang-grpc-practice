#openssl genrsa -out ssl/server.key 2048
#
#openssl req -nodes -new -x509 -sha256 -days 1825 -config ssl/cert.conf -extensions 'req_ext' -key ssl/server.key -out ssl/server.crt

#SERVER_CN=localhost
openssl genrsa -passout pass:123456 -des3 -out ca.key 4096
openssl req -passin pass:123456 -new -x509 -days 3650 -key ca.key -out ca.crt

#openssl genrsa -passout pass:123456 -des3 -out server.key 4096
#openssl req -passin pass:123456 -new -key server.key -out server.csr
#openssl x509 -req -passin pass:123456 -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
#openssl pkcs8 -topk8 -nocrypt -passin pass:123456 -in server.key -out server.pem