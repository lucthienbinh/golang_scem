# run command in terminal: sh createTLS.sh 

openssl genrsa -out key.pem 2048
openssl ecparam -genkey -name secp384r1 -out key.pem
openssl req -new -x509 -sha256 -key key.pem -out cert.bem -days 3650

# run command in terminal: go run `go env GOROOT`/src/crypto/tls/generate_cert.go --host=127.0.0.1:5000 --ecdsa-curve=P256