#!/bin/bash

set -e

gen-signedexchange \
    -uri https://example.org/standalone.html \
    -responseHeader 'Cache-Control: public, max-age=604800' \
    -responseHeader 'Content-Length: 37' \
    -responseHeader 'Content-Type: text/html;charset=utf-8' \
    -content standalone.html \
    -date 2019-04-22T19:30:00Z \
    -expire 168h \
    -certUrl https://example.org/cert.cbor \
    -certificate ../certs/test.pem \
    -privateKey ../certs/test.key \
    -validityUrl https://example.org/standalone.html.validity.1555961400 \
    -o standalone.sxg

gen-signedexchange \
    -uri https://example.org/preloading.html \
    -responseHeader 'Cache-Control: public, max-age=604800' \
    -responseHeader 'Content-Length: 78' \
    -responseHeader 'Content-Type: text/html;charset=utf-8' \
    -responseHeader 'Link: <https://example.org/style.css>;rel="allowed-alt-sxg";header-integrity="dummy-integrity"' \
    -responseHeader 'Link: <https://example.org/style.css>;rel="preload";as="style"' \
    -content preloading.html \
    -date 2019-04-22T19:30:00Z \
    -expire 168h \
    -certUrl https://example.org/cert.cbor \
    -certificate ../certs/test.pem \
    -privateKey ../certs/test.key \
    -validityUrl https://example.org/preloading.html.validity.1555961400 \
    -o preloading.sxg
    
gen-signedexchange \
    -uri https://example.org/incomplete.html \
    -responseHeader 'Cache-Control: public, max-age=604800' \
    -responseHeader 'Content-Length: 78' \
    -responseHeader 'Content-Type: text/html;charset=utf-8' \
    -responseHeader 'Link: <https://example.org/style.css>;rel="preload";as="style"' \
    -content incomplete.html \
    -date 2019-04-22T19:30:00Z \
    -expire 168h \
    -certUrl https://example.org/cert.cbor \
    -certificate ../certs/test.pem \
    -privateKey ../certs/test.key \
    -validityUrl https://example.org/incomplete.html.validity.1555961400 \
    -o incomplete.sxg
