# Analysis of supported-media meta tag

This is an analysis of the [supported-media specification](supported_media.md).

## Limitations

This proposal requires client-side code on the referrer to evaluate the media
query. Thus, it is not able to address the use-case for same-origin caching
gateways (e.g. CDNs). Doing so is out of scope for this document; hopefully
other proposals such as [Sec-CH-UA](https://github.com/WICG/ua-client-hints)
can address this eventually.

This proposal does not allow publishers to enumerate all of a page's form
factor variants. That would be an interesting area for further design. Such a
specification would allow increased cache hit rate, as publishers could
instruct caches to fetch multiple variants per page. Without this, caches could
heuristically attempt to enumerate the variants through multiple fetches, but
shouldn't; doing so would introduce unnecessary traffic. Hopefully other
proposals such as
[Variants](https://tools.ietf.org/html/draft-ietf-httpbis-variants-06) can
address this eventually.

Referrers need to evaluate the intersection of user-agent support and page
support. Referrers with their own separate crawl can do so locally,
independently of the cache. Otherwise, they will need to query the cache for
its knowledge of the cached page's support. Such a query could be a batched or
streamed push from caches to referrers of recently ingested documents, or
referrers might make bulk requests for new and refreshed documents, per their
crawl schedulers. In either case, both entities will need to guard against
version skew.

## Privacy considerations

Allowing the referrer to evaluate the `supported-media` server-side would
reduce the number of bytes needed in the referring page, but would require
browser support to send sufficient information. This would introduce another
passive fingerprinting vector; this is at odds with efforts to limit
fingerprinting such as [Privacy
Budget](https://github.com/bslassey/privacy-budget) and [Brave's
protections](https://github.com/brave/brave-browser/wiki/Fingerprinting-Protections).

In order to evaluate the intersection of user-agent support and page support,
referrers could send RPCs to caches, containing (client information, referent
URL) tuples. They shouldn't. Doing so widens the sphere of knowledge of this
information beyond the referrer. This specification does not change the risk:

 1. It is already possible, and:
 2. Referrers have no increased incentive to do so, as the alternative approach
    (outlined in [Limitations](#limitations)) is both easier to implement and a
    better user experience, and does not change which entities know which
    information about which users.

    It is easier to implement because referrers and caches can rely on existing
    browser support for evaluating media queries, rather than developing
    client-side code for collecting device information, protocols for
    serializing this information, and libraries for evaluating it server-side.
    It is also easier to maintain, as it is more resilient to changes in
    browser behavior that reduce scripting access to this information.

    It is a better user experience because it can be evaluated in
    microseconds, rather than hundreds of milliseconds (requiring several RPC
    hops to the referrer and from there to the cache).

Referent pages could determine whether they've been loaded from a cross-origin
cache or from origin, and evaluate the media query locally, and compare the
responses. If they differ, pages might learn some information about user
preferences on the referrer. For instance, perhaps a User Agent offers
per-origin configuration of `orientation` or `prefers-color-scheme` and the
user has configured that specifically for the referrer. For a given user this
is a fractional bit, but in aggregate may reveal a distribution. Referrers
should be aware of this risk.

## Alternatives considered

### Syntax

#### UA sniffing

Referrers or caches would need to run arbitrary sniffing algorithms. These can
be arbitrarily complex (read: Turing complete), hence sandboxing would be
necessary.

Publishers would need to upload their sniffing algorithms to caches, either
through bilateral coordination, or by publishing them on the web somewhere
(e.g. a `.well-known` URI). In some cases, this may be infeasible, e.g. if the
publisher relies on a closed-source third-party sniffer.

While it would allow for more flexibility in responsiveness than media queries,
its complexity is in disproportion to the need.

These considerations don't change with regard to the `Sec-CH-UA` proposal.

#### Arbitrary media queries

Referrers could allow any media queries rather than the subset specified above.
This creates a few risks:

  - Queries might evaluate incorrectly. Of the set of browsers for which a
    referrer provides caching, the referrer would need to keep track of which
    media types/features each supports, in order not to [fail
    closed](#fail-closed-on-malformed-media-queries). This complexity is
    disproportionate to the need.
  - Queries might hurt user experience on the referring page. For instance,
    overly long queries may peg the CPU. (Or perhaps there are some
    expensive-to-compute features? These features might be safe to add to the
    web platform where a page author may choose the appropriate resource usage
    trade-offs, but not in a shared environment containing multiple queries
    from different authors.)
  - Queries might vary over time. For instance, `orientation` may evaluate to
    `portait` while on the referrer, but the user may change it to `landscape`
    after visiting the referent. (This is true of some of the allowed subset,
    too, but it seemed like a reasonable trade-off between utility and risk of
    inaccuracy.)

The above subset was chosen as a safe starting point, based on a cursory
examination of common media queries. Proposals for additions are welcome.

#### Predefined aliases

The supported-media value could be a comma-separated list of form factors such
as `desktop`, `mobile`, and `tablet`. Referrers and publishers would need to
standardize such a list. This is likely difficult to do for many reasons; for
instance, different regions may have different standards for mobile sizes. In
addition, doing so risks ossification of the predefined aliases, which then may
become outdated as user expectations change; for instance, phones might change
size distribution, or users might expect specialized [foldable
interfaces](https://github.com/MicrosoftEdge/MSEdgeExplainers/blob/main/Foldables/explainer.md)
in the future.

#### meta name=viewport

Referrers could use `<meta name=viewport>` as an indication that a site is
responsive. This is an opt-in mechanism; lacking this signal, referrers would have to
assume the site is non-responsive. See the [Assume non-responsive by
default](#assume-non-responsive-by-default) for why this is not preferred.

#### Response header

Referrers could look for a `supported-media` HTTP response header rather than a
`<meta>` tag. This has the advantage of supporting non-HTML referents, such as
PDFs; however, those are rarely published in multiple form factors.

It risks increased maintenance burden for publishers; HTML and HTTP teams often
being different, there would be communication overhead in maintaining the
proper mapping from page to appropriate response header.

Referrers could offer the HTTP header as an optional alternative serialization
to the `<meta>` tag (i.e. allow either spelling). Lacking evidence of need, it
seems safest to simplify to a single spelling for now. It is easier to add
options later than to remove them.

#### Variants
a
Referrers could look for a hypothetical `Variants` header such as:

```http
variants: User-Agent=(desktop mobile)
variant-key: (mobile)
```

This combines the risks of the [Response header](#response-header) approach and
the [Predefined aliases](#predefined-aliases) approach, as there is no standard
[content negotiation
mechanism](https://tools.ietf.org/html/draft-ietf-httpbis-variants-06#appendix-A)
for `User-Agent`.

### Semantics

#### Nothing

Referrers could offer no support for UA-Varying pages; publishers wishing to opt
into cross-origin caching would need to convert any such pages to be responsive
or opt them out of caching. This creates several risks:

  - Converting a template to be responsive can be costly compared to auditing
    its responsiveness and adding the appropriate meta tag.
  - HTML and HTTP teams are often different. It could be an ongoing maintenance
    burden for the HTTP team to maintain a list of pages to exclude from
    cross-origin caching, where such knowledge is better maintained by the
    HTML team.
  - Responsive pages require more bytes (either upfront or lazy-loaded) than
    non-responsive. The existience of [media-dependent CSS
    imports](https://developer.mozilla.org/en-US/docs/Web/CSS/@import) is one
    indication of that. See also JS libraries for mobile- or desktop- specific
    needs. The increase may counteract the benefit of cross-origin caching.
  - Some consider responsive design to be at odds with accessibility, for a
    variety of reasons, such as:
    - Encouraging reuse of design patterns that are not equally accessible for
      both mobile and desktop (e.g. Fitts's law applies differently).
    - Negative interactions with screen readers.

In short, it is best in the hands of publishers to make trade-off decisions
that work for their individual cases.

#### Assume non-responsive by default

If the `meta` element is not present or not well-formed, referrers could assume
that the page is non-responsive. Not knowing for which devices it's suited, the
referrer could never link to the cache entry.

The purpose of such a default would be to decrease the risk that existing pages
are viewed on unsupported media. There is already a mitigation of this risk:
cross-origin caching is only viable under some publisher opt-in such as [signed
exchanges](https://wicg.github.io/webpackage/draft-yasskin-http-origin-signed-responses.html)
or [content-based
origins](https://tools.ietf.org/html/draft-thomson-wpack-content-origin-00).
Publishers need to factor in auditing their pages for UA-variance into the cost
of opting into cross-origin caching.

On the other hand, changing the default to be non-responsive would increase the
cost:benefit ratio for publishers. Those who already know their pages are
responsive would need to modify their templates in order for the caching opt in
to be effective. In either default, those who don't know would need to audit.

#### Assume mobile or desktop by default, by heuristic

Similarly to above, referrers could assume that a non-annotated page is mobile
or desktop by default, based on proprietary heuristics such as
[rel=alternate/canonical](https://developers.google.com/search/mobile-sites/mobile-seo/separate-urls).
Those heuristics work well in environment without cross-origin caching; if a
referrer makes an incorrect call, the publisher can sniff the UA server-side
and send a redirect to the appropriate page. When server-side UA-Varying is not
available, the publisher would need to either:

  - Accept the risk of incorrect rendering
  - Implement client-side UA sniffing with redirect (this costs bytes)
  - Implement a client-side XHR to lazy-UA-sniff (this costs latency)

#### Fail closed on malformed media queries

Referrers, on seeing pages with `supported-media` not meeting the media query
syntax requirements, could assume the pages are non-responsive and thus
available on no media. Doing so increases the risk of ossifying the format.

Consider the following series of events:

  1. A page includes in its `supported-media` a nascent media feature.
  2. Referrers never refer to its cached copy.
  3. The publisher tests its site, and assumes the query is working as
     intended; the fact that the origin copy is loaded rather than the cached copy
     is easy to miss.
  4. After standardization, referrers allow the feature.
  5. The media query starts taking effect, and the cached page loads, but the
     query was buggy, and so the page displays on media for which it was not
     designed.

If this scenario happens to a significant number of pages, referrers would be
reticent to allow the new feature, as it harms user experience on pages
referred from their site.

On the other hand, if step 2 was "Referrers always refer to its cached copy",
then step 3 will result in a more visible error (the copy loading on incorrect
media), and thus increase the chance that the bug is fixed before step 5.

This is different from [CSS
error-handling](https://www.w3.org/TR/2012/REC-css3-mediaqueries-20120619/#error-handling).
The reason for that difference is in the result of the error-handling. In CSS,
whether the fallback value is `true` or `false`, the result is an incorrect
display. In prefetching, the result may be harder for the publisher to notice,
as outlined above.
