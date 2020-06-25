# testdata/certs/cbor

## ecdsap256_nosct.cbor

```shell
gen-certurl -pem certs/chain/ecdsap256.pem -ocsp ocsp/ecdsap256_7days.ocsp \
    > certs/cbor/ecdsap256_nosct.cbor
```

## ecdsap384_nosct.cbor

```shell
gen-certurl -pem certs/chain/ecdsap384.pem -ocsp ocsp/ecdsap384_7days.ocsp \
    > certs/cbor/ecdsap384_nosct.cbor
```

## self_signed.cbor

```shell
gen-certurl -pem certs/chain/self_signed.pem -ocsp <(echo -n 'dummy-ocsp') \
    > certs/cbor/self_signed.cbor
```

## certmanager_0401_0409.cbor

```shell
gen-certurl -pem certs/chain/certmanager_0401.pem \
    -ocsp ocsp/certmanager_0401_0409.ocsp \
    > certs/cbor/certmanager_0401_0409.cbor
```

## certmanager_0401_0413.cbor

```shell
gen-certurl -pem certs/chain/certmanager_0401.pem \
    -ocsp ocsp/certmanager_0401_0413.ocsp \
    > certs/cbor/certmanager_0401_0413.cbor
```

## certmanager_0415.cbor

```shell
gen-certurl -pem certs/chain/certmanager_0415.pem \
    -ocsp ocsp/certmanager_0415_0415.ocsp \
    > certs/cbor/certmanager_0415_0415.cbor
```
