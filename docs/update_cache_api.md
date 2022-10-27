# Signed Exchange Update Cache API Reference

Allows site owners to request deletion of a
[signed exchange](https://web.dev/signed-exchanges/) (SXG)
document in the Google SXG cache through the use of an HTTP API. After deletion, a subsequent
request for the signed exchange from the cache will return a cache miss and will redirect the
user to the original URL. The cache will fetch a new copy of the signed exchange in the background.
    
# Quick Start
      
Generate public/private key pairs. This example is for EC private keys, but you can also use
RSA keys. After generating, make sure to keep them in a safe place.  You will use the private
key to create requests to update the Google cache, while Google will use your public key to
verify that the request was authentic.

```
$ openssl ecparam -genkey -name prime256v1 -out ecprivkey.pem
$ openssl ec -in ecprivkey.pem -pubout -out ecpubkey.pem
```

Place your public key in your website’s ./well-known directory.  Google servers will fetch the
public key from this location when it’s time to verify your request.  Google will also keep a
copy of your public key in its cache with a 24 hour expiration date. If your website is
"[https://www.example.com](https://www.example.com)" then the URL for your public key should be:
“[https://www.example.com/.well-known/sxg-update-publickey.pem](https://www.example.com/.well-known/sxg-update-publickey.pem)".

Send the command to do the actual deletion. Replace the URL in the following script with your
origin URL, and CURL_URL with the corresponding SXG cache URL, as described in the 
[details section](#details).

```
#!/bin/bash
URL="https://example.com/document"
CURL_URL="https://example-com.webpkgcache.com/doc/-/s/example.com/document"
TIMESTAMP=$(date +%s)
MESSAGE=$({ \
  echo -n $URL; \
  echo -n " "; \
  echo -n $TIMESTAMP; })
SIGNATURE=$( \
  echo -n $MESSAGE | openssl dgst -sign ecprivkey.pem | base64 -w0)
curl \
-X DELETE \
--data-urlencode "timestamp=$TIMESTAMP" \
--data-urlencode "signature=$SIGNATURE" \
$CURL_URL
``` 

If the delete request was successful, the response status code will be 202. Otherwise, the
body will be a JSON message indicating the reason.

If the response indicates the signature was invalid, you can download and build the
[verify signature](https://github.com/google/libsxg/blob/main/src/verifysignature.c)
utility and use it for debugging. You can either
download just the source file and build it using gcc or clang, or you can follow the build
instructions for the entire package in
[https://github.com/google/libsxg#readme](https://github.com/google/libsxg#readme).
Once you’ve built it, verify that the utility works by running:

```
$ ./verifysignature sha256 ecpubkey.pem <(echo -n
'https://www.example.com/document.sxg 12345678' | openssl dgst -sha256
-sign ecprivkey.pem) <(echo -n 'https://www.example.com/document.sxg 12345678')
Signature verified. OK
```

# Details

## <b>HTTP Method</b>: DELETE

<b>HTTP Endpoint</b>: The same URL that represents the SXG. For example, the URL at 
[https://www.example.com/index.html](https://www.example.com/index.html)
is cached at the URL <code>https://www-example-com.webpkgcache.com/doc/-/s/www.example.com/index.html</code>
per this [algorithm](https://developer.google.com/search/docs/advanced/experience/signed-exchange#debug-the-google-sxg-cache). To purge the cache entry, call the same endpoint with a DELETE request method.

To determine the subdomain, do one of the following:
- Use the [SXG Validator Chrome extension](https://chrome.google.com/webstore/detail/sxg-validator/hiijcdgcphjeljafieaejfhodfbpmgoe) to determine the cache URL.
- GET the superdomain URL (such as <code>https://webpkgcache.com/doc/-/s/www.example.com/index.html</code>) and follow the redirect.
- Use this [algorithm](https://github.com/google/sxg-validator/blob/2c738d64ef6848e0074cd6ad5ad2a77e00f1f30f/dialog.js#L42-L64) to calculate it locally.

(Note that subresources are cached at a different path than documents. The best way to update
them is to purge the documents that embed them [in order to update their subresource integrity
digests]. For an intermediary that doesn't know whether a resource is a document or a
subresource, it is safe to request a deletion at the /doc/ URL; it will be a no-op for
subresources.)

## DELETE request parameters

The parameters should be specified via the request body, encoded as <code>application/x-www-form-urlencoded</code>.

- <b>signature</b> (required) 
  This parameter represents the signature of the entire request path. The signature is
  generated using the private key and the public key must be available on the origin at
  /.well-known/sxg-update-publickey.pem. It consists of the endpoint URL followed by a space
  followed by a timestamp, signed using the private key and then base64 encoded.
- <b>timestamp</b> (required)
  A decimal integer representing seconds since epoch. The request will be considered valid
  only if the timestamp is within 5 minutes of the current time, as known by the cache server.

There may be backwards-compatible changes to this format (e.g. new optional parameters).
The `version` parameter is reserved for incompatible changes in the future. (Its value is
currently unspecified.)

## Public key file
The .well-known resource must be an ASCII encoding of a
[PEM-encoded](https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail) file containing 1 to 10 unencrypted public keys in [SubjectPublicKeyInfo](https://datatracker.ietf.org/doc/html/rfc5280#section-4.1)
format. The signature will be valid if it matches any of the keys. (The goal of multiple keys is to allow cache purge by the end-point or anybody in its chain of authorized gateways. CDNs could concatenate their own public key to whatever the origin serves.)

Supported signature algorithms include [rsaEncryption](https://datatracker.ietf.org/doc/html/rfc8017#appendix-A.1) and [id-ecPublicKey](https://datatracker.ietf.org/doc/html/rfc5480#section-2.1.1)
with a parameter of [secp256r1](https://datatracker.ietf.org/doc/html/rfc5480#section-2.1.1.1)
(aka prime256v1). To allow for future updates, the cache will merely skip keys in
the file using other algorithms, rather than mark the entire file as invalid.

## Signature format
The signed message should be the request path exactly as sent to webpkgcache.com, e.g. 
<code>`/doc/-/s/signed-exchange-testing.dev/sxgs/valid.html`</code>, followed by a space (ASCII 0x20), followed by the value of the timestamp parameter. The signature should be base64-encoded using the
[URL-safe without base64 padding](https://datatracker.ietf.org/doc/html/rfc4648#section-5) variant.

# Output / Errors

The delete cache key operation is a high latency operation, involving deleting the
signed exchange resource from several edge nodes and multi-homed storage servers. Thus,
this API does not block on completion. If the deletion request has been successfully
initiated, it returns immediately with an HTTP 202 response, pending the actual delete which
may take several seconds.

After the deletion has been initiated, it may not complete in very rare cases. For instance:

- If the cluster manager preempts the request handling process, it may not achieve quorum.
- If a datacenter is netsplit, a caching tier may fail to clear for a short time.

All responses will be of type application/json containing the following fields:
- <b>success</b>: true or false
- <b>message</b>: A description if the success is false. For example: “The SXG URL is not
  found in the cache.”, or “Invalid URL signature, using public key &lt;key url&gt;”.

success: true is returned only for 202 responses. So the client may avoid parsing JSON upon
seeing a 202 response. All other non-202 responses will return JSON response
<code>{success: false, message: “&lt;reason for failure&gt;”}</code>

# Stay informed
Subscribe to the [webpackaging-announce](https://groups.google.com/g/webpackaging-announce) mailing list to stay up-to-date on significant changes to the update cache API.

If you have questions about SXG on Google Search, visit the [Search Central Help Community](https://support.google.com/webmasters/community).
