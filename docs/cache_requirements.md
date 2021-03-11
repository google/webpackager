# Audience

The audience for this document is people intending on implementing their own
signed exchange generator, independent of webpackager, and those implementing
their own SXG cache for the purposes of privacy-preserving prefetch. Users of
webpkgserver need not read this, as the tool should automatically guarantee the
following requirements are met.

# Google SXG cache

The Google SXG cache sets these requirements in addition to the ones set by the
[SXG spec][]:
 - The signed `fallback URL` must equal the URL at which the SXG was served.
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
     `imagesrcset`, or `imagesizes`.
   - All `rel` parameters must be either `preload` or `allowed-alt-sxg`.
   - All `imagesrcset` values must parse as a [srcset attribute](https://html.spec.whatwg.org/multipage/images.html#srcset-attribute).
   - There may be no more than 20 `rel=preload`s.
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
[Link spec]: https://tools.ietf.org/html/rfc5988#section-5
[SXG-preferring `Accept` header]: https://github.com/google/webpackager/tree/master/cmd/webpkgserver#content-negotiation

Some of the above limitations are overly strict for an SXG cache's needs, and
were implemented as such for the sake of expediency. They may be loosened over
time, especially in response to publisher feedback.

# Other SXG caches

Other SXG caches could define their own set of requirements. It would be most
useful for publishers and users, however, if the requirements were the same
across all caches. If you see a need for a different requirement on your cache,
please file an issue.

# Testing

There is no known publicly available tool for validating an SXG against the
above requirements, though one is certainly welcome. In the interim, one may
send an HTTP request to the Google SXG Cache and see if the response is a valid
SXG.

Meets requirements:

```
$ curl -s -i -H 'Accept: application/signed-exchange;v=b3' https://signed--exchange--testing-dev.webpkgcache.com/doc/-/s/signed-exchange-testing.dev/sxgs/valid.html | grep -a -i content-type:
content-type: application/signed-exchange;v=b3
```

Does not meet requirements:

```
$ curl -s -i -H 'Accept: application/signed-exchange;v=b3' https://signed--exchange--testing-dev.webpkgcache.com/doc/-/s/signed-exchange-testing.dev/sxgs/invalid-signature-date.html | grep -a -i content-type:
content-type: text/html; charset=UTF-8
```

Note that new documents may appear not to meet the requirements at first; cache
ingestion is asynchronous.
