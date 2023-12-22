// https://blog.csdn.net/dl962454/article/details/124330966

CA证书生成:

``` 
# 生成CA私钥(ca.key)
openssl genrsa -des3 -out ca.key 2048 
# 生成CA证书签名请求(ca.csr)
openssl req -new -key ca.key -out ca.csr
# 生成自签名CA证书(ca.cert)
openssl x509 -req -days 3650 -in ca.csr -signkey ca.key -out ca.crt
```

``` 
// 服务端私钥
openssl genrsa -out server.key 2048
//服务端签名请求
openssl req -new -sha256 -out server.csr -key server.key -config server.conf
//用根证书签发服务端证书server.pem
openssl x509 -req -days 3650 -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
```