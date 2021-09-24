# supported-media meta tag

## Status

This is a draft proposal. It is not yet implemented anywhere.

Google Search may implement it or some subset in the future. To increase chance
of forward compatibility, use the following exact tag for mobile-only HTML:

```html
<meta name=supported-media content="only screen and (max-width: 640px)">
```

and the following exact tag for desktop-only HTML:

```html
<meta name=supported-media content="only screen and (min-width: 640px)">
```

This aligns with their [`rel=alternate`
recommendation](https://developers.google.com/search/mobile-sites/mobile-seo/separate-urls#annotation-in-the-html).

## Problem

Many sites rely on serving different HTML to different devices (e.g.
mobile/desktop). This is typically done by server-side `User-Agent` (UA)
sniffing.  This technique conflicts with caching by upstream intermediaries,
but the purposes served by both are important.

## Proposed solution

Pages that `Vary` by UA should specify the set of devices that they support via
a new meta tag. For example:

```html
  <meta name=supported-media content="(max-width: 8in) and (hover: none)">
```

Referrers should not link to cached copies unless the device meets the given
query.

User agents should ignore this meta tag.

## Syntax

A `<meta>` element, as a child of the document's `<head>` element, where its
`name=` is "supported-media", and its `content=` is a well-formed [level 3
media query list](https://www.w3.org/TR/2012/REC-css3-mediaqueries-20120619/)
in the same character encoding as the page, less than 200 characters, and
containing only the below supported media types and features:

 - `all`
 - `screen`
 - `hover`
 - `pointer`
 - the following, plus their `min-` and `max-` variants:
   - `height`
   - `width`
   - `aspect-ratio`
   - `device-height`
   - `device-width`
   - `device-aspect-ratio`

Referrers should make no restrictions as to the byte position of this element.

## Semantics

If the meta element is not present or doesn't meet the above syntax
requirements, referrers may assume that the page renders on all devices.

For pages that already vary by `User-Agent` header, this media query would
likely not exactly match the behavior of their UA sniffing algorithm.
Publishers should err on the side of specificity over sensitivity.  The cost of
a false positive in the media query is the page renders incorrectly; the cost
of a false negative is the page renders slower (because the cache entry is
skipped and the referrer links to origin).

## Analysis

See [this analysis](supported_media_analysis.md) of the limitations of this
specification and its alternatives.
