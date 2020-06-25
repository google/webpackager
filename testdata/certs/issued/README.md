# testdata/certs/issued

This directory contains certificates issued by Test Intermediate CA, which
were created by running the following commands in the `testdata` directory.

## ecdsap256_sxg_60days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/ecdsap256_sxg_60days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap256_sxg_60days.csr \
    -out certs/issued/ecdsap256_sxg_60days.crt \
    -startdate 20200401000000Z -enddate 20200531000000Z -notext -batch
```

## ecdsap384_sxg_60days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap384.key \
    -out certs/issued/ecdsap384_sxg_60days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap384_sxg_60days.csr \
    -out certs/issued/ecdsap384_sxg_60days.crt \
    -startdate 20200401000000Z -enddate 20200531000000Z -notext -batch
```

## ecdsap521_sxg_60days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap521.key \
    -out certs/issued/ecdsap521_sxg_60days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap521_sxg_60days.csr \
    -out certs/issued/ecdsap521_sxg_60days.crt \
    -startdate 20200401000000Z -enddate 20200531000000Z -notext -batch
```

## rsa4096_sxg_60days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/rsa4096.key \
    -out certs/issued/rsa4096_sxg_60days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/rsa4096_sxg_60days.csr \
    -out certs/issued/rsa4096_sxg_60days.crt \
    -startdate 20200401000000Z -enddate 20200531000000Z -notext -batch
```

## ecdsap256_sxg_90days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/ecdsap256_sxg_90days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap256_sxg_90days.csr \
    -out certs/issued/ecdsap256_sxg_90days.crt \
    -startdate 20200401000000Z -enddate 20200630000000Z -notext -batch
```

## ecdsap256_sxg_91days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/ecdsap256_sxg_91days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap256_sxg_91days.csr \
    -out certs/issued/ecdsap256_sxg_91days.crt \
    -startdate 20200401000000Z -enddate 20200701000000Z -notext -batch
```

## ecdsap256_sxg_365days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/ecdsap256_sxg_365days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap256_sxg_365days.csr \
    -out certs/issued/ecdsap256_sxg_365days.crt \
    -startdate 20200401000000Z -enddate 20210401000000Z -notext -batch
```

## ecdsap256_sxg_-1days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/ecdsap256_sxg_-1days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap256_sxg_-1days.csr \
    -out certs/issued/ecdsap256_sxg_-1days.crt \
    -startdate 20200401000000Z -enddate 20200331000000Z -notext -batch
```

## ecdsap256_tls_60days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/ecdsap256_tls_60days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions tls_cert \
    -in certs/issued/ecdsap256_tls_60days.csr \
    -out certs/issued/ecdsap256_tls_60days.crt \
    -startdate 20200401000000Z -enddate 20200531000000Z -notext -batch
```

## ecdsap256_sxg_revoked.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/ecdsap256_sxg_revoked.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/ecdsap256_sxg_revoked.csr \
    -out certs/issued/ecdsap256_sxg_revoked.crt \
    -startdate 20200401000000Z -enddate 20200531000000Z -notext -batch

NO_FAKE_STAT=1 TZ=UTC faketime '2020-04-01 12:00:00' \
    openssl ca -config openssl.cnf -name CA_inter \
    -revoke certs/issued/ecdsap256_sxg_revoked.crt -crl_reason keyCompromise
```

## certmanager_0401_15days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/certmanager_0401_15days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/certmanager_0401_15days.csr \
    -out certs/issued/certmanager_0401_15days.crt \
    -startdate 20200401000000Z -enddate 20200416000000Z -notext -batch
```

## certmanager_0415_15days.crt

```shell
openssl req -new -subj '/CN=webpackager.test' -key keys/ecdsap256.key \
    -out certs/issued/certmanager_0415_15days.csr

openssl ca -config openssl.cnf -name CA_inter -extensions sxg_cert \
    -in certs/issued/certmanager_0415_15days.csr \
    -out certs/issued/certmanager_0415_15days.crt \
    -startdate 20200415000000Z -enddate 20200430000000Z -notext -batch
```
