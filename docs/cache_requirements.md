# Audience

The audience for this document is people intending on implementing their own
signed exchange generator, independent of webpackager, and those implementing
their own SXG cache for the purposes of privacy-preserving prefetch. Users of
webpkgserver need not read this, as the tool should automatically guarantee the
following requirements are met.

# Google SXG cache

The Google SXG cache sets these requirements in addition to the ones set by the
[SXG spec][]:
 - The SXG must have a freshness lifetime of at least 120 seconds, as [computed
   for a shared cache](https://tools.ietf.org/html/rfc7234#section-4.2.1) from
   its outer headers.
 - The signed `fallback URL` must approximately equal the URL at which the SXG
   was served. Where possible, aim to make them byte-equal. The set of allowed
   differences is not precisely specified, but approximately:
   - Characters may be substituted by their percent encodings, and vice versa,
     with the exception of meaningful delimiters like `/`, `;`, `?`, `&`, and
     `=`.
   - Query parameters may be re-ordered.
   - Valueless query parameters may be encoded with or without a trailing `=`.
   - Extra `&`s in the query string are allowed.
 - The signed `cert-url` must be `https`.
 - The signature header must contain only:
   - One parameterised identifier.
   - Parameter values of type string, binary, or identifier.
 - The payload must be non-empty.
 - The signed `cache-control` header cannot have a `no-cache` or `private`
   directive, even with a value (e.g. `no-cache=some-header` is disallowed).
 - The `content-type` must satisfy the [`media-type` grammar][].
 - The `link` header, if present, must lead to successful substitution per the
   [Loading spec][].
   Specifically, it must meet these requirements, in addition to the ones set by
   the [Link spec][]:
   - Each `URI-Reference` must be an absolute `https` URL.
   - Parameter names can only be `as`, `header-integrity`, `media`, `rel`,
     `imagesrcset`, `imagesizes`, or `crossorigin`.
   - All `rel` parameters must be either `preload` or `allowed-alt-sxg`.
   - All `imagesrcset` values must parse as a [srcset attribute](https://html.spec.whatwg.org/multipage/images.html#srcset-attribute).
   - There may be no more than 20 `rel=preload`s.
   - All `crossorigin` values must either be the empty string, or `anonymous`.
   - Every `rel=preload` must have a corresponding `rel=allowed-alt-sxg` with
     the same URI, which in turn must contain a `header-integrity` parameter
     with a value that satisfies the [CSP `hash-source` grammar](https://w3c.github.io/webappsec-csp/#grammardef-hash-source)
     using the `sha256` variant.
   - The preloaded URLs, when requested with an [SXG-preferring `Accept` header][],
     must respond with valid SXGs that match their given `header-integrity`.
 - The `link` header must not be present on subresources, i.e. SXGs that are
   themselves preloaded from other SXGs.
 - There must not be a signed `variant-key-04` or `variants-04` header.
 - The signature's lifetime (`expires` minutes request time) must be >= 120
   seconds.
 - The SXG must be no larger than 8 megabytes.
 - The page should be responsive, i.e. correct on all media. (In the future, a
   [supported-media](supported_media.md) annotation should allow this
   constraint to be removed.)

[SXG spec]: https://wicg.github.io/webpackage/draft-yasskin-http-origin-signed-responses.html
[`media-type` grammar]: https://tools.ietf.org/html/rfc7231#section-3.1.1.5
[Loading spec]: https://wicg.github.io/webpackage/loading.html#subresource-substitution
[Link spec]: https://datatracker.ietf.org/doc/html/rfc8288#section-3
[SXG-preferring `Accept` header]: https://github.com/google/webpackager/tree/main/cmd/webpkgserver#content-negotiation

Some of the above limitations are overly strict for an SXG cache's needs, and
were implemented as such for the sake of expediency. They may be loosened over
time, especially in response to publisher feedback.

# Other SXG caches

Other SXG caches could define their own set of requirements. It would be most
useful for publishers and users, however, if the requirements were the same
across all caches. If you see a need for a different requirement on your cache,
please file an issue.

The Google AMP Cache has a different set of requirements for SXGs. See [advice
for sites with a mix of AMP and non-AMP](amp_cache_differences.md).

# Testing

For SXGs on the internet, one can use the [SXG Validator Chrome extension](https://chrome.google.com/webstore/detail/sxg-validator/hiijcdgcphjeljafieaejfhodfbpmgoe). This queries the Google SXG Cache to see if the SXG meets the above requirements.

Alternatively, one can query the cache directly. This is an example that meets the requirements:

```
$ curl -siH 'Accept: application/signed-exchange;v=b3' https://signed--exchange--testing-dev.webpkgcache.com/doc/-/s/signed-exchange-testing.dev/sxgs/valid.html | grep -aiE 'content-type:|warning:'
content-type: application/signed-exchange;v=b3
```

and this does not meet requirements:

```
$ curl -siH 'Accept: application/signed-exchange;v=b3' https://signed--exchange--testing-dev.webpkgcache.com/doc/-/s/signed-exchange-testing.dev/sxgs/invalid-signature-date.html | grep -aiE 'content-type:|warning:'
warning: 199 - "debug: content has ingestion error: SXG ingestion failure: sig_date is in the future"
content-type: text/html; charset=UTF-8
```

Cache ingestion is asynchronous. Documents not yet ingested will have a `text/html` content type but no `Warning` header.
