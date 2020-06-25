# testdata/keys

This directory contains private keys for testing. They were created
by the following commands:

## ecdsap256.key

```shell
openssl ecparam -name prime256v1 -genkey -out ecdsap256.key
```

## ecdsap384.key

```shell
openssl ecparam -name secp384r1 -genkey -out ecdsap384.key
```

## ecdsap521.key

```shell
openssl ecparam -name secp521r1 -genkey -out ecdsap521.key
```

## rsa4096.key

```shell
openssl genrsa -out rsa4096.key 4096
```

## pkcs8-ecdsa.key

```shell
openssl pkcs8 -nocrypt -topk8 -in ecdsap256.key -out pkcs8-ecdsa.key
```

## pkcs8-rsa.key

```shell
openssl pkcs8 -nocrypt -topk8 -in rsa4096.key -out pkcs8-rsa.key
```

## pkcs8-multi.key

```shell
cat pkcs8-rsa.key pkcs8-ecdsa.key > pkcs8-multi.key
```
