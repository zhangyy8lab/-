- create OpenSSl setting file -> openssl.conf
```go
[ req ]
default_bits        = 2048
default_md          = sha256
prompt              = no
encrypt_key         = no
default_keyfile     = server.key

distinguished_name  = req_distinguished_name
req_extensions      = req_ext

[ req_distinguished_name ]
countryName            = US
stateOrProvinceName    = California
localityName           = San Francisco
organizationName       = Your Company
commonName             = grpc.demo.test 

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
```

- generate .key and .crt for openss.cof
```go
// 生成私钥
openssl genrsa -out server.key 2048

# 生成证书签名请求 (CSR) 使用自定义配置
openssl req -new -key server.key -out server.csr -config openssl.cnf

# 生成自签名证书
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt -extensions req_ext -extfile openssl.cnf

```

- use ssl see server/main.go && client/main.go
