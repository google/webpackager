# testdata/ocsp

This directory contains OCSP responses in DER form, which were created by
running the following commands in the `testdata/` directory.

## ecdsap256_7days.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-05-01 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
     -issuer CA/inter/cert.pem -cert certs/issued/ecdsap256_sxg_60days.crt \
    -ndays 7 -respout ocsp/ecdsap256_7days.ocsp -no_nonce -resp_no_certs
```

## ecdsap256_8days.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-05-01 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
     -issuer CA/inter/cert.pem -cert certs/issued/ecdsap256_sxg_60days.crt \
    -ndays 8 -respout ocsp/ecdsap256_8days.ocsp -no_nonce -resp_no_certs
```

## ecdsap384_7days.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-05-01 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
     -issuer CA/inter/cert.pem -cert certs/issued/ecdsap384_sxg_60days.crt \
    -ndays 7 -respout ocsp/ecdsap384_7days.ocsp -no_nonce -resp_no_certs
```

## revoked_7days.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-05-01 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
     -issuer CA/inter/cert.pem -cert certs/issued/ecdsap256_sxg_revoked.crt \
    -ndays 7 -respout ocsp/revoked_7days.ocsp -no_nonce -resp_no_certs
```

## certmanager_0401_0401.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-04-01 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
    -issuer CA/inter/cert.pem -cert certs/issued/certmanager_0401_15days.crt \
    -ndays 7 -respout ocsp/certmanager_0401_0401.ocsp -no_nonce -resp_no_certs
```

## certmanager_0401_0409.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-04-09 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
    -issuer CA/inter/cert.pem -cert certs/issued/certmanager_0401_15days.crt \
    -ndays 7 -respout ocsp/certmanager_0401_0409.ocsp -no_nonce -resp_no_certs
```

## certmanager_0401_0413.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-04-13 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
    -issuer CA/inter/cert.pem -cert certs/issued/certmanager_0401_15days.crt \
    -ndays 7 -respout ocsp/certmanager_0401_0413.ocsp -no_nonce -resp_no_certs
```

## certmanager_0415_0415.ocsp

```shell
NO_FAKE_STAT=1 TZ=UTC faketime '2020-04-15 00:00:00' \
    openssl ocsp -index CA/inter/index.txt -CA CA/inter/cert.pem \
    -rkey CA/inter/key.pem -rsigner CA/inter/cert.pem \
    -issuer CA/inter/cert.pem -cert certs/issued/certmanager_0415_15days.crt \
    -ndays 7 -respout ocsp/certmanager_0415_0415.ocsp -no_nonce -resp_no_certs
```
