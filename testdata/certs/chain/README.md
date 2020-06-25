# testdata/certs/chain

This directory contains concatenated certificates, which were created by
running the following commands in the `testdata/` directory.

## ecdsap256.pem

```shell
cat certs/issued/ecdsap256_sxg_60days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/ecdsap256.pem
```

## ecdsap384.pem

```shell
cat certs/issued/ecdsap384_sxg_60days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/ecdsap384.pem
```

## ecdsap521.pem

```shell
cat certs/issued/ecdsap521_sxg_60days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/ecdsap521.pem
```

## rsa4096.pem

```shell
cat certs/issued/rsa4096_sxg_60days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/rsa4096.pem
```

## non_sxg_cert.pem

```shell
cat certs/issued/ecdsap256_tls_60days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/non_sxg_cert.pem
```

## lasting_90days.pem

```shell
cat certs/issued/ecdsap256_sxg_90days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/lasting_90days.pem
```

## lasting_91days.pem

```shell
cat certs/issued/ecdsap256_sxg_91days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/lasting_91days.pem
```

## lasting_365days.pem

```shell
cat certs/issued/ecdsap256_sxg_365days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/lasting_365days.pem
```

## lasting_-1days.pem

```shell
cat certs/issued/ecdsap256_sxg_-1days.crt CA/inter/cert.pem CA/root/cert.pem \
    > certs/chain/lasting_-1days.pem
```

## without_root.pem

```shell
cat certs/issued/ecdsap256_sxg_60days.crt CA/inter/cert.pem \
    > certs/chain/without_root.pem
```

## self_signed.pem

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-03-31 00:00:00' \
    openssl req -new -x509 -config openssl.cnf \
    -extensions sxg_cert_selfsign -days 60 -subj '/CN=webpackager.test' \
    -key keys/ecdsap256.key -out certs/chain/self_signed.pem
```

## certmanager_0401.pem

```shell
cat certs/issued/certmanager_0401_15days.crt CA/inter/cert.pem \
    CA/root/cert.pem > certs/chain/certmanager_0401.pem
```

## certmanager_0415.pem

```shell
cat certs/issued/certmanager_0415_15days.crt CA/inter/cert.pem \
    CA/root/cert.pem > certs/chain/certmanager_0415.pem
```

## fake_acme_cert.pem

This is just an arbitrary certificate generated via:
https://docs.digicert.com/manage-certificates/certificate-profile-options/get-your-signed-http-exchange-certificate/
