# Web Packager

[![Build Status](https://travis-ci.org/google/webpackager.svg?branch=main)](https://travis-ci.org/google/webpackager)
[![GoDoc](https://godoc.org/github.com/google/webpackager?status.svg)](https://godoc.org/github.com/google/webpackager)

Web Packager is a command-line tool to "package" websites in accordance
with the specifications proposed at [WICG/webpackage][]. It may look like
[gen-signedexchange][], but is rather *based on* gen-signedexchange and
focuses on automating generation of Signed HTTP Exchanges (aka. SXGs) and
optimizing the page loading.

[WICG/webpackage]: https://github.com/WICG/webpackage/
[gen-signedexchange]: https://github.com/WICG/webpackage/tree/main/go/signedexchange

Web Packager HTTP Server is an HTTP server built on top of Web Packager.
It functions like a reverse-proxy, receiving signing requests over HTTP.
For more detail, see [cmd/webpkgserver/README.md](cmd/webpkgserver/README.md).
This README focuses on the command-line tool.

Web Packager retrieves HTTP responses from servers and turns them into
signed exchanges. Those signed exchanges are written into files in a way to
preserve the URL path structure, so can be deployed easily in some typical
cases. In addition, Web Packager applies some optimizations to the signed
exchanges to help the content get rendered quicker.

Web Packager is purposed primarily for a showcase of how to speed up the
page loading with [privacy-preserving prefetch][]. Web developers may port
the logic from this codebase to their systems or integrate Web Packager into
their systems. The Web Packager's code is designed to allow some injections
of custom logic; see the GoDoc comments for details. Note, however, that Web
Packager is currently at an early stage now: see Limitations below.

[privacy-preserving prefetch]: https://wicg.github.io/webpackage/draft-yasskin-webpackage-use-cases.html#private-prefetch

Web Packager is *not* related to [webpack][].

[webpack]: https://webpack.js.org/


## Prerequisite

Web Packager is written in the Go language thus requires a Go system to run.
See [Getting Started on golang.org](https://golang.org/doc/install) for how
to install Go on your computer.

You will also need a certificate and private key pair to use for the signing
the exchanges. Note the certificate must:

*   use an ECDSA private key (e.g. prime256v1) and
*   have [CanSignHttpExchanges extension][].

(For example, [DigiCert][] offers the right kind of certificates.)

[CanSignHttpExchanges extension]: https://wicg.github.io/webpackage/draft-yasskin-http-origin-signed-responses.html#cross-origin-cert-req
[DigiCert]: https://www.digicert.com/account/ietf/http-signed-exchange.php

Then you will need to convert your certificate into
the [application/cert-chain+cbor][] format, which you can do using the
instructions at:

*   [Creating our first signed exchange](https://github.com/WICG/webpackage/blob/main/go/signedexchange/README.md#creating-our-first-signed-exchange)
    to generate a self-signed certificate for testing.
*   [Creating a signed exchange using a trusted certificate](https://github.com/WICG/webpackage/blob/main/go/signedexchange/README.md#creating-a-signed-exchange-using-a-trusted-certificate)
    to use a CA-issued certificate.

[application/cert-chain+cbor]: https://wicg.github.io/webpackage/draft-yasskin-http-origin-signed-responses.html#cert-chain-format


## Limitations

In this early phase, we may make backward-breaking changes to the commandline
or API.

Web Packager aims to automatically meet most but not all [Google SXG Cache
requirements](docs/cache_requirements.md). In particular, pages that do not use
responsive design should specify a [`supported-media`
annotation](docs/supported_media.md).

Web Packager does not handle [request matching][] correctly. It should not
matter unless your web server implements content negotiation using the
`Variants` and `Variant-Key` headers (*not* the `Vary` header). We plan to
support the request matching in future, but there is no ETA (estimated time of
availability) at this moment.

**Note:** The above limitation is not expected to be a big deal even if your
    server serves signed exchanges conditionally using content negotiation:
    if you already have signed exchanges, you should not need Web Packager.

[request matching]: https://wicg.github.io/webpackage/loading.html#request-matching


## Install

```shell
go get -u github.com/google/webpackager/cmd/...
```


## Usage

The simplest command looks like:

```shell
webpackager \
    --cert_cbor=cert.cbor \
    --private_key=priv.key \
    --cert_url=https://example.com/cert.cbor \
    --url=https://example.com/hello.html
```

It will retrieve an HTTP response from https://example.com/, generate
a signed exchange with the given pair of certificate (`cert.cbor`) and
private key (`priv.key`), then write it to `./sxg/hello.html.sxg`.
If `hello.html` had subresources that could be preloaded together,
`webpackager` would also retrieve those resources and generate their signed
exchanges under `./sxg`. Web Packager recognizes `<link rel="preload">`
and equivalent `Link` HTTP headers. It also adds the preload links for CSS
(stylesheets) used in HTML, and may use more heuristics in future. See the
defaultproc package to find how exactly the HTTP response is processed.

`--cert_url` specifies where the client will expect to find the CBOR-format
certificate chain. `--cert_cbor` is optional when it can be fetched from
`--cert_url`. Note the reverse is not true: `--cert_url` is always required.

The `--url` flag can be repeated as many times as you want. For example:

```shell
webpackager \
    --cert_cbor=cert.cbor \
    --private_key=priv.key \
    --cert_url=https://example.com/cert.cbor \
    --url=https://example.com/foo/ \
    --url=https://example.com/bar/ \
    --url=https://example.com/baz/
```

would generate the following three files:

*   `./sxg/foo/index.html.sxg` for `https://example.com/foo/`
*   `./sxg/bar/index.html.sxg` for `https://example.com/bar/`
*   `./sxg/baz/index.html.sxg` for `https://example.com/baz/`

**Note:** `webpackager` expects all target URLs to have the same origin.
    In particular, the output files collide if you specify more than one URL
    that has the same path but a different domain.

### Using URL File

`webpackage` also accepts `--url_file=FILE`. `FILE` is a plain text file
with one URL on each line. For example, you could create `urls.txt` with:

```
# This is a comment.
https://example.com/foo/
https://example.com/bar/
https://example.com/baz/
```

then run:

```
webpackager \
    --cert_cbor=cert.cbor \
    --private_key=priv.key \
    --cert_url=https://example.com/cert.cbor \
    --url_file=urls.txt
```

### Changing Output Directory

You can change the output directory with the `--sxg_dir` flag:

```shell
webpackager \
    --cert_cbor=cert.cbor \
    --private_key=priv.key \
    --cert_url=https://example.com/cert.cbor \
    --sxg_dir=/tmp/sxg \
    --url=https://example.com/hello.html
```

### Setting Expiration

The signed exchanges last one hour by default. You can change the duration
with the `--expiry` flag. For example:

```shell
webpackager \
    --cert_cbor=cert.cbor \
    --private_key=priv.key \
    --cert_url=https://example.com/cert.cbor \
    --expiry=72h \
    --url=https://example.com/hello.html
```

would make the signed exchanges valid for 72 hours (3 days). The maximum
is `168h` (7 days), due to the specification.

### Other Flags

`webpackager` provides more flags for advanced usage (e.g. to set request
headers). Run the tool with `--help` to see those flags.


## Appendix: Deploying SXGs

The steps below illustrate an example of deploying Signed HTTP Exchanges on
an Apache server.

1.  Upload `cert.cbor` to your server. Make it available at `--cert_url`.

2.  Upload `*.sxg` files to your server. Put them next to the original files
    (e.g. `hello.html.sxg` should stay in the same directory as `hello.html`).
    For example, if you are using the `sftp` command to upload, you can:

    ```
    sftp> cd public_html
    sftp> put -r sxg/*
    ```

    assuming `public_html` to be the document root and `sxg` to be where you
    generated the `*.sxg` files.

3.  Edit or create `.htaccess` in `public_html` (or the Apache's config file)
    to add the following settings:

    ```
    AddType application/signed-exchange;v=b3 .sxg

    <Files "cert.cbor">
      AddType application/cert-chain+cbor .cbor
    </Files>

    RewriteEngine On
    RewriteCond %{HTTP:Accept} application/signed-exchange
    RewriteCond %{REQUEST_FILENAME} !\.sxg$
    RewriteCond %{REQUEST_FILENAME}\.sxg -s
    RewriteRule .+ %{REQUEST_URI}.sxg [L]

    Header set X-Content-Type-Options: "nosniff"
    ```
