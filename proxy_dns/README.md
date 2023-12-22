新建ca.conf,填入以下内容：
``` 
[ req ]
default_bits       = 4096
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
countryName                 = GB
countryName_default         = BeiJing
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = JiangSu
localityName                = Locality Name (eg, city)
localityName_default        = NanJing
organizationName            = Organization Name (eg, company)
organizationName_default    = Step
commonName                  = liuqh.icu
commonName_max              = 64
commonName_default          = liuqh.icu
```


生成ca.key：
``` 
 openssl genrsa -out ca.key 4096
```

生成CA证书：

``` 
$ openssl req -new -x509 -days 365 -subj "/C=GB/L=Beijing/O=github/CN=liuqh.icu" \
-key ca.key -out ca.crt -config ca.conf


C=GB: C代表的是国家名称代码。
L=Beijing: 代表地方名称,例如城市。
O=gobook: 代表组织单位名称。
CN=liuqh.icu: 代表关联的域名，

```

2.server证书生成：
在hello/keys目录下：
新建server.conf,填入以下内容：

``` 
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = JiangSu
localityName                = Locality Name (eg, city)
localityName_default        = NanJing
organizationName            = Organization Name (eg, company)
organizationName_default    = Step
commonName                  = CommonName (e.g. server FQDN or YOUR name)
commonName_max              = 64
commonName_default          = liuqh
[ req_ext ]
subjectAltName = @alt_names
[alt_names]
DNS.1   = liuqh.icu
IP      = 127.0.0.1
```

生成公私钥

``` 
$ openssl genrsa -out server.key 2048
```

. 生成CSR

``` 
$ openssl req -new  -subj "/C=GB/L=Beijing/O=github/CN=liuqh.icu" \
-key server.key -out server.csr -config server.conf
```

. 基于CA签发证书

``` 
$ openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 365 \
-in server.csr -out server.crt -extensions req_ext -extfile server.conf
```

生成客户端证书

1. 生成公私钥
``` 
$ openssl genrsa -out client.key 2048 
```

2. 生成CSR
``` 
   $ openssl req -new -subj "/C=GB/L=Beijing/O=github/CN=liuqh.icu"  \
   -key client.key -out client.csr 
```

3.基于CA签发证书

``` 
   $ openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -days 365 \
   -in client.csr -out client.crt
```   